package services

import "github.com/Edilberto-Vazquez/weahter-services/src/usecases"

type Services struct {
	WeatherStationService *usecases.WeatherStation
}

func NewServices() *Services {
	return &Services{
		WeatherStationService: usecases.NewStation(usecases.WithMongoWeatherStationRepository()),
	}
}
