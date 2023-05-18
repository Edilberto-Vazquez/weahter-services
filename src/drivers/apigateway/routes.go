package apigateway

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetRoutes(s Server, r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	stationsRoutes(s, v1)
}
