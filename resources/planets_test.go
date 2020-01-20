package resources

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallacebenevides/star-wars-api/dao"
	"github.com/wallacebenevides/star-wars-api/mocks"
	"github.com/wallacebenevides/star-wars-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPlanetHandler_GetAll(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/planets", nil)
	if err != nil {
		t.Fatal(err)
	}

	dbHelper := &mocks.DatabaseHelper{}
	planetDao := dao.NewPlanetsDao(dbHelper)
	planetDao = &mocks.PlanetsDAO{}

	dataMock := []models.Planet{{Name: "mocked-planet"}}

	planetDao.(*mocks.PlanetsDAO).
		On("FindAll").
		Return(dataMock, nil)

	planetDao.(*mocks.PlanetsDAO).
		On("FindAll").
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
	expected := `[{"name":"mocked-planet"}]`
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
	dataMock := []models.Planet{{Name: "mocked-planet"}}

	planetDao.
		On("Create").
		Return(dataMock, nil)

	rr := httptest.NewRecorder()

	create := NewPlanetHandler(planetDao).Create()
	handler := http.HandlerFunc(create)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := payload
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_GetByID(t *testing.T) {
	id := "12345"
	objectID, _ := primitive.ObjectIDFromHex(id)
	req, err := http.NewRequest(http.MethodGet, "/api/planets/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}
	dataMock := models.Planet{ID: objectID}

	planetDao.
		On("FindOne").
		Return(dataMock, nil)

	rr := httptest.NewRecorder()

	getByID := NewPlanetHandler(planetDao).GetByID()
	handler := http.HandlerFunc(getByID)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected, _ := json.Marshal(dataMock)
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
	dataMock := models.Planet{Name: "mocked-planet"}

	planetDao.
		On("FindByName").
		Return(dataMock, nil)

	rr := httptest.NewRecorder()

	findByName := NewPlanetHandler(planetDao).FindByName()
	handler := http.HandlerFunc(findByName)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected, _ := json.Marshal(dataMock)
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Delete(t *testing.T) {
	payload := `{"name":"mocked-planet"}`
	jsonStr := []byte(payload)

	req, err := http.NewRequest(http.MethodDelete, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}

	planetDao.
		On("Delete").
		Return(nil)

	rr := httptest.NewRecorder()

	delete := NewPlanetHandler(planetDao).Delete()
	handler := http.HandlerFunc(delete)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := payload
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}
