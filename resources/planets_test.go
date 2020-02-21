package resources

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wallacebenevides/star-wars-api/dao"
	"github.com/wallacebenevides/star-wars-api/mocks"
	"github.com/wallacebenevides/star-wars-api/models"
)

func TestPlanetHandler_GetAll(t *testing.T) {
	// mock declaration
	planetDao := &mocks.PlanetsDAO{}
	dataMock := []models.Planet{{Name: "mocked-planet"}}
	planetDao.
		On("FindAll", context.TODO()).
		Once().
		Return(dataMock, nil)

	// request logic
	req, err := http.NewRequest(http.MethodGet, "/api/planets", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getAll := NewPlanetHandler(planetDao).GetAll()
	handler := http.HandlerFunc(getAll)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	got := rr.Body.String()
	expected := `[{"id":"","name":"mocked-planet","climate":"","terrain":"","films":0}]`
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_GetAll_with_error(t *testing.T) {
	// mock declaration
	planetDao := &mocks.PlanetsDAO{}
	planetDao.
		On("FindAll", context.TODO()).
		Once().
		Return(nil, errors.New("mocked-error"))

	// request logic
	req, err := http.NewRequest(http.MethodGet, "/api/planets", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getAll := NewPlanetHandler(planetDao).GetAll()
	handler := http.HandlerFunc(getAll)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	got := rr.Body.String()
	expected := `{"error":"` + INTERNAL_SERVER_ERROR_MESSAGE + `"}`
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Create(t *testing.T) {
	// mock declaration
	planetDao := &mocks.PlanetsDAO{}
	mockedID := "5e27096d0c326694932a4cc8"
	planetDao.
		On("Create", context.TODO(), mock.Anything).
		Once().
		Return(mockedID, nil)

	payload := `{"name":"mocked-planet"}`
	jsonStr := []byte(payload)
	// request logic
	req, err := http.NewRequest(http.MethodPost, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	rr := httptest.NewRecorder()
	create := NewPlanetHandler(planetDao).Create()
	handler := http.HandlerFunc(create)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := fmt.Sprintf(`{"id":"%s","name":"mocked-planet","climate":"","terrain":"","films":0}`, mockedID)
	got := rr.Body.String()
	assert.Equal(t, expected, got)
	// Check the location header
	location := rr.Header().Get("Location")
	assert.True(t, strings.Contains(location, mockedID))
}

func TestPlanetHandler_Create_with_error(t *testing.T) {
	// mock declaration
	planetDao := &mocks.PlanetsDAO{}
	planetDao.
		On("Create", context.TODO(), mock.Anything).
		Once().
		Return("", errors.New("mocked-error"))

	payload := `{"name":"mocked-planet"}`
	jsonStr := []byte(payload)

	// request logic
	req, err := http.NewRequest(http.MethodPost, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	rr := httptest.NewRecorder()
	create := NewPlanetHandler(planetDao).Create()
	handler := http.HandlerFunc(create)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"error":"` + INTERNAL_SERVER_ERROR_MESSAGE + `"}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Create_with_bad_request_error(t *testing.T) {
	// mock declaration
	planetDao := &mocks.PlanetsDAO{}

	payload := `{"name":0}`
	jsonStr := []byte(payload)

	// request logic
	req, err := http.NewRequest(http.MethodPost, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	rr := httptest.NewRecorder()
	create := NewPlanetHandler(planetDao).Create()
	handler := http.HandlerFunc(create)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"error":"Invalid request payload"}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_GetByID(t *testing.T) {
	// mock declation
	id := "5e27096d0c326694932a4cc8"
	planetDao := &mocks.PlanetsDAO{}
	dataMock := models.Planet{ID: id}
	planetDao.
		On("FindByID", context.TODO(), id).
		Once().
		Return(dataMock, nil)

	// request logic
	req, err := http.NewRequest(http.MethodGet, "/api/planets/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	getByID := NewPlanetHandler(planetDao).GetByID()
	router.HandleFunc("/api/planets/{id}", getByID)
	router.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"id":"5e27096d0c326694932a4cc8","name":"","climate":"","terrain":"","films":0}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_GetByID_with_error(t *testing.T) {
	// mock declaration
	id := "5e27096d0c326694932a4cc8"
	planetDao := &mocks.PlanetsDAO{}
	var planet models.Planet
	planetDao.
		On("FindByID", context.TODO(), id).
		Once().
		Return(planet, errors.New("mocked-error"))

	// request logic
	path := "/api/planets/" + id
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	getByID := NewPlanetHandler(planetDao).GetByID()
	router.HandleFunc("/api/planets/{id}", getByID)
	router.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"error":"` + INTERNAL_SERVER_ERROR_MESSAGE + `"}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_GetByID_with_bad_request_error(t *testing.T) {
	// mock declaration
	id := "invalidId"
	planetDao := &mocks.PlanetsDAO{}
	var planet models.Planet
	planetDao.
		On("FindByID", context.TODO(), id).
		Once().
		Return(planet, errors.New(dao.INVALID_ID_ERROR_MESSAGE))

	// request logic
	path := "/api/planets/" + id
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	getByID := NewPlanetHandler(planetDao).GetByID()
	router.HandleFunc("/api/planets/{id}", getByID)
	router.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"error":"Invalid Planet ID"}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_GetByID_with_not_found(t *testing.T) {
	// mock declaration
	id := "5e27096d0c326694932a4cc8"
	planetDao := &mocks.PlanetsDAO{}
	var planet models.Planet
	planetDao.
		On("FindByID", context.TODO(), id).
		Return(planet, errors.New(dao.NOT_FOUND_ERROR_MESSAGE))

	// request logic
	path := "/api/planets/" + id
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	getByID := NewPlanetHandler(planetDao).GetByID()
	router.HandleFunc("/api/planets/{id}", getByID)
	router.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"error":"document not found"}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_FindByName(t *testing.T) {
	// mock declaration
	planetDao := &mocks.PlanetsDAO{}
	mockedName := "mocked-planet"
	dataMock := []models.Planet{{Name: mockedName}}
	planetDao.
		On("FindByName", context.TODO(), mockedName).
		Once().
		Return(dataMock, nil)

	// request logic
	req, err := http.NewRequest(http.MethodGet, "/api/planets/findByName?name="+mockedName, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	rr := httptest.NewRecorder()
	findByName := NewPlanetHandler(planetDao).FindByName()
	handler := http.HandlerFunc(findByName)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `[{"id":"","name":"mocked-planet","climate":"","terrain":"","films":0}]`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_FindByName_with_error(t *testing.T) {
	// mock declaration
	mockedName := "mocked-planet"
	planetDao := &mocks.PlanetsDAO{}
	planetDao.
		On("FindByName", context.TODO(), mockedName).
		Once().
		Return(nil, errors.New("mocked-error"))

	// request declation
	req, err := http.NewRequest(http.MethodGet, "/api/planets/findByName?name="+mockedName, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	rr := httptest.NewRecorder()
	findByName := NewPlanetHandler(planetDao).FindByName()
	handler := http.HandlerFunc(findByName)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"error":"` + INTERNAL_SERVER_ERROR_MESSAGE + `"}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_FindByName_with_not_found(t *testing.T) {
	// mock declation
	mockedName := "mocked-planet"
	planetDao := &mocks.PlanetsDAO{}
	dataMock := []models.Planet{}
	planetDao.
		On("FindByName", context.TODO(), mockedName).
		Once().
		Return(dataMock, nil)

	// request logic
	req, err := http.NewRequest(http.MethodGet, "/api/planets/findByName?name="+mockedName, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	rr := httptest.NewRecorder()
	findByName := NewPlanetHandler(planetDao).FindByName()
	handler := http.HandlerFunc(findByName)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"error":"document not found"}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Delete(t *testing.T) {
	// mock declaration
	mockedID := "5e270a857247f2102f213565"
	planetDao := &mocks.PlanetsDAO{}
	planetDao.
		On("Delete", context.TODO(), mockedID).
		Once().
		Return(nil)

	// request logic
	path := "/api/planets/" + mockedID
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	planetHandler := NewPlanetHandler(planetDao)
	router.HandleFunc("/api"+planetHandler.Routes().PLANETS_ID, planetHandler.Delete()).Methods(http.MethodDelete)
	router.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"result":"success"}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

/* func TestPlanetHandler_Delete_with_bad_request_error(t *testing.T) {
	// mock declaration
	planetDao := &mocks.PlanetsDAO{}
	payload := `{"id":0}`
	jsonStr := []byte(payload)

	// request logic
	req, err := http.NewRequest(http.MethodDelete, "/api/planets", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	rr := httptest.NewRecorder()
	delete := NewPlanetHandler(planetDao).Delete()
	handler := http.HandlerFunc(delete)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"error":"Invalid request payload"}`
	got := rr.Body.String()
	assert.Equal(t, expected, got)
} */

func TestPlanetHandler_Delete_with_bad_request_hexadecimal_id_error(t *testing.T) {
	// mock declaration
	planetDao := &mocks.PlanetsDAO{}
	planetDao.On("Delete", mock.Anything, mock.Anything).
		Once().
		Return(errors.New(dao.INVALID_ID_ERROR_MESSAGE))
	// payload not is hexadecimal
	mockedID := `{"id":"5e270a857247f2102f21356z"}`

	// request logic
	req, err := http.NewRequest(http.MethodDelete, "/api/planets/"+mockedID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	rr := httptest.NewRecorder()
	planetHandler := NewPlanetHandler(planetDao)
	router := mux.NewRouter()
	router.HandleFunc("/api"+planetHandler.Routes().PLANETS_ID, planetHandler.Delete()).Methods(http.MethodDelete)
	router.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := fmt.Sprintf(`{"error":"%s"}`, dao.INVALID_ID_ERROR_MESSAGE)
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}

func TestPlanetHandler_Delete_with_not_found(t *testing.T) {
	// mock declaration
	planetDao := &mocks.PlanetsDAO{}
	planetDao.
		On("Delete", mock.Anything, mock.Anything).
		Once().
		Return(errors.New(dao.NOT_FOUND_ERROR_MESSAGE))
	mockedID := `{"id":"5e270a857247f2102f213565"}`

	// request logic
	req, err := http.NewRequest(http.MethodDelete, "/api/planets/"+mockedID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	rr := httptest.NewRecorder()
	planetHandler := NewPlanetHandler(planetDao)
	router := mux.NewRouter()
	router.HandleFunc("/api"+planetHandler.Routes().PLANETS_ID, planetHandler.Delete()).Methods(http.MethodDelete)
	router.ServeHTTP(rr, req)

	// Check the status code.
	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := fmt.Sprintf(`{"error":"%s"}`, dao.NOT_FOUND_ERROR_MESSAGE)
	got := rr.Body.String()
	assert.Equal(t, expected, got)
}
