package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/internal/app/starwars/routes"
	"github.com/wallacebenevides/star-wars-api/internal/pkg/config"
	"github.com/wallacebenevides/star-wars-api/internal/pkg/db"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	config := config.Config{}
	config.Read()
	log.Info(config)
	database := initializeDB(config)

	r := mux.NewRouter()
	api := newRouterAPI(r)

	api.Use(loggingMiddleware)
	routes.Routes(api, database)

	log.Info("star wars planets api is listening on port ", config.Server.Port)
	log.Fatal(http.ListenAndServe(":"+config.Server.Port, r))
}

func newRouterAPI(r *mux.Router) *mux.Router {
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "api v1")
	})
	return api
}

func initializeDB(config config.Config) db.DatabaseHelper {
	// initialize db config
	client, err := db.NewClient(&config.Database)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Connected to MongoDB!")
	return db.NewDatabase(&config.Database, client)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
