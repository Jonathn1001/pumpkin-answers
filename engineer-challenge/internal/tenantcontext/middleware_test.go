package tenantcontext_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/tenantcontext"
	"github.com/gin-gonic/gin"
)

type fakeFetcher struct{ notFound bool }

func (f fakeFetcher) GetTenant(_ context.Context, slug string) (domain.Tenant, error) {
	if f.notFound {
		return domain.Tenant{}, domain.ErrTenantNotFound
	}
	return domain.Tenant{Slug: slug, Name: "X"}, nil
}

func serve(f tenantcontext.TenantFetcher, path string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/tenants/:slug/x", tenantcontext.Middleware(f), func(c *gin.Context) {
		tn, ok := tenantcontext.FromContext(c)
		if !ok {
			c.String(http.StatusInternalServerError, "no tenant")
			return
		}
		c.String(http.StatusOK, tn.Slug)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
}

func TestMiddlewareSetsTenant(t *testing.T) {
	w := serve(fakeFetcher{}, "/tenants/acme/x")
	if w.Code != http.StatusOK || w.Body.String() != "acme" {
		t.Fatalf("got %d %q", w.Code, w.Body.String())
	}
}

func TestMiddleware404OnUnknownTenant(t *testing.T) {
	w := serve(fakeFetcher{notFound: true}, "/tenants/ghost/x")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
