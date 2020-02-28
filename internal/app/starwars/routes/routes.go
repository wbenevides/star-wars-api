package routes

import (
	"github.com/gorilla/mux"
	"github.com/wallacebenevides/star-wars-api/internal/pkg/db"
)

func Routes(router *mux.Router, db db.DatabaseHelper) {
	planetsRoutes(router, db)
}
