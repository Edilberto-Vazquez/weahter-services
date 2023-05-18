package main

import (
	"context"
	"log"

	"github.com/Edilberto-Vazquez/weahter-services/src/config"
	"github.com/Edilberto-Vazquez/weahter-services/src/drivers/apigateway"
	"github.com/Edilberto-Vazquez/weahter-services/src/services"
)

func main() {
	config.SetEnvironment()
	log.Println("PORT: " + config.ENVS["PORT"])

	s, err := apigateway.NewServer(
		context.Background(),
		&apigateway.Config{Port: config.ENVS["PORT"]},
		services.NewServices(),
	)
	if err != nil {
		log.Fatal(err)
	}
	s.Start(apigateway.GetRoutes)

}
