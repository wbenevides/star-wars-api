package resources

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/dao"
	"github.com/wallacebenevides/star-wars-api/models"
)

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

func GetAllPlanets(w http.ResponseWriter, r *http.Request) {
	log.Info("Fetching all planets")
	planets, err := dao.GetAllPlanets()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, planets)
}

func CreatePlanet(w http.ResponseWriter, r *http.Request) {
	log.Info("Creating a planet")
	defer r.Body.Close()
	var planet models.Planet
	if err := json.NewDecoder(r.Body).Decode(&planet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.CreatePlanet(planet); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, planet)

}

func GetPlanetByID(w http.ResponseWriter, r *http.Request) {
	log.Info("Fetching a planet by ID")
	params := mux.Vars(r)
	planet, err := dao.GetPlanetByID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Planet ID")
		return
	}
	respondWithJson(w, http.StatusOK, planet)
}

func FindPlanetByName(w http.ResponseWriter, r *http.Request) {
	log.Info("Finding a planet by name")
	name := r.URL.Query().Get("name")

	planet, err := dao.FindPlanetByName(name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, planet)
}

func DeletePlanet(w http.ResponseWriter, r *http.Request) {
	log.Info("Deleting a planet")
	defer r.Body.Close()
	var planet models.Planet
	if err := json.NewDecoder(r.Body).Decode(&planet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
	if err := dao.DeletePlanet(planet); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "sucess"})
}
