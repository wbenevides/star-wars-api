package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wallacebenevides/star-wars-api/dao"
	"github.com/wallacebenevides/star-wars-api/db"
	"github.com/wallacebenevides/star-wars-api/resources"
)

func planetsRoutes(r *mux.Router, db db.DatabaseHelper) {
	dao := dao.NewPlanetsDao(db)
	handler := resources.NewPlanetHandler(dao)
	r.HandleFunc("/planets", handler.GetAll()).Methods(http.MethodGet)
	r.HandleFunc("/planets", handler.Create()).Methods(http.MethodPost)
	r.HandleFunc("/planets", handler.Delete()).Methods(http.MethodDelete)
	r.HandleFunc("/planets/findByName", handler.FindByName()).Methods(http.MethodGet)
	r.HandleFunc("/planets/{id}", handler.GetByID()).Methods(http.MethodGet)
}
