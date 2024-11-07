package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Luiggy102/go-rest-ws/handlers"
	"github.com/Luiggy102/go-rest-ws/middleware"
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
	// use the middleware
	r.Use(
		middleware.Log(s),
		middleware.CheckAuth(s),
	)
	// handler func and http methods
	// --- no token auth --- //
	r.HandleFunc("/", handlers.HomeHandler(s)).
		Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).
		Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LogInHandler(s)).
		Methods(http.MethodPost)

		// token auth
	r.HandleFunc("/me", handlers.MeHandler(s)).
		Methods(http.MethodGet)
	r.HandleFunc("/posts", handlers.InsertPostHandler(s)).
		Methods(http.MethodPost)
}
