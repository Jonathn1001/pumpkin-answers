package httpapi

import (
	"net/http"

	"claimsplatform/internal/tenantcontext"
	"github.com/gin-gonic/gin"
)

type createTenantReq struct {
	Name      string `json:"name" binding:"required"`
	CloneFrom string `json:"cloneFrom"`
}

func (h *handlers) createTenant(c *gin.Context) {
	var req createTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	tn, err := h.svc.CreateTenant(c.Request.Context(), req.Name, req.CloneFrom)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, tn)
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
