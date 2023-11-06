package mysql

import (
	"database/sql"
	"fmt"

	"chatapp/app/biz/models"
)

// CreateUser create a new user
func (r *Repo) CreateUser(nu *models.User) (uint64, error) {
	if !r.isReady() {
		r.logger.Printf(connectionUnavailableMsg)
		return 0, fmt.Errorf(connectionUnavailableMsg)
	}
	query := `INSERT INTO ` + r.DBName + `.users (name, email, password) VALUES (?,?,?)`
	res, err := r.SQLDB.Exec(query, nu.Name, nu.Email, nu.Password)
	if err != nil {
		return 0, err
	}
	idx, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(idx), nil
}

// GetUser create a new user
func (r *Repo) GetUser(email, pass string) (*models.User, error) {
	if !r.isReady() {
		r.logger.Printf(connectionUnavailableMsg)
		return nil, fmt.Errorf(connectionUnavailableMsg)
	}
	query := `SELECT id, name FROM ` + r.DBName + `.users where email = ? and password = ?`
	user := &models.User{
		Email:    email,
		Password: pass,
	}
	err := r.SQLDB.QueryRow(query, email, pass).Scan(&user.ID, &user.Name)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}
		r.logger.Printf("repo query GetUser execution failed: %v", err)
		return nil, err
	}
	return user, nil
}
