package apigateway

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Edilberto-Vazquez/weahter-services/src/services"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port string
}

type Server interface {
	Config() *Config
	Services() *services.Services
}

type Broker struct {
	config   *Config
	router   *gin.Engine
	services *services.Services
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Services() *services.Services {
	return b.services
}

func NewServer(ctx context.Context, config *Config, services *services.Services) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	broker := &Broker{
		config:   config,
		services: services,
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *gin.Engine)) {
	b.router = gin.Default()
	b.router.SetTrustedProxies([]string{"127.0.0.1"})
	binder(b, b.router)
	log.Println("[SERVER] starting server on port", b.config.Port)
	if err := b.router.Run(b.config.Port); err != nil {
		log.Println(fmt.Errorf("[SERVER] error starting server: %w", err))
	} else {
		log.Fatalf("[SERVER] server stopped")
	}
}
