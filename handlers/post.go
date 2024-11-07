package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Luiggy102/go-rest-ws/models"
	"github.com/Luiggy102/go-rest-ws/repository"
	"github.com/Luiggy102/go-rest-ws/server"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
)

type PostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(
			tokenStr,
			&models.AppClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			// correct token pass
			id, err := ksuid.NewRandom() // (new) post id
			userId := claims.UserId      //       user id
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// json post content
			var request = PostRequest{}
			err = json.NewDecoder(r.Body).Decode(&request) // post content
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// struct for the db
			p := models.Post{
				Id:          id.String(),
				PostContent: request.PostContent,
				UserId:      userId,
			}

			// add to the db
			err = repository.InsertPost(r.Context(), &p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// prepare response
			var response = PostResponse{
				Id:          p.Id,
				PostContent: p.PostContent,
			}

			// send the response
			err = json.NewEncoder(w).Encode(&response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// incorrect token or error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
