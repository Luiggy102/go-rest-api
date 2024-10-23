package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Luiggy102/go-rest-ws/models"
	"github.com/Luiggy102/go-rest-ws/repository"
	"github.com/Luiggy102/go-rest-ws/server"
	"github.com/segmentio/ksuid"
)

type SignUpRequest struct {
	Email    string
	Password string
}
type SignUpResponse struct {
	Id    string
	Email string
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signUpResquest := SignUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&signUpResquest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// generate user id
		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// create user struct
		u := models.User{
			Id:       id.String(),
			Email:    signUpResquest.Email,
			Password: signUpResquest.Password,
		}
		// create a new user
		// using the db repository
		err = repository.InsertUser(r.Context(), &u) // add the context of the Request
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// -- user created and added successfully -- //
		// create the header
		w.Header().Set("Content-Type", "application/json") // json response
		// send the response
		signUpResponse := SignUpResponse{
			Id:    u.Id,
			Email: u.Email,
		}
		// sent response as a json
		err = json.NewEncoder(w).Encode(&signUpResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
