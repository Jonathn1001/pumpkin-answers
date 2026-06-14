package httpapi

import (
	"net/http"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/tenantcontext"
	"github.com/gin-gonic/gin"
)

// createTenantReq carries the full starter config from the client. config is
// optional: when omitted, the server seeds the tenant from DefaultConfig().
// The slug is always derived server-side from name (never client-supplied).
type createTenantReq struct {
	Name   string                 `json:"name" binding:"required"`
	Config *domain.ConfigDocument `json:"config"`
}

func (h *handlers) createTenant(c *gin.Context) {
	var req createTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	cfg := h.svc.DefaultConfig()
	if req.Config != nil {
		cfg = *req.Config
	}
	tn, err := h.svc.CreateTenant(c.Request.Context(), req.Name, cfg)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, tn)
}

// configDefault returns the starter config the create wizard edits when the
// user picks the "default" source.
func (h *handlers) configDefault(c *gin.Context) {
	c.JSON(http.StatusOK, h.svc.DefaultConfig())
}

func (h *handlers) listTenants(c *gin.Context) {
	ts, err := h.svc.ListTenants(c.Request.Context())
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, ts)
}

func (h *handlers) getTenant(c *gin.Context) {
	tn, _ := tenantcontext.FromContext(c) // middleware guarantees presence
	c.JSON(http.StatusOK, tn)
}

type updateTenantReq struct {
	Name   string `json:"name" binding:"required"`
	Status string `json:"status"`
}

func (h *handlers) updateTenant(c *gin.Context) {
	var req updateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	tn, err := h.svc.UpdateTenant(c.Request.Context(), c.Param("slug"), req.Name, req.Status)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, tn)
}
