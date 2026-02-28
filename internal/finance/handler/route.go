package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *FinanceHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	v1 := router.Group("/api/v1/finance")
	v1.Use(authMiddleware)
	{
		// Should be protected by RBAC
		v1.GET("/schedules", h.GetSchedules)
		v1.POST("/payments", h.ProcessPayment)
	}
}
