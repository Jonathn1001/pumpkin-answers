package httpapi

import (
	"net/http"
	"strconv"

	"claimsplatform/internal/domain"
	"github.com/gin-gonic/gin"
)

// actor reads the optional X-Actor header (no real auth; PRD out-of-scope).
func actor(c *gin.Context) string {
	if a := c.GetHeader("X-Actor"); a != "" {
		return a
	}
	return "system"
}

func (h *handlers) getActiveConfig(c *gin.Context) {
	cfg, err := h.svc.GetActiveConfig(c.Request.Context(), c.Param("slug"))
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *handlers) listVersions(c *gin.Context) {
	vs, err := h.svc.ListVersions(c.Request.Context(), c.Param("slug"))
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, vs)
}

func (h *handlers) getVersion(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil {
		badRequest(c, err)
		return
	}
	v, err := h.svc.GetVersion(c.Request.Context(), c.Param("slug"), n)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}

type createDraftReq struct {
	Config domain.ConfigDocument `json:"config"`
	Note   string                `json:"note"`
}

func (h *handlers) createDraft(c *gin.Context) {
	var req createDraftReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	v, err := h.svc.SaveDraftConfig(c.Request.Context(), c.Param("slug"), req.Config, req.Note, actor(c))
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *handlers) publish(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("n"))
	if err != nil {
		badRequest(c, err)
		return
	}
	if err := h.svc.PublishVersion(c.Request.Context(), c.Param("slug"), n); err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"published": n})
}

type rollbackReq struct {
	TargetVersion int `json:"targetVersion" binding:"required"`
}

func (h *handlers) rollback(c *gin.Context) {
	var req rollbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	v, err := h.svc.RollbackVersion(c.Request.Context(), c.Param("slug"), req.TargetVersion, actor(c))
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, v)
}
