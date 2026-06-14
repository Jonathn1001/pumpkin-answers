// Package httpapi adapts HTTP requests to use cases. Handlers stay thin; domain
// errors map to HTTP status here, in one place.
package httpapi

import (
	"errors"
	"net/http"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/tenantcontext"
	"claimsplatform/internal/usecase"
	"github.com/gin-gonic/gin"
)

type handlers struct{ svc *usecase.Service }

// NewRouter wires all /api routes onto a fresh Gin engine.
func NewRouter(svc *usecase.Service) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	h := &handlers{svc: svc}

	api := r.Group("/api")
	api.GET("/tenants", h.listTenants)
	api.POST("/tenants", h.createTenant)
	api.GET("/config-schema", h.configSchema)
	api.GET("/diff", h.diff)

	tg := api.Group("/tenants/:slug")
	tg.Use(tenantcontext.Middleware(svc))
	tg.GET("", h.getTenant)
	tg.PATCH("", h.updateTenant)
	tg.GET("/config", h.getActiveConfig)
	tg.GET("/versions", h.listVersions)
	tg.GET("/versions/:n", h.getVersion)
	tg.POST("/versions", h.createDraft)
	tg.POST("/versions/:n/publish", h.publish)
	tg.POST("/rollback", h.rollback)
	tg.POST("/process", h.process)
	tg.POST("/preview", h.preview)

	return r
}

// fail maps domain errors to HTTP responses (single source of HTTP status logic).
func fail(c *gin.Context, err error) {
	var ve domain.ValidationError
	switch {
	case errors.As(err, &ve):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "validation_failed", "message": "config validation failed", "fields": ve.Fields}})
	case errors.Is(err, domain.ErrTenantNotFound), errors.Is(err, domain.ErrVersionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": err.Error()}})
	case errors.Is(err, domain.ErrSlugTaken):
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "conflict", "message": err.Error()}})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal", "message": "internal server error"}})
	}
}

func badRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "bad_request", "message": err.Error()}})
}
