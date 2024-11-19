package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Luiggy102/go-rest-ws/database"
	"github.com/Luiggy102/go-rest-ws/repository"
	"github.com/Luiggy102/go-rest-ws/websocket"
	"github.com/gorilla/mux"
)

type Config struct {
	Port      string
	JWTSecret string
	DBUrl     string
}
type Server interface { // to satisfy the interface(for broker) must return a Config
	Config() *Config
	Hub() *websocket.Hub // hub for websocket connection
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub // ws
}

// satisfy interface(broker)
func (b *Broker) Config() *Config {
	return b.config
}
func (b *Broker) Hub() *websocket.Hub {
	return b.hub // ws
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
		hub:    websocket.NewHub(), // added the hub to the broker
	}, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter() // handler
	binder(b, b.router)        // bind routers
	// init db
	// start the postgres repo
	repo, err := database.NewPostgresRepo(b.config.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	b.hub.Run() // start the ws connection
	repository.SetRepo(repo)
	// show server status
	log.Println("Starting server on port", b.config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", b.config.Port), b.router); err != nil {
		log.Fatalln("Failed to start server:", err)
	}
}
