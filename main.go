package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wallacebenevides/star-wars-api/dao"
	"github.com/wallacebenevides/star-wars-api/resources"
)

const (
	hosts    = "localhost"
	database = "db"
)

var planetDao = dao.PlanetsDAO{}

func init() {
	planetDao.Hosts = hosts
	planetDao.Database = database
	planetDao.Connect()
}

func main() {
	r := mux.NewRouter()
	log.Println("star wars planets api")
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "api v1")
	})

	api.HandleFunc("/planets", resources.GetPlanets).Methods(http.MethodGet)
	/* api.HandleFunc("/planets", resources.PostPlanet).Methods(http.MethodPost)
	api.HandleFunc("/planets/{id}", resources.GetPlanetById).Methods(http.MethodGet)
	api.HandleFunc("/planets/{id}", resources.DeletePlanet).Methods(http.MethodDelete)
	api.HandleFunc("/planets/findByName", resources.GetPlanetByName).Methods(http.MethodGet) */

	log.Fatal(http.ListenAndServe(":8080", r))
}
