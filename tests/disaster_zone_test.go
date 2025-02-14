package tests

import (
	"database/sql"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
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
	return m.DB.Query(query, args...)
}

func (m *MockDB) Close() error {
	return m.DB.Close()
}

// ✅ Test: Get Disaster Zones
// func TestGetDisasterZones(t *testing.T) {
// 	mockDB, _ := NewMockDB()
// 	handler := handlers.NewDisasterZoneHandler(mockDB)

// 	router := gin.Default()
// 	router.GET("/zones", handler.GetDisasterZones)

// 	// ✅ Mock a database response
// 	mockDB.Mock.ExpectQuery(regexp.QuoteMeta("SELECT i.incident_id as incident_id, t.type_name AS incident_name, i.latitude as latitude, i.longitude as longitude, i.severity_id as severity_id FROM incident i JOIN incident_type t ON i.type_id = t.type_id; ")).
// 		WillReturnRows(sqlmock.NewRows([]string{"incident_id", "incident_name", "latitude", "longitude", "severity_id"}).
// 			AddRow(1, "Flood Zone", 53.349805, -6.26031, 500))

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/zones", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// }
