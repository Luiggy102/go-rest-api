package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Luiggy102/go-rest-ws/models"
	"github.com/Luiggy102/go-rest-ws/repository"
	"github.com/Luiggy102/go-rest-ws/server"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
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
		// check if user already exist by checking the input email
		usersExists, err := repository.UserEmailExists(context.Background(), signUpResquest.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if usersExists {
			http.Error(w,
				fmt.Sprintf("Email: %s already exist!\n", signUpResquest.Email),
				http.StatusBadRequest)
			return
			// fmt.Fprintf(w, "Email: %s already exist!\n", signUpResquest.Email)
		}
		// user email don't exist
		// generate user id
		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// generate password hash
		hassedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpResquest.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// create user struct
		u := models.User{
			Id:       id.String(),
			Email:    signUpResquest.Email,
			Password: string(hassedPassword),
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
