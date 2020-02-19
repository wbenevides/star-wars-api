package resources

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	planetDao := &mocks.PlanetsDAO{}
	dataMock := []models.Planet{{Name: "mocked-planet"}}
	planetDao.
		On("FindAll", context.TODO()).
		Once().
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

func TestPlanetHandler_GetAll_with_error(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/planets", nil)
	if err != nil {
		t.Fatal(err)
	}
	planetDao := &mocks.PlanetsDAO{}
	planetDao.
		On("FindAll", context.TODO()).
		Once().
		Return(nil, errors.New("mocked-error"))

	rr := httptest.NewRecorder()
	getAll := NewPlanetHandler(planetDao).GetAll()
	handler := http.HandlerFunc(getAll)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check the response body is what we expect.
	got := rr.Body.String()
	expected := `{"error":"` + INTERNAL_SERVER_ERROR_MESSAGE + `"}`

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
		Once().
		Return(nil, nil)

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

func TestPlanetHandler_Create_with_error(t *testing.T) {
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
		Once().
		Return(nil, errors.New("mocked-error"))

	rr := httptest.NewRecorder()

	create := NewPlanetHandler(planetDao).Create()
	handler := http.HandlerFunc(create)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"` + INTERNAL_SERVER_ERROR_MESSAGE + `"}`
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Create_with_bad_request_error(t *testing.T) {

	payload := `{"name":0}`
	jsonStr := []byte(payload)

	req, err := http.NewRequest(http.MethodPost, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}

	rr := httptest.NewRecorder()

	create := NewPlanetHandler(planetDao).Create()
	handler := http.HandlerFunc(create)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error":"Invalid request payload"}`
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
	planetDao.
		On("FindByID", context.TODO(), id).
		Once().
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

func TestPlanetHandler_GetByID_with_error(t *testing.T) {
	id := "5e27096d0c326694932a4cc8"
	path := fmt.Sprintf("/api/planets/%s", id)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	planetDao := &mocks.PlanetsDAO{}
	planetDao.
		On("FindByID", context.TODO(), id).
		Once().
		Return(nil, errors.New("mocked-error"))

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	getByID := NewPlanetHandler(planetDao).GetByID()
	router.HandleFunc("/api/planets/{id}", getByID)
	router.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"` + INTERNAL_SERVER_ERROR_MESSAGE + `"}`

	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_GetByID_with_bad_request_error(t *testing.T) {
	id := "invalidId"
	path := fmt.Sprintf("/api/planets/%s", id)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	planetDao := &mocks.PlanetsDAO{}
	planetDao.
		On("FindByID", context.TODO(), id).
		Once().
		Return(nil, errors.New(dao.INVALID_ID_ERROR_MESSAGE))

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	getByID := NewPlanetHandler(planetDao).GetByID()
	router.HandleFunc("/api/planets/{id}", getByID)
	router.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error":"Invalid Planet ID"}`
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_GetByID_with_not_found(t *testing.T) {
	id := "5e27096d0c326694932a4cc8"
	path := fmt.Sprintf("/api/planets/%s", id)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	planetDao := &mocks.PlanetsDAO{}
	planetDao.
		On("FindByID", context.TODO(), id).
		Return(nil, errors.New(dao.NOT_FOUND_ERROR_MESSAGE))

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	getByID := NewPlanetHandler(planetDao).GetByID()
	router.HandleFunc("/api/planets/{id}", getByID)
	router.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := `{"error":"document not found"}`
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
		Once().
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

func TestPlanetHandler_FindByName_with_error(t *testing.T) {
	name := "mocked-planet"
	req, err := http.NewRequest(http.MethodGet, "/api/planets/findByName?name="+name, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}

	planetDao.
		On("FindByName", context.TODO(), name).
		Once().
		Return(nil, errors.New("mocked-error"))

	rr := httptest.NewRecorder()

	findByName := NewPlanetHandler(planetDao).FindByName()
	handler := http.HandlerFunc(findByName)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"error":"` + INTERNAL_SERVER_ERROR_MESSAGE + `"}`

	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_FindByName_with_not_found(t *testing.T) {
	name := "mocked-planet"
	req, err := http.NewRequest(http.MethodGet, "/api/planets/findByName?name="+name, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}
	dataMock := []models.Planet{}

	planetDao.
		On("FindByName", context.TODO(), name).
		Once().
		Return(dataMock, nil)

	rr := httptest.NewRecorder()

	findByName := NewPlanetHandler(planetDao).FindByName()
	handler := http.HandlerFunc(findByName)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := `{"error":"document not found"}`

	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Delete(t *testing.T) {
	id := "5e270a857247f2102f213565"
	payload := fmt.Sprintf(`{"id": "%s"}`, id)
	jsonStr := []byte(payload)

	req, err := http.NewRequest(http.MethodDelete, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}

	planetDao.
		On("Delete", context.TODO(), id).
		Once().
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

func TestPlanetHandler_Delete_with_bad_request_error(t *testing.T) {
	payload := `{"id":0}`
	jsonStr := []byte(payload)

	req, err := http.NewRequest(http.MethodDelete, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}

	rr := httptest.NewRecorder()

	delete := NewPlanetHandler(planetDao).Delete()
	handler := http.HandlerFunc(delete)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error":"Invalid request payload"}`
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Delete_with_bad_request_hexadecimal_id_error(t *testing.T) {
	// payload not in hexadecimal
	payload := `{"id":"5e270a857247f2102f21356z"}`
	jsonStr := []byte(payload)

	req, err := http.NewRequest(http.MethodDelete, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}
	planetDao.On("Delete", mock.Anything, mock.Anything).
		Once().
		Return(errors.New(dao.INVALID_ID_ERROR_MESSAGE))

	rr := httptest.NewRecorder()

	delete := NewPlanetHandler(planetDao).Delete()
	handler := http.HandlerFunc(delete)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"error":"%s"}`, dao.INVALID_ID_ERROR_MESSAGE)
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Delete_with_not_found(t *testing.T) {
	payload := `{"id":"5e270a857247f2102f213565"}`
	jsonStr := []byte(payload)

	req, err := http.NewRequest(http.MethodDelete, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	planetDao := &mocks.PlanetsDAO{}

	planetDao.
		On("Delete", mock.Anything, mock.Anything).
		Once().
		Return(errors.New(dao.NOT_FOUND_ERROR_MESSAGE))

	rr := httptest.NewRecorder()

	delete := NewPlanetHandler(planetDao).Delete()
	handler := http.HandlerFunc(delete)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := fmt.Sprintf(`{"error":"%s"}`, dao.NOT_FOUND_ERROR_MESSAGE)
	got := rr.Body.String()

	assert.Equal(t, expected, got)
}
