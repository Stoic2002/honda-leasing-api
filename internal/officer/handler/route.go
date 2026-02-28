package handler

import "github.com/gin-gonic/gin"

func (h *OfficerHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc, rbac func(...string) gin.HandlerFunc) {
	v1 := router.Group("/api/v1/officer")
	v1.Use(authMiddleware)
	v1.Use(rbac("ADMIN_CABANG", "SALES", "SURVEYOR", "FINANCE", "COLLECTION"))
	{
		v1.GET("/contracts", h.GetIncomingContracts)
		v1.GET("/tasks", h.GetMyTasks)
		v1.POST("/tasks/:taskId/process", h.ProcessTask)
	}

	// Serve the static uploaded files so they can be viewed via browser
	router.Static("/uploads", "./uploads")
}
