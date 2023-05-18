package apigateway

import (
	"net/http"
	"strings"
	"time"

	"github.com/Edilberto-Vazquez/weahter-services/src/models"
	"github.com/Edilberto-Vazquez/weahter-services/src/services"

	"github.com/gin-gonic/gin"
)

func Records(svcs *services.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		collection := ctx.Param("collection")
		db := ctx.Param("database")

		weatherStationQ := ctx.Query("weatherStation")
		weatherStationF := strings.Split(weatherStationQ, ",")
		fieldMonitorQ := ctx.Query("fieldMonitor")
		fieldMonitorF := strings.Split(fieldMonitorQ, ",")

		dateStart, _ := time.Parse(time.RFC3339, "2019-04-00T00:00:00Z")
		dateEnd, _ := time.Parse(time.RFC3339, "2019-04-30T11:59:59Z")

		query := models.FindRecords{
			DB:         db,
			Collection: collection,
			DateStart:  dateStart,
			DateEnd:    dateEnd,
			Fields:     weatherStationF,
		}

		records, err := svcs.WeatherStationService.GetRecords(ctx, query)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ocurrió un error al obtener los registros del clima",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"db":             db,
			"collection":     collection,
			"weatherStation": weatherStationF,
			"fieldMonitor":   fieldMonitorF,
			"records":        records,
		})
	}
}
