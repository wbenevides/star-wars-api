package dao

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

type PlanetsDAO struct {
	ConnectionUri string
	Database      string
}

const (
	COLLECTION = "planets"
)

func (m *PlanetsDAO) Connect() {
	clientOptions := options.Client().ApplyURI(m.ConnectionUri)

	log.Info("initializing a session with db", m.Database)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to MongoDB!")
	db = client.Database(m.Database)
}

func GetAllPlanets() ([]models.Planet, error) {
	var planets []models.Planet
	cursor, err := db.Collection(COLLECTION).Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), &planets); err != nil {
		return nil, err
	}

	return planets, nil
}

func CreatePlanet(planet models.Planet) error {
	_, err := db.Collection(COLLECTION).InsertOne(context.TODO(), &planet)
	if err != nil {
		log.WithField("name", planet.Name).Error("There was an error creating the planet")
		return err
	}
	log.WithField("name", planet.Name).Debug("Planet created")
	return nil
}

func GetPlanetByID(id string) (models.Planet, error) {
	var planet models.Planet
	filter := bson.D{{"_id", id}}
	err := db.Collection(COLLECTION).FindOne(context.TODO(), filter).Decode(&planet)
	return planet, err
}

func FindPlanetByName(name string) ([]models.Planet, error) {
	filter := bson.D{{"name", name}}
	var planets []models.Planet
	cursor, err := db.Collection(COLLECTION).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var planet models.Planet
		if err := cursor.Decode(&planet); err != nil {
			return nil, err
		}
		planets = append(planets, planet)
	}
	return planets, err
}

func DeletePlanet(planet models.Planet) error {
	withFild := log.WithField("name", planet.Name)
	_, err := db.Collection(COLLECTION).DeleteOne(context.TODO(), &planet)
	if err != nil {
		withFild.Error("The was an error removing the planet")
	}
	withFild.Debug("Planet removed")
	return nil
}
