package dao

import (
	"log"

	"github.com/wallacebenevides/star-wars-api/models"
	. "github.com/wallacebenevides/star-wars-api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

type PlanetsDAO struct {
	Hosts    string
	Database string
}

const (
	COLLECTION = "planets"
)

func (m *PlanetsDAO) Connect() {

	session, err := mgo.Dial(m.Hosts)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func GetAllPlanets() ([]Planet, error) {
	var movies []Planet
	err := db.C(COLLECTION).Find(bson.M{}).All(&movies)
	return movies, err
}

func CreatePlanet(planet models.Planet) error {
	err := db.C(COLLECTION).Insert(&planet)
	return err
}

func GetPlanetByID(id string) (Planet, error) {
	var planet models.Planet
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&planet)
	return planet, err
}
