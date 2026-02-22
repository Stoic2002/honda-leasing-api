package handler

import "github.com/gin-gonic/gin"

func (h *AuthHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	public := router.Group("/api/v1/auth")
	{
		public.POST("/login", h.Login)
		public.POST("/refresh", h.Refresh)
		public.POST("/logout", h.Logout)
	}

	private := router.Group("/api/v1/user")
	private.Use(authMiddleware)
	{
		private.GET("/me", h.GetProfile)
	}
}
