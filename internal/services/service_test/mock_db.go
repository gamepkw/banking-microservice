package service

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type MockDatabase struct {
	QueryRowResult *sql.Row
	QueryResult    *sql.Rows
	ExecResult     sql.Result
	// Add other fields for expected results
}

func NewMockDB(t *testing.T) sqlmock.Sqlmock {
	_, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}

	// Return the mock database and SQL mock instance
	return mock
}

// func NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("Error creating mock DB: %v", err)
// 	}

// 	// Return the mock database and SQL mock instance
// 	return db, mock
// }

func (m *MockDatabase) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.QueryRowResult
}

func (m *MockDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.QueryResult, nil
}

func (m *MockDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.ExecResult, nil
}
