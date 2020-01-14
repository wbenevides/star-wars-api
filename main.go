package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/config"
	"github.com/wallacebenevides/star-wars-api/dao"
	"github.com/wallacebenevides/star-wars-api/resources"
)

var configuration config.Config

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// initialize db config
	configuration = config.Config{}
	configuration.Read()
	connectionUri := configuration.Database.ConnectionUri
	database := configuration.Database.Database

	log.Info(connectionUri, database)

	var planetDao = dao.PlanetsDAO{
		connectionUri,
		database,
	}
	planetDao.Connect()
}

func main() {
	port := configuration.Server.Port

	r := mux.NewRouter()
	log.Info("star wars planets api is listening on port ", port)
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "api v1")
	})
	api.Use(loggingMiddleware)

	api.HandleFunc("/planets", resources.GetAllPlanets).Methods(http.MethodGet)
	api.HandleFunc("/planets", resources.CreatePlanet).Methods(http.MethodPost)
	api.HandleFunc("/planets", resources.DeletePlanet).Methods(http.MethodDelete)
	api.HandleFunc("/planets/findByName", resources.FindPlanetByName).Methods(http.MethodGet)
	api.HandleFunc("/planets/{id}", resources.GetPlanetByID).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
