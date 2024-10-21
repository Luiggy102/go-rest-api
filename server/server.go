package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port      string
	JWTSecret string
	DBUrl     string
}
type Server interface { // to satisfy the interface(for broker) must return a Config
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

// satisfy interface(broker)
func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	switch { // check config values
	case config.Port == "":
		return nil, errors.New("Port is required")
	case config.DBUrl == "":
		return nil, errors.New("Database is required")
	case config.JWTSecret == "":
		return nil, errors.New("JWT is required")
	}
	// create server
	return &Broker{
		config: config,
		router: mux.NewRouter(),
	}, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter() // handler
	binder(b, b.router)
	log.Println("Starting server on port:", b.config.Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatalln("Failed to start server:", err)
	}
}
