// Package tenantcontext resolves the :slug path segment to a Tenant and stores it on
// the request, so handlers operate within one tenant's scope.
package tenantcontext

import (
	"context"
	"errors"
	"net/http"

	"claimsplatform/internal/domain"
	"github.com/gin-gonic/gin"
)

const tenantKey = "tenant"

// TenantFetcher is the narrow slice of the use case this middleware needs.
type TenantFetcher interface {
	GetTenant(ctx context.Context, slug string) (domain.Tenant, error)
}

// Middleware resolves :slug to a Tenant; 404 if not found, 500 on other errors.
func Middleware(f TenantFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		tn, err := f.GetTenant(c.Request.Context(), c.Param("slug"))
		if err != nil {
			if errors.Is(err, domain.ErrTenantNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound,
					gin.H{"error": gin.H{"code": "not_found", "message": "tenant not found"}})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				gin.H{"error": gin.H{"code": "internal", "message": err.Error()}})
			return
		}
		c.Set(tenantKey, tn)
		c.Next()
	}
}

// FromContext returns the tenant resolved by Middleware.
func FromContext(c *gin.Context) (domain.Tenant, bool) {
	v, ok := c.Get(tenantKey)
	if !ok {
		return domain.Tenant{}, false
	}
	tn, ok := v.(domain.Tenant)
	return tn, ok
}
