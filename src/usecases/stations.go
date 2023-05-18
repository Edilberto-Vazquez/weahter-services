package usecases

import (
	"context"
	"log"
	"reflect"

	"github.com/Edilberto-Vazquez/weahter-services/src/config"
	"github.com/Edilberto-Vazquez/weahter-services/src/drivers/db"
	"github.com/Edilberto-Vazquez/weahter-services/src/models"
	"github.com/Edilberto-Vazquez/weahter-services/src/repository"
)

type WeatherStation struct {
	stations repository.StationRepository
}

type WeatherStationConfig func(ss *WeatherStation) error

func NewStation(cfgs ...WeatherStationConfig) *WeatherStation {
	ws := &WeatherStation{}
	for _, cfg := range cfgs {
		err := cfg(ws)
		if err != nil {
			log.Fatal(err)
		}
	}
	return ws
}

func WithMongoWeatherStationRepository() WeatherStationConfig {
	return func(ws *WeatherStation) error {
		dbConfig := models.DBConfig{URI: config.ENVS["DB_URI"]}
		repo, err := db.NewMongoDBConnection(dbConfig)
		if err != nil {
			return err
		}
		ws.stations = repo
		return nil
	}
}

func (ws *WeatherStation) Records(ctx context.Context, query models.FindRecords) ([]map[string]interface{}, error) {

	results, err := ws.stations.GetRecords(query, ctx)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (ws *WeatherStation) LineChart(ctx context.Context, query models.FindRecords) (*models.LineChart, error) {

	results, err := ws.stations.GetLineChart(query, ctx)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (ws *WeatherStation) BarChart(ctx context.Context, query models.FindRecords) ([]map[string]interface{}, error) {

	results, err := ws.stations.GetRecords(query, ctx)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (ws *WeatherStation) RadialChart(ctx context.Context, query models.FindRecords) (*models.RadialChart, error) {

	results, err := ws.stations.GetRadialChart(query, ctx)

	labels := make([]string, 0)
	series := []int64{results.Range1, results.Range2, results.Range3, results.Range4}

	radialType := reflect.TypeOf(results).Elem()
	for i := 0; i < radialType.NumField(); i++ {
		field := radialType.Field(i)
		labels = append(labels, field.Name)
	}

	if err != nil {
		return nil, err
	}

	return &models.RadialChart{Labels: labels, Series: series}, nil
}
