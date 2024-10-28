package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Luiggy102/go-rest-ws/models"
	"github.com/Luiggy102/go-rest-ws/repository"
	"github.com/Luiggy102/go-rest-ws/server"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type SignUpLoginRequest struct {
	Email    string
	Password string
}
type SignUpResponse struct {
	Id    string
	Email string
}
type LogInResponse struct {
	Token string `json:"token"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signUpResquest := SignUpLoginRequest{}
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
func LogInHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// see what values the client sents
		var logInRequest = SignUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&logInRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// get user
		u, err := repository.GetUserByEmail(r.Context(), logInRequest.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if u == nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		}
		// compare password
		if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(logInRequest.Password)); err != nil {
			// invalid password
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// correct credentials
		// bengin to generate token response
		claims := models.AppClaims{
			UserId: u.Id,
			StandardClaims: jwt.StandardClaims{
				// 2 days for token expiration
				ExpiresAt: time.Now().Add(time.Hour * 24 * 2).Unix(),
			},
		}
		// generate new token
		// algorithm hs256 for making the token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// sign the token with the env password (added in the config)
		tokenStr, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// share the generated token
		logInResponse := LogInResponse{Token: tokenStr}
		err = json.NewEncoder(w).Encode(logInResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
