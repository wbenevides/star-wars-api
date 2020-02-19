package dao

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/db"
	"github.com/wallacebenevides/star-wars-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	COLLECTION = "planets"
)

const (
	NOT_FOUND_ERROR_MESSAGE  = "document not found"
	INVALID_ID_ERROR_MESSAGE = "Invalid Planet ID"
)

type PlanetsDAO interface {
	FindAll(ctx context.Context) ([]models.Planet, error)
	Create(ctx context.Context, planets *models.Planet) (interface{}, error)
	FindByID(cxt context.Context, id string) (models.Planet, error)
	FindByName(cxt context.Context, name string) ([]models.Planet, error)
	Delete(cxt context.Context, id string) error
}

type planetsDAO struct {
	db db.DatabaseHelper
}

func NewPlanetsDao(db db.DatabaseHelper) PlanetsDAO {
	return &planetsDAO{db: db}
}

func (pd *planetsDAO) FindAll(ctx context.Context) ([]models.Planet, error) {
	filter := bson.D{{}}
	return pd.find(ctx, filter)
}

func (pd *planetsDAO) Create(ctx context.Context, planet *models.Planet) (interface{}, error) {
	result, err := pd.db.Collection(COLLECTION).InsertOne(ctx, planet)
	if err != nil {
		log.WithField("name", planet.Name).Error("There was an error creating the planet::", err.Error())
		return nil, err
	}
	log.WithField("name", planet.Name).Info("Planet created:", result.InsertedID.(primitive.ObjectID).String())
	return result.InsertedID, nil
}

func (pd *planetsDAO) FindByID(ctx context.Context, id string) (models.Planet, error) {
	objectID, err := createObjectIDFromHex(id)
	var planet models.Planet
	if err != nil {
		log.WithField("id", id).Error("There was an error find the planet by id")
		return planet, err
	}
	filter := bson.M{"_id": objectID}
	planet, err = pd.findOne(ctx, filter)

	if err != nil {
		log.WithField("id", id).Error("There was an error find the planet by id")
		return planet, err
	}
	return planet, nil
}

func (pd *planetsDAO) FindByName(ctx context.Context, name string) ([]models.Planet, error) {
	filter := bson.D{{"name", primitive.Regex{Pattern: name, Options: "i"}}}
	return pd.find(ctx, filter)
}

func (pd *planetsDAO) Delete(ctx context.Context, id string) error {
	objectID, err := createObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectID}

	result, err := pd.db.Collection(COLLECTION).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New(NOT_FOUND_ERROR_MESSAGE)
	}
	log.Debug("Planet removed")
	return nil
}

func (pd *planetsDAO) find(ctx context.Context, filter interface{}) ([]models.Planet, error) {
	var planets []models.Planet
	cursor, err := pd.db.Collection(COLLECTION).Find(ctx, filter)
	if err != nil {
		log.WithField("filter", filter).Error("There was an error finding the planets::", err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &planets); err != nil {
		log.Error(err)
		return nil, err
	}
	return planets, nil
}

func (pd *planetsDAO) findOne(ctx context.Context, filter interface{}) (models.Planet, error) {
	var planet models.Planet
	if err := pd.db.Collection(COLLECTION).FindOne(ctx, filter).Decode(&planet); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Error(err)
			return planet, errors.New(NOT_FOUND_ERROR_MESSAGE)
		}
		return planet, err
	}
	return planet, nil
}

func createObjectIDFromHex(id string) (primitive.ObjectID, error) {
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return idPrimitive, errors.New(INVALID_ID_ERROR_MESSAGE)
	}
	return idPrimitive, nil
}
