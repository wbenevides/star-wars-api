package resources

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/internal/app/starwars/dao"
	"github.com/wallacebenevides/star-wars-api/internal/app/starwars/models"
)

type PlanetHandler struct {
	db dao.PlanetsDAO
}

const (
	INVALID_REQUEST_PAYLOAD_ERROR_MESSAGE = "Invalid request payload"
	INTERNAL_SERVER_ERROR_MESSAGE         = "Operation could not be performed"
)

type routes struct {
	PLANETS_PATH         string
	PLANETS_ID           string
	PLANETS_FIND_BY_NAME string
}

func NewPlanetHandler(dao dao.PlanetsDAO) *PlanetHandler {
	return &PlanetHandler{dao}
}

func (h PlanetHandler) Routes() routes {
	return routes{
		PLANETS_PATH:         "/planets",
		PLANETS_ID:           "/planets/{id}",
		PLANETS_FIND_BY_NAME: "/planets/findByName",
	}
}

func (h *PlanetHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Finding all planets")
		planets, err := h.db.FindAll(context.TODO())
		if err != nil {
			errorHandler(w, err)
			return
		}
		respondWithJson(w, http.StatusOK, planets)
	}
}

func (h *PlanetHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var planet models.Planet
		if err := json.NewDecoder(r.Body).Decode(&planet); err != nil {
			log.Debug(err.Error(), planet)
			errorHandler(w, errors.New(INVALID_REQUEST_PAYLOAD_ERROR_MESSAGE))
			return
		}
		log.Info("Creating a planet")
		id, err := h.db.Create(context.TODO(), &planet)
		if err != nil {
			errorHandler(w, err)
			return
		}
		planet.ID = id
		location := r.URL.String() + "/" + id
		w.Header().Add("Location", location)
		respondWithJson(w, http.StatusCreated, planet)
	}
}

func (h *PlanetHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Info("Finding a planet by ID")
		planet, err := h.db.FindByID(context.TODO(), params["id"])
		if err != nil {
			errorHandler(w, err)
			return
		}
		respondWithJson(w, http.StatusOK, planet)
	}
}

func (h *PlanetHandler) FindByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		log.Info("Finding planets by name")
		planets, err := h.db.FindByName(context.TODO(), name)
		if err != nil {
			errorHandler(w, err)
			return
		}
		if len(planets) == 0 {
			errorHandler(w, errors.New(dao.NOT_FOUND_ERROR_MESSAGE))
			return
		}
		respondWithJson(w, http.StatusOK, planets)
	}
}

func (h *PlanetHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		log.Info("params", params)
		log.Info("Deleting a planet", params["id"])
		if err := h.db.Delete(context.TODO(), params["id"]); err != nil {
			errorHandler(w, err)
			return
		}
		result := createSuccessResult()
		respondWithJson(w, http.StatusOK, result)
	}
}

func errorHandler(w http.ResponseWriter, err error) {
	switch err.Error() {
	case dao.INVALID_ID_ERROR_MESSAGE,
		INVALID_REQUEST_PAYLOAD_ERROR_MESSAGE:
		respondWithError(w, http.StatusBadRequest, err.Error())
	case dao.NOT_FOUND_ERROR_MESSAGE:
		respondWithError(w, http.StatusNotFound, err.Error())
	default:
		log.Error(err)
		respondWithError(w, http.StatusInternalServerError, INTERNAL_SERVER_ERROR_MESSAGE)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func createSuccessResult() map[string]string {
	return map[string]string{"result": "success"}
}
