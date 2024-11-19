package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Luiggy102/go-rest-ws/handlers"
	"github.com/Luiggy102/go-rest-ws/middleware"
	"github.com/Luiggy102/go-rest-ws/server"
	"github.com/Luiggy102/go-rest-ws/websocket"
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
	// add the hub for websockets
	hub := websocket.NewHub()
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
		Methods(http.MethodPost) // creatre
	r.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).
		Methods(http.MethodGet) // read
	r.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).
		Methods(http.MethodPut) // update
	r.HandleFunc("/posts/{id}", handlers.DelelePostHandler(s)).
		Methods(http.MethodDelete) // delete

		// pagination
	r.HandleFunc("/posts", handlers.ListPosts(s)).
		Methods(http.MethodGet)
	// websockets
	r.HandleFunc("/ws", hub.HandleWebSocket)
}
