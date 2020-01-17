package dao

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wallacebenevides/star-wars-api/db"
	"github.com/wallacebenevides/star-wars-api/mocks"
	"github.com/wallacebenevides/star-wars-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_planetsDAO_FindAll(t *testing.T) {

	// Define variables for interfaces
	var dbHelper db.DatabaseHelper
	var collectionHelper db.CollectionHelper
	var cursorHelperErr db.CursorHelper
	var cursorHelperCorrect db.CursorHelper

	//Set interfaces implementation to mocked structures
	dbHelper = &mocks.DatabaseHelper{}
	collectionHelper = &mocks.CollectionHelper{}
	cursorHelperErr = &mocks.CursorHelper{}
	cursorHelperCorrect = &mocks.CursorHelper{}

	cursorHelperErr.(*mocks.CursorHelper).
		On("All", context.Background(), mock.AnythingOfType("*[]models.Planet")).
		Return(nil, errors.New("mocked-error"))

	cursorHelperCorrect.(*mocks.CursorHelper).
		On("All", context.Background(), mock.AnythingOfType("*[]models.Planet")).
		Return(nil, nil).Run(func(args mock.Arguments) {
		planets := args.Get(0).([]models.Planet)
		for _, planet := range planets {
			planet.Name = "mocked-planet"
		}
	})

	collectionHelper.(*mocks.CollectionHelper).
		On("Find", context.Background(), bson.M{"error": true}).
		Return(cursorHelperErr)

	collectionHelper.(*mocks.CollectionHelper).
		On("Find", context.Background(), bson.M{"error": false}).
		Return(cursorHelperCorrect)

	dbHelper.(*mocks.DatabaseHelper).
		On("Collection", "planets").
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)
	planets, err := planetDao.FindAll(context.Background(), bson.M{"error": true})
	assert.Empty(t, planets)
	assert.EqualError(t, err, "mocked-error")

	planets, err = planetDao.FindAll(context.Background(), bson.M{"error": false})
	expected := &[]models.Planet{
		{Name: "mocked-planet"},
	}
	assert.Equal(t, expected, planets)
	assert.NoError(t, err)
}

func Test_planetsDAO_FindOne(t *testing.T) {

	// Define variables for interfaces
	var dbHelper db.DatabaseHelper
	var collectionHelper db.CollectionHelper
	var srHelperErr db.SingleResultHelper
	var srHelperCorrect db.SingleResultHelper

	//Set interfaces implementation to mocked structures
	dbHelper = &mocks.DatabaseHelper{}
	collectionHelper = &mocks.CollectionHelper{}
	srHelperErr = &mocks.SingleResultHelper{}
	srHelperCorrect = &mocks.SingleResultHelper{}

	srHelperErr.(*mocks.SingleResultHelper).
		On("Decode", mock.AnythingOfType("*models.Planet")).
		Return(errors.New("mocked-error"))

	srHelperCorrect.(*mocks.SingleResultHelper).
		On("Decode", mock.AnythingOfType("*models.Planet")).
		Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Planet)
		arg.Name = "mocked-planet"
	})

	collectionHelper.(*mocks.CollectionHelper).
		On("FindOne", context.Background(), bson.M{"error": true}).
		Return(srHelperErr)

	collectionHelper.(*mocks.CollectionHelper).
		On("FindOne", context.Background(), bson.M{"error": false}).
		Return(srHelperCorrect)

	dbHelper.(*mocks.DatabaseHelper).
		On("Collection", "planets").
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	planet, err := planetDao.FindOne(context.Background(), bson.M{"error": true})
	assert.Empty(t, planet)
	assert.EqualError(t, err, "mocked-error")

	planet, err = planetDao.FindOne(context.Background(), bson.M{"error": false})
	assert.Equal(t, &models.Planet{Name: "mocked-planet"}, planet)
	assert.NoError(t, err)
}

func Test_planetsDAO_Create(t *testing.T) {
	// Define variables for interfaces
	var dbHelper db.DatabaseHelper
	var collectionHelper db.CollectionHelper

	//Set interfaces implementation to mocked structures
	dbHelper = &mocks.DatabaseHelper{}
	collectionHelper = &mocks.CollectionHelper{}

	collectionHelper.(*mocks.CollectionHelper).
		On("InsertOne", context.Background(), &models.Planet{Name: "mocked-planet-error"}).
		Return(nil, errors.New("mocked-error"))

	collectionHelper.(*mocks.CollectionHelper).
		On("InsertOne", context.Background(), &models.Planet{Name: "mocked-planet-correct"}).
		Return(nil, nil)

	dbHelper.(*mocks.DatabaseHelper).
		On("Collection", "planets").
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	err := planetDao.Create(context.Background(), &models.Planet{Name: "mocked-planet-error"})
	assert.EqualError(t, err, "mocked-error")

	err = planetDao.Create(context.Background(), &models.Planet{Name: "mocked-planet-correct"})
	assert.NoError(t, err)
}

func Test_planetsDAO_Delete(t *testing.T) {
	// Define variables for interfaces
	var dbHelper db.DatabaseHelper
	var collectionHelper db.CollectionHelper

	//Set interfaces implementation to mocked structures
	dbHelper = &mocks.DatabaseHelper{}
	collectionHelper = &mocks.CollectionHelper{}
	deleteResultCorrect := mongo.DeleteResult{1}

	collectionHelper.(*mocks.CollectionHelper).
		On("DeleteOne", context.Background(), bson.M{"db-error": true}).
		Return(nil, errors.New("mocked-db-error"))

	collectionHelper.(*mocks.CollectionHelper).
		On("DeleteOne", context.Background(), bson.M{"notFound-error": true}).
		Return(nil, errors.New("document not found"))

	collectionHelper.(*mocks.CollectionHelper).
		On("DeleteOne", context.Background(), bson.M{"error": false}).
		Return(&deleteResultCorrect, nil)

	dbHelper.(*mocks.DatabaseHelper).
		On("Collection", "planets").
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	err := planetDao.Delete(context.Background(), bson.M{"db-error": true})
	assert.EqualError(t, err, "mocked-db-error")

	err = planetDao.Delete(context.Background(), bson.M{"notFound-error": true})
	assert.EqualError(t, err, "document not found")

	err = planetDao.Delete(context.Background(), bson.M{"error": false})
	assert.NoError(t, err, "document not found")
}

func Test_planetsDAO_FindByName(t *testing.T) {
	// Define variables for interfaces
	var dbHelper db.DatabaseHelper

	//Set interfaces implementation to mocked structures
	dbHelper = &mocks.DatabaseHelper{}

	planetDao := NewPlanetsDao(dbHelper)
	planetDao = &mocks.PlanetsDAO{}

	planetDao.(*mocks.PlanetsDAO).
		On("FindByName", context.Background(), "mocked-planet-error").
		Return(nil, errors.New("mocked-error"))

	planetDao.(*mocks.PlanetsDAO).
		On("FindByName", context.Background(), "mocked-planet-correct").
		Return([]models.Planet{
			{Name: "mocked-planet"},
		}, nil)

	planets, err := planetDao.FindByName(context.Background(), "mocked-planet-error")
	assert.Empty(t, planets)
	assert.EqualError(t, err, "mocked-error")

	planets, err = planetDao.FindByName(context.Background(), "mocked-planet-correct")
	expected := []models.Planet{
		{Name: "mocked-planet"},
	}
	assert.Equal(t, expected, planets)
	assert.NoError(t, err)
}
