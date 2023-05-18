package apigateway

import (
	"github.com/gin-gonic/gin"
)

func stationsRoutes(s Server, rg *gin.RouterGroup) {
	services := s.Services()
	stations := rg.Group("/weather")
	stations.GET("/stations/:database/:collection/", Records(services))
}
