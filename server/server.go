package server

import "net/http"

func Start() {

	mux := http.NewServeMux()

	InitRoute(mux)

	http.ListenAndServe(":8080", mux)

}
