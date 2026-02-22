package handler

import "github.com/gin-gonic/gin"

func (h *LeasingHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc, rbac func(...string) gin.HandlerFunc) {
	v1 := router.Group("/api/v1/customer")
	v1.Use(authMiddleware)
	v1.Use(rbac("CUSTOMER"))
	{
		v1.GET("/orders", h.GetMyOrders)
		v1.POST("/orders", h.SubmitOrder)
		v1.GET("/orders/:id/progress", h.GetContractProgress)
	}
}
