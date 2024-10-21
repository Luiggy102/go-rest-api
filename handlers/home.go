package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Luiggy102/go-rest-ws/server"
)

type HomeResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)                       // set the status code
		w.Header().Set("Content-Type", "application/json") // set the response header

		// prepare data response
		if err := json.NewEncoder(w).Encode(HomeResponse{
			Message: "Welcome to my REST Api",
			Status:  true,
		}); err != nil {
			log.Fatalln("Json encoding error", http.StatusInternalServerError)
		}
	}
}
