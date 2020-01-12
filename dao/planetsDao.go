package dao

import (
	"log"

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
