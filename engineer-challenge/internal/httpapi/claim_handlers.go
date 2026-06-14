package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"claimsplatform/internal/configschema"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/usecase"
	"github.com/gin-gonic/gin"
)

func (h *handlers) process(c *gin.Context) {
	var claim domain.Claim
	if err := c.ShouldBindJSON(&claim); err != nil {
		badRequest(c, err)
		return
	}
	dec, err := h.svc.ProcessClaim(c.Request.Context(), c.Param("slug"), claim)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, dec)
}

type previewReq struct {
	Claim         domain.Claim           `json:"claim"`
	VersionNumber *int                   `json:"versionNumber"`
	Config        *domain.ConfigDocument `json:"config"`
}

func (h *handlers) preview(c *gin.Context) {
	var req previewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	dec, err := h.svc.PreviewClaim(c.Request.Context(), c.Param("slug"), req.Claim, req.VersionNumber, req.Config)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, dec)
}

func (h *handlers) diff(c *gin.Context) {
	left, err := parseRef(c.Query("left"))
	if err != nil {
		badRequest(c, err)
		return
	}
	right, err := parseRef(c.Query("right"))
	if err != nil {
		badRequest(c, err)
		return
	}
	changes, err := h.svc.CompareConfigs(c.Request.Context(), left, right)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"changes": changes})
}

func (h *handlers) configSchema(c *gin.Context) {
	c.JSON(http.StatusOK, configschema.Get())
}

// parseRef parses "slug" (active config) or "slug@n" (version n) into a ConfigRef.
func parseRef(s string) (usecase.ConfigRef, error) {
	if s == "" {
		return usecase.ConfigRef{}, errors.New("ref is required")
	}
	parts := strings.SplitN(s, "@", 2)
	ref := usecase.ConfigRef{Slug: parts[0]}
	if len(parts) == 2 {
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return usecase.ConfigRef{}, fmt.Errorf("invalid version in ref %q", s)
		}
		ref.Version = &n
	}
	return ref, nil
}
