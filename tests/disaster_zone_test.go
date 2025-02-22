package tests

import (
	"database/sql"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

type MockDB struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
}

func NewMockDB() (*MockDB, error) {
	gin.SetMode(gin.ReleaseMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	return &MockDB{DB: db, Mock: mock}, nil
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	driverArgs := make([]driver.Value, len(args))
	for i, arg := range args {
		driverArgs[i] = arg
	}

	m.Mock.ExpectExec(query).WithArgs(driverArgs...).WillReturnResult(sqlmock.NewResult(1, 1))

	return m.DB.Exec(query, args...)
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.DB.Query(query, args...)
}

func (m *MockDB) Close() error {
	return m.DB.Close()
}
