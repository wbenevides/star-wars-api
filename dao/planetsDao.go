package dao

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/db"
	"github.com/wallacebenevides/star-wars-api/models"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	COLLECTION = "planets"
)

type PlanetsDAO interface {
	FindAll() ([]models.Planet, error)
	Create(*models.Planet) error
	FindById(id string) (models.Planet, error)
	FindByName(name string) ([]models.Planet, error)
	Delete(*models.Planet) error
}

type planetsDAO struct {
	db db.DatabaseHelper
}

func NewPlanetsDao(db db.DatabaseHelper) PlanetsDAO {
	return &planetsDAO{db: db}
}

func (pd *planetsDAO) FindAll() ([]models.Planet, error) {
	var planets []models.Planet
	cursor, err := pd.db.Collection(COLLECTION).Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), &planets); err != nil {
		return nil, err
	}

	return planets, nil
}

func (pd *planetsDAO) Create(planet *models.Planet) error {
	_, err := pd.db.Collection(COLLECTION).InsertOne(context.TODO(), planet)
	if err != nil {
		log.WithField("name", planet.Name).Error("There was an error creating the planet")
		return err
	}
	log.WithField("name", planet.Name).Debug("Planet created")
	return nil
}

func (pd *planetsDAO) FindById(id string) (models.Planet, error) {
	var planet models.Planet
	filter := bson.D{{"_id", id}}
	err := pd.db.Collection(COLLECTION).FindOne(context.TODO(), filter).Decode(&planet)
	return planet, err
}

func (pd *planetsDAO) FindByName(name string) ([]models.Planet, error) {
	filter := bson.D{{"name", name}}
	var planets []models.Planet
	cursor, err := pd.db.Collection(COLLECTION).Find(context.TODO(), filter)
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

func (pd *planetsDAO) Delete(planet *models.Planet) error {
	withFild := log.WithField("name", planet.Name)
	_, err := pd.db.Collection(COLLECTION).DeleteOne(context.TODO(), planet)
	if err != nil {
		withFild.Error("The was an error removing the planet")
	}
	withFild.Debug("Planet removed")
	return nil
}
