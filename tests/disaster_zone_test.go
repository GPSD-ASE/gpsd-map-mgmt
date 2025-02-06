package tests

import (
	"database/sql"
	"database/sql/driver"
	"disaster-response-map-api/internal/handlers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ✅ MockDB struct to hold the mock database
type MockDB struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
}

// ✅ Create a new mock database instance
func NewMockDB() (*MockDB, error) {
	gin.SetMode(gin.ReleaseMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	return &MockDB{DB: db, Mock: mock}, nil
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	// ✅ Convert args to []driver.Value
	driverArgs := make([]driver.Value, len(args))
	for i, arg := range args {
		driverArgs[i] = arg
	}

	// ✅ Expect the query execution
	m.Mock.ExpectExec(query).WithArgs(driverArgs...).WillReturnResult(sqlmock.NewResult(1, 1))

	// ✅ Execute the query and return results
	return m.DB.Exec(query, args...)
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// ✅ Mock a database response with proper execution
	m.Mock.ExpectQuery(query).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "latitude", "longitude", "radius"}).
			AddRow(1, "Flood Zone", 53.349805, -6.26031, 500),
	)
	return m.DB.Query(query, args...)
}

func (m *MockDB) Close() error {
	return m.DB.Close()
}

// ✅ Test: Create Disaster Zone
func TestCreateDisasterZone(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := handlers.NewDisasterZoneHandler(mockDB)

	router := gin.Default()
	router.POST("/zones", handler.CreateDisasterZone)

	// ✅ Mock expected database behavior
	mockDB.Mock.ExpectExec("INSERT INTO disaster_zones").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Send a valid JSON payload
	body := `{"name": "Test Zone", "latitude": 53.349805, "longitude": -6.26031, "radius": 500}`
	req, _ := http.NewRequest("POST", "/zones", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// ✅ Test: Get Disaster Zones
func TestGetDisasterZones(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := handlers.NewDisasterZoneHandler(mockDB)

	router := gin.Default()
	router.GET("/zones", handler.GetDisasterZones)

	// ✅ Mock a database response
	mockDB.Mock.ExpectQuery("SELECT id, name, latitude, longitude, radius FROM disaster_zones").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "latitude", "longitude", "radius"}).
			AddRow(1, "Flood Zone", 53.349805, -6.26031, 500))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/zones", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ✅ Test: Delete Disaster Zone
func TestDeleteDisasterZone(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := handlers.NewDisasterZoneHandler(mockDB)

	router := gin.Default()
	router.DELETE("/zones/:id", handler.DeleteDisasterZone)

	// ✅ Mock the delete behavior
	mockDB.Mock.ExpectExec("DELETE FROM disaster_zones WHERE id = ?").
		WillReturnResult(sqlmock.NewResult(0, 1)) // ✅ Simulates successful delete

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/zones/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
