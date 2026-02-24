package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *FinanceHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1/finance")
	{
		// Should be protected by RBAC
		v1.GET("/schedules", h.GetSchedules)

		// Typically secured by API Key or specific webhook auth
		v1.POST("/payments/webhook", h.ProcessWebhook)
	}
}
