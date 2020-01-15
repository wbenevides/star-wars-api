package dao

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/db"
	"github.com/wallacebenevides/star-wars-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	COLLECTION = "planets"
)

type PlanetsDAO interface {
	FindAll() ([]models.Planet, error)
	Create(*models.Planet) (interface{}, error)
	FindById(id string) (*models.Planet, error)
	FindByName(name string) ([]models.Planet, error)
	Delete(id string) error
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

func (pd *planetsDAO) Create(planet *models.Planet) (interface{}, error) {
	insertedID, err := pd.db.Collection(COLLECTION).InsertOne(context.TODO(), planet)
	if err != nil {
		log.WithField("name", planet.Name).Error("There was an error creating the planet")
		return nil, err
	}
	log.WithField("name", planet.Name).Debug("Planet created")
	if oid, ok := insertedID.(primitive.ObjectID); ok {
		return map[string]interface{}{
			"id": oid.String(),
		}, nil
	}
	// Not objectid.ObjectID
	return map[string]interface{}{
		"id": "",
	}, nil

}

func (pd *planetsDAO) FindById(id string) (*models.Planet, error) {
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": idPrimitive}
	var planet models.Planet
	if err := pd.db.Collection(COLLECTION).FindOne(context.TODO(), filter).Decode(&planet); err != nil {
		return nil, err
	}
	return &planet, err
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

func (pd *planetsDAO) Delete(id string) error {
	withFild := log.WithField("id", id)
	// Declare a primitive ObjectID from a hexadecimal string
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": idPrimitive}
	result, err := pd.db.Collection(COLLECTION).DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("document not found")
	}
	withFild.Debug("Planet removed")
	return nil
}
