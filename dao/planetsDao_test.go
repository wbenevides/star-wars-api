package dao

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wallacebenevides/star-wars-api/mocks"
	"github.com/wallacebenevides/star-wars-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_planetsDAO_FindAll(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	cursor := &mocks.CursorHelper{}

	id, _ := primitive.ObjectIDFromHex("12345")

	expected := []models.Planet{
		{ID: id, Name: "mocked-planet", Climate: ""},
	}

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	collectionHelper.
		On("Find", context.Background(), mock.Anything).
		Once().
		Return(cursor, nil)

	cursor.On("Close", context.Background()).Return(nil)

	cursor.On("All", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*[]models.Planet)
			*arg = expected
		}).
		Return(nil)

	dao := NewPlanetsDao(dbHelper)
	planets, err := dao.FindAll(context.Background(), bson.M{"error": false})

	assert.Equal(t, expected, planets)
	assert.NoError(t, err)
}

func Test_planetsDAO_FindAll_with_error_on_find(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	collectionHelper.
		On("Find", context.Background(), bson.M{"mocked-error": true}).
		Once().
		Return(nil, errors.New("mocked-error"))

	dao := NewPlanetsDao(dbHelper)

	planets, err := dao.FindAll(context.Background(), bson.M{"mocked-error": true})
	assert.Empty(t, planets)
	assert.EqualError(t, err, "mocked-error")
}

func Test_planetsDAO_FindAll_with_error_on_all(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	cursor := &mocks.CursorHelper{}

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	collectionHelper.
		On("Find", context.Background(), bson.M{"mocked-error": true}).
		Once().
		Return(cursor, nil)

	cursor.On("All", mock.Anything, mock.Anything).
		Once().
		Return(errors.New("mocked-error"))
	cursor.On("Close", context.Background()).Return(nil)

	dao := NewPlanetsDao(dbHelper)

	planets, err := dao.FindAll(context.Background(), bson.M{"mocked-error": true})
	assert.Empty(t, planets)
	assert.EqualError(t, err, "mocked-error")
}

func Test_planetsDAO_FindOne(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	srHelperCorrect := &mocks.SingleResultHelper{}

	srHelperCorrect.
		On("Decode", mock.AnythingOfType("*models.Planet")).
		Once().
		Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Planet)
		arg.Name = "mocked-planet"
	})

	collectionHelper.
		On("FindOne", context.Background(), bson.M{"error": false}).
		Once().
		Return(srHelperCorrect)

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	planet, err := planetDao.FindOne(context.Background(), bson.M{"error": false})
	assert.Equal(t, &models.Planet{Name: "mocked-planet"}, planet)
	assert.NoError(t, err)
}

func Test_planetsDAO_FindOne_with_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	srHelperErr := &mocks.SingleResultHelper{}

	srHelperErr.
		On("Decode", mock.AnythingOfType("*models.Planet")).
		Once().
		Return(errors.New("mocked-error"))

	collectionHelper.
		On("FindOne", context.Background(), bson.M{"error": true}).
		Once().
		Return(srHelperErr)

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	planet, err := planetDao.FindOne(context.Background(), bson.M{"error": true})
	assert.Empty(t, planet)
	assert.EqualError(t, err, "mocked-error")
}

func Test_planetsDAO_Create(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	collectionHelper.
		On("InsertOne", context.Background(), &models.Planet{Name: "mocked-planet-correct"}).
		Once().
		Return(nil, nil)

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	err := planetDao.Create(context.Background(), &models.Planet{Name: "mocked-planet-correct"})
	assert.NoError(t, err)
}

func Test_planetsDAO_Create_with_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	collectionHelper.
		On("InsertOne", context.Background(), &models.Planet{Name: "mocked-planet-error"}).
		Once().
		Return(nil, errors.New("mocked-error"))

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	err := planetDao.Create(context.Background(), &models.Planet{Name: "mocked-planet-error"})
	assert.EqualError(t, err, "mocked-error")
}

func Test_planetsDAO_Delete(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	deleteResultCorrect := mongo.DeleteResult{DeletedCount: 1}

	collectionHelper.
		On("DeleteOne", context.Background(), bson.M{"error": false}).
		Once().
		Return(&deleteResultCorrect, nil)

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	err := planetDao.Delete(context.Background(), bson.M{"error": false})
	assert.NoError(t, err, "document not found")
}

func Test_planetsDAO_Delete_with_notFound_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	collectionHelper.
		On("DeleteOne", context.Background(), bson.M{"notFound-error": true}).
		Once().
		Return(nil, errors.New("document not found"))

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	err := planetDao.Delete(context.Background(), bson.M{"notFound-error": true})
	assert.EqualError(t, err, "document not found")
}

func Test_planetsDAO_Delete_with_db_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	collectionHelper.
		On("DeleteOne", context.Background(), bson.M{"db-error": true}).
		Once().
		Return(nil, errors.New("mocked-db-error"))

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	err := planetDao.Delete(context.Background(), bson.M{"db-error": true})
	assert.EqualError(t, err, "mocked-db-error")
}

func Test_planetsDAO_FindByName(t *testing.T) {
	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	cursor := &mocks.CursorHelper{}

	id, _ := primitive.ObjectIDFromHex("12345")

	expected := []models.Planet{
		{ID: id, Name: "mocked-planet", Climate: ""},
	}

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	collectionHelper.
		On("Find", context.Background(), mock.Anything).
		Once().
		Return(cursor, nil)

	cursor.On("Close", context.Background()).Return(nil)

	cursor.On("All", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*[]models.Planet)
			*arg = expected
		}).
		Return(nil)

	dao := NewPlanetsDao(dbHelper)
	planets, err := dao.FindByName(context.Background(), "mocked-planet")

	assert.Equal(t, expected, planets)
	assert.NoError(t, err)
}
