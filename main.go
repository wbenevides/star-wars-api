package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wallacebenevides/star-wars-api/resources"
)

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/planets", resources.GetPlanets).Methods(http.MethodGet)
	/* 	api.HandleFunc("/planets", resources.PostPlanet).Methods(http.MethodPost)

	   	api.HandleFunc("/planets/{id}", resources.GetPlanetById).Methods(http.MethodGet)
	   	api.HandleFunc("/planets/{id}", resources.DeletePlanet).Methods(http.MethodDelete)

	   	api.HandleFunc("/planets/findByName", resources.GetPlanetByName).Methods(http.MethodGet) */

	log.Fatal(http.ListenAndServe(":8080", r))
}
