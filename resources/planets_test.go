package resources

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wallacebenevides/star-wars-api/mocks"
	"github.com/wallacebenevides/star-wars-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPlanetHandler_GetAll(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/planets", nil)
	if err != nil {
		t.Fatal(err)
	}
	planetDao := &mocks.PlanetsDAO{}
	dataMock := []models.Planet{{Name: "mocked-planet"}}
	planetDao.
		On("FindAll", context.TODO(), bson.D{{}}).
		Return(dataMock, nil)

	rr := httptest.NewRecorder()
	getAll := NewPlanetHandler(planetDao).GetAll()
	handler := http.HandlerFunc(getAll)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	got := rr.Body.String()
	expected := `[{"id":"000000000000000000000000","name":"mocked-planet","climate":"","terrain":"","films":0}]`

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Create(t *testing.T) {
	payload := `{"name":"mocked-planet"}`
	jsonStr := []byte(payload)

	req, err := http.NewRequest(http.MethodPost, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}
	objectID := &mocks.ObjectIDHelper{}
	id, _ := primitive.ObjectIDFromHex("12345")
	objectID.On("NewObjectID").Return(id)

	planetDao.
		On("Create", context.TODO(), mock.Anything).
		Return(nil)

	rr := httptest.NewRecorder()

	create := NewPlanetHandler(planetDao).Create()
	handler := http.HandlerFunc(create)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `{"result":"success"}`
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_GetByID(t *testing.T) {
	id := "5e27096d0c326694932a4cc8"
	path := fmt.Sprintf("/api/planets/%s", id)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	planetDao := &mocks.PlanetsDAO{}
	dataMock := models.Planet{ID: objectID}
	filter := bson.M{"_id": objectID}
	planetDao.
		On("FindOne", context.TODO(), filter).
		Return(&dataMock, nil)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	getByID := NewPlanetHandler(planetDao).GetByID()
	router.HandleFunc("/api/planets/{id}", getByID)
	router.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":"5e27096d0c326694932a4cc8","name":"","climate":"","terrain":"","films":0}`
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_FindByName(t *testing.T) {
	name := "mocked-planet"
	req, err := http.NewRequest(http.MethodGet, "/api/planets/findByName?name="+name, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}
	dataMock := []models.Planet{{Name: "mocked-planet"}}

	planetDao.
		On("FindByName", context.TODO(), name).
		Return(dataMock, nil)

	rr := httptest.NewRecorder()

	findByName := NewPlanetHandler(planetDao).FindByName()
	handler := http.HandlerFunc(findByName)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":"000000000000000000000000","name":"mocked-planet","climate":"","terrain":"","films":0}]`

	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Delete(t *testing.T) {
	payload := `{"id":"5e270a857247f2102f213565"}`
	jsonStr := []byte(payload)

	req, err := http.NewRequest(http.MethodDelete, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}

	idPrimitive, _ := primitive.ObjectIDFromHex("5e270a857247f2102f213565")
	filter := bson.M{"_id": idPrimitive}

	planetDao.
		On("Delete", context.TODO(), filter).
		Return(nil)

	rr := httptest.NewRecorder()

	delete := NewPlanetHandler(planetDao).Delete()
	handler := http.HandlerFunc(delete)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"result":"success"}`
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}
