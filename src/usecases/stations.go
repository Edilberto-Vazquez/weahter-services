package usecases

import (
	"context"
	"log"

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

func (ws *WeatherStation) GetRecords(ctx context.Context, query models.FindRecords) ([]map[string]interface{}, error) {

	results, err := ws.stations.GetRecords(query, ctx)

	if err != nil {
		return nil, err
	}

	return results, nil
}
