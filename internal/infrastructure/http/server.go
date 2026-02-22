package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
	Port   string
}

func NewServer(port string, env string) *Server {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Generic health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return &Server{
		Router: router,
		Port:   port,
	}
}

func (s *Server) Start() error {
	log.Printf("Starting HTTP server on port %s", s.Port)
	return s.Router.Run(":" + s.Port)
}
