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
	COLLECTION              = "planets"
	NOT_FOUND_ERROR_MESSAGE = "document not found"
)

type PlanetsDAO interface {
	FindAll(ctx context.Context, filter interface{}) ([]models.Planet, error)
	Create(ctx context.Context, planets *models.Planet) error
	FindOne(cxt context.Context, filter interface{}) (*models.Planet, error)
	FindByName(cxt context.Context, name string) ([]models.Planet, error)
	Delete(cxt context.Context, filter interface{}) error
}

type planetsDAO struct {
	db db.DatabaseHelper
}

func NewPlanetsDao(db db.DatabaseHelper) PlanetsDAO {
	return &planetsDAO{db: db}
}

func (pd *planetsDAO) FindAll(ctx context.Context, filter interface{}) ([]models.Planet, error) {
	var planets []models.Planet
	cursor, err := pd.db.Collection(COLLECTION).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &planets); err != nil {
		return nil, err
	}
	return planets, nil
}

func (pd *planetsDAO) Create(ctx context.Context, planet *models.Planet) error {
	_, err := pd.db.Collection(COLLECTION).InsertOne(ctx, planet)
	if err != nil {
		log.WithField("name", planet.Name).Error("There was an error creating the planet")
		return err
	}
	log.WithField("name", planet.Name).Debug("Planet created")
	return nil
}

func (pd *planetsDAO) FindOne(ctx context.Context, filter interface{}) (*models.Planet, error) {
	var planet models.Planet
	if err := pd.db.Collection(COLLECTION).FindOne(ctx, filter).Decode(&planet); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(NOT_FOUND_ERROR_MESSAGE)
		}
		return nil, err
	}
	return &planet, nil
}

func (pd *planetsDAO) FindByName(ctx context.Context, name string) ([]models.Planet, error) {
	filter := bson.D{{"name", primitive.Regex{Pattern: name, Options: "i"}}}
	return pd.FindAll(ctx, filter)
}

func (pd *planetsDAO) Delete(ctx context.Context, filter interface{}) error {
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
