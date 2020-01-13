package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/dao"
	"github.com/wallacebenevides/star-wars-api/resources"
)

const (
	hosts    = "mongodb"
	database = "star_wars_db"
)

const (
	port = ":8080"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	var planetDao = dao.PlanetsDAO{}
	planetDao.Hosts = hosts
	planetDao.Database = database
	planetDao.Connect()

}

func main() {

	r := mux.NewRouter()
	log.Info("star wars planets api is listening on port ", port)
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "api v1")
	})

	api.HandleFunc("/planets", resources.GetAllPlanets).Methods(http.MethodGet)
	api.HandleFunc("/planets", resources.CreatePlanet).Methods(http.MethodPost)
	api.HandleFunc("/planets", resources.DeletePlanet).Methods(http.MethodDelete)
	api.HandleFunc("/planets/findByName", resources.FindPlanetByName).Methods(http.MethodGet)
	api.HandleFunc("/planets/{id}", resources.GetPlanetByID).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(port, r))
}
