package main

import (
	"net/http"

	"github.com/born2ngopi/eel/pkg/ell"
)

func main() {

	ell.Init("token", "http://localhost:8080")
	// watch new token
	opt := ell.WatchOption{
		Interval: "@hourly",
		Driver:   ell.RABBITMQ_DRIVER,
		Username: "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     "5672",
	}
	ell.Watch(opt, []string{"service-b"})

}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := ell.GetToken("service-b")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if r.Header.Get("Authorization") != token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
