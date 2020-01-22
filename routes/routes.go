package routes

import "github.com/gorilla/mux"

import "github.com/wallacebenevides/star-wars-api/db"

func Routes(router *mux.Router, db db.DatabaseHelper) {
	planetsRoutes(router, db)
}
