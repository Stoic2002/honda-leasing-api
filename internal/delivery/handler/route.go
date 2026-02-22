package handler

import "github.com/gin-gonic/gin"

func (h *DeliveryHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc, rbac func(...string) gin.HandlerFunc) {
	v1 := router.Group("/api/v1/delivery")
	v1.Use(authMiddleware)
	v1.Use(rbac("SALES"))
	{
		v1.GET("/orders", h.GetDeliveryOrders)
		v1.GET("/tasks", h.GetDeliveryTasks)
		v1.POST("/tasks/:taskId/complete", h.CompleteDeliveryTask)
	}

	// Serve the static uploaded files so they can be viewed via browser
	router.Static("/uploads", "./uploads")
}
