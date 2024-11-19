package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Luiggy102/go-rest-ws/models"
	"github.com/Luiggy102/go-rest-ws/repository"
	"github.com/Luiggy102/go-rest-ws/server"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type PostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

type PostUpdateResponse struct {
	Message string `json:"message"`
}

// create
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

			// send a websocket message
			wsmessage := models.WebsocketMessage{
				Type:    "Post_created",
				Payload: p,
			}
			s.Hub().Broadcast(wsmessage, nil)

			// prepare response
			var response = PostResponse{
				Id:          p.Id,
				PostContent: p.PostContent,
			}

			// send the response
			w.Header().Set("Content-Type", "application/json")
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

// read
func GetPostByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		postId := params["id"]

		// call the db
		p, err := repository.GetPostById(r.Context(), postId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// send the response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// update
func UpdatePostHandler(s server.Server) http.HandlerFunc {
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
			params := mux.Vars(r)
			// decode de update
			var request PostRequest
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// get the user id and the post id
			postId := params["id"]
			userId := claims.UserId

			// struct for the db
			p := models.Post{
				Id:          postId,
				UserId:      userId,
				PostContent: request.PostContent,
			}

			// update the db
			err = repository.UpdatePost(r.Context(), &p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// send the response
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(&PostUpdateResponse{
				Message: "Post Updated",
			})
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

// delete
func DelelePostHandler(s server.Server) http.HandlerFunc {
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
			params := mux.Vars(r)
			// get the user id and the post id
			postId := params["id"]
			userId := claims.UserId

			err = repository.DeletePost(r.Context(), postId, userId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// send the response
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(&PostUpdateResponse{
				Message: "Post Deleted",
			})
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

// list the posts
func ListPosts(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		// get the query param
		pageStr := r.URL.Query().Get("page")
		// default value
		var page = uint64(1)

		if pageStr != "" {
			// not empty = now value
			page, err = strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				// an error like `/posts?page=foo`
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}

		// call the db for the structs
		posts, err := repository.ListPosts(r.Context(), page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// send the response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
