package dao

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wallacebenevides/star-wars-api/internal/app/starwars/models"
	"github.com/wallacebenevides/star-wars-api/internal/pkg/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_planetsDAO_FindAll(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	cursor := &mocks.CursorHelper{}

	id := "5e27096d0c326694932a4cc8"

	expected := []models.Planet{
		{ID: id, Name: "mocked-planet", Climate: ""},
	}

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	collectionHelper.
		On("Find", context.Background(), primitive.D{{}}).
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
	planets, err := dao.FindAll(context.Background())

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
		On("Find", context.Background(), mock.Anything).
		Once().
		Return(nil, errors.New("mocked-error"))

	dao := NewPlanetsDao(dbHelper)

	planets, err := dao.FindAll(context.Background())
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
		On("Find", context.Background(), primitive.D{{}}).
		Once().
		Return(cursor, nil)

	cursor.On("All", mock.Anything, mock.Anything).
		Once().
		Return(errors.New("mocked-error"))
	cursor.On("Close", context.Background()).Return(nil)

	dao := NewPlanetsDao(dbHelper)

	planets, err := dao.FindAll(context.Background())
	assert.Empty(t, planets)
	assert.EqualError(t, err, "mocked-error")
}

func Test_planetsDAO_FindByID(t *testing.T) {

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
		On("FindOne", context.Background(), mock.Anything).
		Once().
		Return(srHelperCorrect)

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	id := "5e27096d0c326694932a4cc8"
	planet, err := planetDao.FindByID(context.Background(), id)
	assert.Equal(t, models.Planet{Name: "mocked-planet"}, planet)
	assert.NoError(t, err)
}

func Test_planetsDAO_FindById_with_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	srHelperErr := &mocks.SingleResultHelper{}

	srHelperErr.
		On("Decode", mock.AnythingOfType("*models.Planet")).
		Once().
		Return(errors.New("mocked-error"))

	collectionHelper.
		On("FindOne", context.Background(), mock.Anything).
		Once().
		Return(srHelperErr)

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	id := "5e27096d0c326694932a4cc8"
	planet, err := planetDao.FindByID(context.Background(), id)
	assert.Empty(t, planet)
	assert.EqualError(t, err, "mocked-error")
}

func Test_planetsDAO_FindById_with_invalid_id_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}

	planetDao := NewPlanetsDao(dbHelper)

	id := "invalid id"
	planet, err := planetDao.FindByID(context.Background(), id)
	assert.Empty(t, planet)
	assert.EqualError(t, err, INVALID_ID_ERROR_MESSAGE)
}

func Test_planetsDAO_Create(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	idExpected := "5e27096d0c326694932a4cc8"
	objectID, err := primitive.ObjectIDFromHex(idExpected)
	if err != nil {
		log.Panic(err.Error())
	}
	insertOneResult := &mongo.InsertOneResult{InsertedID: objectID}

	collectionHelper.
		On("InsertOne", context.Background(), mock.Anything).
		Once().
		Return(insertOneResult, nil)

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	insertedID, err := planetDao.Create(context.Background(), &models.Planet{Name: "mocked-planet-correct"})
	assert.Equal(t, idExpected, insertedID)
	assert.NoError(t, err)
}

func Test_planetsDAO_Create_with_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	collectionHelper.
		On("InsertOne", context.Background(), &models.PlanetDocument{Name: "mocked-planet-error"}).
		Once().
		Return(nil, errors.New("mocked-error"))

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	_, err := planetDao.Create(context.Background(), &models.Planet{Name: "mocked-planet-error"})
	assert.EqualError(t, err, "mocked-error")
}

func Test_planetsDAO_Delete(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	deleteResultCorrect := mongo.DeleteResult{DeletedCount: 1}

	collectionHelper.
		On("DeleteOne", context.Background(), mock.Anything).
		Once().
		Return(&deleteResultCorrect, nil)

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	// VALID ID
	id := "5e27096d0c326694932a4cc8"
	err := planetDao.Delete(context.Background(), id)
	assert.NoError(t, err, "document not found")
}

func Test_planetsDAO_Delete_with_notFound_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	collectionHelper.
		On("DeleteOne", context.Background(), mock.Anything).
		Once().
		Return(nil, errors.New("document not found"))

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	// VALID ID
	id := "5e27096d0c326694932a4cc8"
	err := planetDao.Delete(context.Background(), id)
	assert.EqualError(t, err, "document not found")
}

func Test_planetsDAO_Delete_with_invalid_id_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}

	planetDao := NewPlanetsDao(dbHelper)

	err := planetDao.Delete(context.Background(), "INVALID ID")
	assert.EqualError(t, err, INVALID_ID_ERROR_MESSAGE)
}

func Test_planetsDAO_Delete_with_db_error(t *testing.T) {

	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}

	collectionHelper.
		On("DeleteOne", context.Background(), mock.Anything).
		Once().
		Return(nil, errors.New("mocked-db-error"))

	dbHelper.
		On("Collection", "planets").
		Once().
		Return(collectionHelper)

	planetDao := NewPlanetsDao(dbHelper)

	// VALID ID
	id := "5e27096d0c326694932a4cc8"
	err := planetDao.Delete(context.Background(), id)
	assert.EqualError(t, err, "mocked-db-error")
}

func Test_planetsDAO_FindByName(t *testing.T) {
	dbHelper := &mocks.DatabaseHelper{}
	collectionHelper := &mocks.CollectionHelper{}
	cursor := &mocks.CursorHelper{}

	id := "5e27096d0c326694932a4cc8"

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
