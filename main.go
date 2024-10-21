package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Luiggy102/go-rest-ws/handlers"
	"github.com/Luiggy102/go-rest-ws/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// load env vars
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	// server config load from env vars
	config := server.Config{
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		DBUrl:     os.Getenv("DATABASE_URL"),
	}

	// new server
	s, err := server.NewServer(context.Background(), &config)
	if err != nil {
		log.Fatal(err)
	}
	s.Start(BindRoutes)
}

// define endpoints for router start
func BindRoutes(s server.Server, r *mux.Router) {
	// handler func and http methods
	r.HandleFunc("/", handlers.HomeHandler(s)).
		Methods(http.MethodGet)
}
