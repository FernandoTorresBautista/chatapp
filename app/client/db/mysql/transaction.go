package mysql

import (
	"database/sql"
	"fmt"
)

// Common failed messages
const (
	connectionUnavailableMsg = "db connection unavailable"
	queryExecFailed          = "repo query execution failed: %v"
)

// GetExample ...
func (r *Repo) GetExample() (interface{}, error) {
	if !r.isReady() {
		r.logger.Printf(connectionUnavailableMsg)
		return nil, fmt.Errorf(connectionUnavailableMsg)
	}
	query := `SELECT * FROM ` + r.DBName + `.table_example`
	var allTasks interface{}
	err := r.SQLDB.Select(&allTasks, query)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}
		r.logger.Printf("repo query GetExample execution failed: %v", err)
		return nil, err
	}
	return allTasks, nil
}

// CreateExample ...
func (r *Repo) CreateExample(fiel string) (interface{}, error) {
	if !r.isReady() {
		r.logger.Printf(connectionUnavailableMsg)
		return nil, fmt.Errorf(connectionUnavailableMsg)
	}
	query := `INSERT INTO ` + r.DBName + `.table_example (field_example) VALUES (?)`
	res, err := r.SQLDB.Exec(query, fiel)
	if err != nil {
		return nil, err
	}
	return res, nil
}
