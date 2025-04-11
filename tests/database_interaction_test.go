package tests

import (
	"database/sql"
	"fmt"
	"testing"

	"disaster-response-map-api/pkg/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_Exec_Happy(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	d := &database.Database{DB: db}

	mock.ExpectExec("INSERT INTO test_table").
		WithArgs("value1", 123).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := d.Exec("INSERT INTO test_table (col1, col2) VALUES ($1, $2)", "value1", 123)
	assert.NoError(t, err)

	rowsAffected, err := result.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)
}

func TestDatabase_Query_Happy(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	d := &database.Database{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob")

	mock.ExpectQuery("SELECT id, name FROM test_table").
		WillReturnRows(rows)

	resultRows, err := d.Query("SELECT id, name FROM test_table")
	assert.NoError(t, err)
	defer resultRows.Close()

	var count int
	for resultRows.Next() {
		var id int
		var name string
		err = resultRows.Scan(&id, &name)
		assert.NoError(t, err)
		count++
	}
	assert.Equal(t, 2, count)
}

func TestNewDatabase_Happy(t *testing.T) {
	dsn := "postgres://test_user:test_pass@localhost:5432/test_db"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Skip("Skipping actual DB connectivity test; using sqlmock is recommended for unit testing")
	}

	fmt.Println("Connected to database successfully")
}
