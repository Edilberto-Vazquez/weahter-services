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

		fieldsStr := ctx.Query("fields")
		fieldsList := strings.Split(fieldsStr, ",")

		dateStart, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
		dateEnd, _ := time.Parse(time.RFC3339, "2019-10-02T11:59:59Z")

		query := models.FindRecords{
			DB:         db,
			Collection: collection,
			DateStart:  dateStart,
			DateEnd:    dateEnd,
			Fields:     fieldsList,
		}

		records, err := svcs.WeatherStationService.Records(ctx, query)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"records": records,
		})
	}
}

func LineChart(svcs *services.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		collection := ctx.Param("collection")
		db := ctx.Param("database")

		fieldsStr := ctx.Query("fields")
		fieldsList := strings.Split(fieldsStr, ",")

		datesStr := ctx.Query("dates")
		datesList := strings.Split(datesStr, ",")

		dateStart, _ := time.Parse(time.RFC3339, datesList[0]+"T00:00:00Z")
		dateEnd, _ := time.Parse(time.RFC3339, datesList[1]+"T23:59:59Z")

		query := models.FindRecords{
			DB:         db,
			Collection: collection,
			DateStart:  dateStart,
			DateEnd:    dateEnd,
			Fields:     fieldsList,
		}

		lineChartData, err := svcs.WeatherStationService.LineChart(ctx, query)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if lineChartData == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"data": []interface{}{},
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": lineChartData.Data,
		})
	}
}
