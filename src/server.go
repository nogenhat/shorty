package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Server tools ====================================================================================

func startServer() {
	bindAddress := config.Server.Address + ":" + config.Server.Port

	logger.Printf("Starting server on %s\n", bindAddress)

	if err := http.ListenAndServe(bindAddress, setupRouter()); err != nil {
		logger.Fatalf("Can't start server: %v", err)
	}
}

func setupRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/shorten", createUrlHandler).Methods("POST")
	router.HandleFunc("/{code}", redirectHandler).Methods("GET")
	router.HandleFunc("/expand/{code}", expandHandler).Methods("GET")
	router.HandleFunc("/statistics/{code}", statisticsHandler).Methods("GET")

	return
}

func serverError(rw http.ResponseWriter, err error, status int) {
	logger.Printf("Server error: %v", err)
	serverResponse(rw, "Internal server error", status)
}

func serverResponse(rw http.ResponseWriter, response string, status int) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(status)
	rw.Write([]byte(response))
}

// end of Server tools
