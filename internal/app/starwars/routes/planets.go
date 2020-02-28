package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wallacebenevides/star-wars-api/internal/app/starwars/dao"
	"github.com/wallacebenevides/star-wars-api/internal/app/starwars/resources"
	"github.com/wallacebenevides/star-wars-api/internal/pkg/db"
)

func planetsRoutes(r *mux.Router, db db.DatabaseHelper) {
	dao := dao.NewPlanetsDao(db)
	handler := resources.NewPlanetHandler(dao)
	routes := handler.Routes()
	r.HandleFunc(routes.PLANETS_PATH, handler.GetAll()).Methods(http.MethodGet)
	r.HandleFunc(routes.PLANETS_PATH, handler.Create()).Methods(http.MethodPost)
	r.HandleFunc(routes.PLANETS_FIND_BY_NAME, handler.FindByName()).Methods(http.MethodGet)
	r.HandleFunc(routes.PLANETS_ID, handler.Delete()).Methods(http.MethodDelete)
	r.HandleFunc(routes.PLANETS_ID, handler.GetByID()).Methods(http.MethodGet)
}
