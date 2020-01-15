package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/config"
	"github.com/wallacebenevides/star-wars-api/db"
	"github.com/wallacebenevides/star-wars-api/resources"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// initialize db config
	config := config.Config{}
	config.Read()
	client, err := db.NewClient(&config.Database)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Connected to MongoDB!")

	database := db.NewDatabase(&config.Database, client)
	planetHandler := resources.NewPlanetHandler(database)

	r := mux.NewRouter()
	log.Info("star wars planets api is listening on port ", config.Server.Port)
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "api v1")
	})

	api.Use(loggingMiddleware)
	api.HandleFunc("/planets", planetHandler.GetAll()).Methods(http.MethodGet)
	api.HandleFunc("/planets", planetHandler.Create()).Methods(http.MethodPost)
	api.HandleFunc("/planets", planetHandler.Delete()).Methods(http.MethodDelete)
	api.HandleFunc("/planets/findByName", planetHandler.FindByName()).Methods(http.MethodGet)
	api.HandleFunc("/planets/{id}", planetHandler.GetByID()).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":"+config.Server.Port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
