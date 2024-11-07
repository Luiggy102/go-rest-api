package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/Luiggy102/go-rest-ws/server"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Log(server server.Server) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler { // middleware
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapper := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(wrapper, r)
			log.Println(wrapper.statusCode, r.Method, r.URL.Path, time.Since(start))

		})
	}

}
