package repository_test

import "database/sql"

type MockDatabase struct {
	QueryRowResult *sql.Row
	QueryResult    *sql.Rows
	ExecResult     sql.Result
	// Add other fields for expected results
}

func (m *MockDatabase) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.QueryRowResult
}

func (m *MockDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.QueryResult, nil
}

func (m *MockDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.ExecResult, nil
}
