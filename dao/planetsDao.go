package dao

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/wallacebenevides/star-wars-api/db"
	"github.com/wallacebenevides/star-wars-api/models"
)

const (
	COLLECTION = "planets"
)

type PlanetsDAO interface {
	FindAll(ctx context.Context, filter interface{}) ([]models.Planet, error)
	Create(ctx context.Context, planets *models.Planet) error
	FindByID(cxt context.Context, filter interface{}) (*models.Planet, error)
	FindByName(cxt context.Context, filter interface{}) ([]models.Planet, error)
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

func (pd *planetsDAO) FindByID(ctx context.Context, filter interface{}) (*models.Planet, error) {
	var planet models.Planet
	if err := pd.db.Collection(COLLECTION).FindOne(ctx, filter).Decode(&planet); err != nil {
		return nil, err
	}
	return &planet, nil
}

func (pd *planetsDAO) FindByName(ctx context.Context, filter interface{}) ([]models.Planet, error) {
	var planets []models.Planet
	cursor, err := pd.db.Collection(COLLECTION).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var planet models.Planet
		if err := cursor.Decode(&planet); err != nil {
			return nil, err
		}
		planets = append(planets, planet)
	}
	return planets, err
}

func (pd *planetsDAO) Delete(ctx context.Context, filter interface{}) error {
	result, err := pd.db.Collection(COLLECTION).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("document not found")
	}
	log.Debug("Planet removed")
	return nil
}
