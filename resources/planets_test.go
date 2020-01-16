package resources

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/wallacebenevides/star-wars-api/dao"
	"github.com/wallacebenevides/star-wars-api/db"
)

func TestNewPlanetHandler(t *testing.T) {
	type args struct {
		db db.DatabaseHelper
	}
	tests := []struct {
		name string
		args args
		want *PlanetHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlanetHandler(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlanetHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanetHandler_GetAll(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/planets", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getAll := NewPlanetHandler(nil).GetAll()
	handler := http.HandlerFunc(getAll)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Error("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	got := rr.Body.String()
	expected := `[{"id":"5e1ffbfb9a01931f37757c58","name":"Alderaan","climate":"temperate","terrain":"grasslands, mountains","films":2},{"id":"5e1ffbfb9a01931f37757c59","name":"Yavin IV","climate":"temperate, tropical","terrain":"jungle, rainforests","films":3},{"id":"5e1ffbfb9a01931f37757c5a","name":"Hoth","climate":"frozen","terrain":"tundra, ice caves, mountain ranges","films":1},{"id":"5e1ffbfb9a01931f37757c5b","name":"Dagobah","climate":"murky","terrain":"swamp, jungles","films":3},{"id":"5e1ffbfb9a01931f37757c5c","name":"Bespin","climate":"temperate","terrain":"gas giant","films":1},{"id":"5e1ffbfb9a01931f37757c5d","name":"Endor","climate":"temperate","terrain":"forests, mountains, lakes","films":1},{"id":"5e1ffbfb9a01931f37757c5e","name":"Naboo","climate":"temperate","terrain":"grassy hills, swamps, forests, mountains","films":4},{"id":"5e1ffbfb9a01931f37757c5f","name":"Coruscant","climate":"temperate","terrain":"cityscape, mountains","films":4},{"id":"5e1ffbfb9a01931f37757c60","name":"Kamino","climate":"temperate","terrain":"ocean","films":1},{"id":"5e1ffbfb9a01931f37757c61","name":"Geonosis","climate":"temperate, arid","terrain":"rock, desert, mountain, barren","films":1},{"id":"5e1ffbfb9a01931f37757c62","name":"Utapau","climate":"temperate, arid, windy","terrain":"scrublands, savanna, canyons, sinkholes","films":1},{"id":"5e1ffbfb9a01931f37757c63","name":"Mustafar","climate":"hot","terrain":"volcanoes, lava rivers, mountains, caves","films":1},{"id":"5e1ffbfb9a01931f37757c64","name":"Kashyyyk","climate":"tropical","terrain":"jungle, forests, lakes, rivers","films":1},{"id":"5e1ffbfb9a01931f37757c65","name":"Polis Massa","climate":"artificial temperate ","terrain":"airless asteroid","films":1},{"id":"5e1ffbfb9a01931f37757c66","name":"Mygeeto","climate":"frigid","terrain":"glaciers, mountains, ice canyons","films":1},{"id":"5e1ffbfb9a01931f37757c67","name":"Felucia","climate":"hot, humid","terrain":"fungus forests","films":1},{"id":"5e1ffbfb9a01931f37757c68","name":"Cato Neimoidia","climate":"temperate, moist","terrain":"mountains, fields, forests, rock arches","films":1},{"id":"5e1ffbfb9a01931f37757c69","name":"Saleucami","climate":"hot","terrain":"caves, desert, mountains, volcanoes","films":1},{"id":"5e1ffbfb9a01931f37757c6a","name":"Stewjon","climate":"temperate","terrain":"grass","films":0},{"id":"5e1ffbfb9a01931f37757c6b","name":"Eriadu","climate":"polluted","terrain":"cityscape","films":0},{"id":"5e1ffc13ad224e96489f0b73","name":"Lua","climate":"temperate","terrain":"grasslands, mountains","films":2},{"id":"5e1ffc27ad224e96489f0b74","name":"Lua","climate":"temperate","terrain":"grasslands, mountains","films":2}]`
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}

func TestPlanetHandler_Create(t *testing.T) {
	payload := `{"id":"5e1ffc27ad224e96489f0b75","name":"Lua","climate":"temperate","terrain":"grasslands, mountains","films":2}`
	jsonStr := []byte(payload)

	req, err := http.NewRequest("POST", "/entry", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")
	rr := httptest.NewRecorder()
	create := NewPlanetHandler(nil).Create()
	handler := http.HandlerFunc(create)
	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := payload
	got := rr.Body.String()
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}

func TestPlanetHandler_GetByID(t *testing.T) {
	type fields struct {
		db dao.PlanetsDAO
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PlanetHandler{
				db: tt.fields.db,
			}
			if got := h.GetByID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanetHandler.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanetHandler_FindByName(t *testing.T) {
	type fields struct {
		db dao.PlanetsDAO
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PlanetHandler{
				db: tt.fields.db,
			}
			if got := h.FindByName(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanetHandler.FindByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanetHandler_Delete(t *testing.T) {
	type fields struct {
		db dao.PlanetsDAO
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &PlanetHandler{
				db: tt.fields.db,
			}
			if got := h.Delete(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlanetHandler.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_respondWithError(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		code int
		msg  string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respondWithError(tt.args.w, tt.args.code, tt.args.msg)
		})
	}
}

func Test_respondWithJson(t *testing.T) {
	type args struct {
		w       http.ResponseWriter
		code    int
		payload interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respondWithJson(tt.args.w, tt.args.code, tt.args.payload)
		})
	}
}

func Test_createSuccessResult(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSuccessResult(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createSuccessResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
