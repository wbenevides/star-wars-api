package resources

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/dao"
	"github.com/wallacebenevides/star-wars-api/db"
	"github.com/wallacebenevides/star-wars-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanetHandler struct {
	db dao.PlanetsDAO
}

func NewPlanetHandler(db db.DatabaseHelper) *PlanetHandler {
	dao := dao.NewPlanetsDao(db)
	return &PlanetHandler{dao}
}

func (h *PlanetHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Fetching all planets")
		planets, err := h.db.FindAll()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, planets)
	}
}

func (h *PlanetHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Creating a planet")
		defer r.Body.Close()
		var planet models.Planet
		if err := json.NewDecoder(r.Body).Decode(&planet); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		planet.ID = primitive.NewObjectID()
		if err := h.db.Create(&planet); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		result := createSuccessResult()
		respondWithJson(w, http.StatusCreated, result)
	}
}

func (h *PlanetHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Fetching a planet by ID")
		params := mux.Vars(r)
		planet, err := h.db.FindById(params["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid Planet ID")
			return
		}
		respondWithJson(w, http.StatusOK, planet)
	}
}

func (h *PlanetHandler) FindByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Finding a planet by name")
		name := r.URL.Query().Get("name")
		planet, err := h.db.FindByName(name)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, planet)
	}
}

func (h *PlanetHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Deleting a planet")
		defer r.Body.Close()
		var body struct{ ID string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		}

		if err := h.db.Delete(body.ID); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		result := createSuccessResult()
		respondWithJson(w, http.StatusOK, result)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	log.Error(msg)
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func createSuccessResult() map[string]string {
	return map[string]string{"result": "sucess"}
}
