package db

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/ambowes87/betechtestv1.1/pkg/data"
	"github.com/ambowes87/betechtestv1.1/pkg/logger"
)

const (
	userUpdateString = `UPDATE users SET first_name=?,last_name=?,nickname=?,password=?,email=?,country=?,updated_at=? WHERE id=?`
	userInsertString = `INSERT INTO users(id,first_name,last_name,nickname,password,email,country,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?)`
	userSelectString = `SELECT id,first_name,last_name,nickname,password,email,country,created_at,updated_at FROM users WHERE id=?`
	userDeleteString = `DELETE FROM users WHERE id=?`
)

type UserSQLStore struct {
	fileName string

	database           *sql.DB
	insertStmt         *sql.Stmt
	updateStmt         *sql.Stmt
	selectStmt         *sql.Stmt
	selectMultipleStmt *sql.Stmt
	deleteStmt         *sql.Stmt
}

// NewUserSQLStore creates a new sql type data store
func NewUserSQLStore() *UserSQLStore {
	return &UserSQLStore{
		fileName: "userstore.db",
	}
}

func (u *UserSQLStore) Open() bool {
	var err error
	u.database, err = sql.Open("sqlite3", u.fileName)
	if err != nil {
		logger.Log(err.Error())
		return false
	}
	u.insertStmt, err = u.database.Prepare(userInsertString)
	u.updateStmt, err = u.database.Prepare(userUpdateString)
	u.selectStmt, err = u.database.Prepare(userSelectString)
	u.deleteStmt, err = u.database.Prepare(userDeleteString)
	if err != nil {
		logger.Log(err.Error())
		return false
	}
	return true
}

// Ping checks the database is alive
func (u *UserSQLStore) Ping() error {
	if u.database != nil {
		return u.database.Ping()
	}
	return errors.New("no database initialised")
}

// Close shuts down any prepared queries and then the datastore itself
func (u *UserSQLStore) Close() {
	if u.database != nil {
		u.insertStmt.Close()
		u.updateStmt.Close()
		u.selectStmt.Close()
		u.deleteStmt.Close()
		u.database.Close()
	}
}

// Add adds a user to the database, the user's ID must be unique
func (u *UserSQLStore) Add(user data.UserData) error {
	now := time.Now()
	_, err := u.insertStmt.Exec(user.ID, user.FirstName, user.LastName, user.Nickname, user.Password, user.Email, user.Country, now, now)
	if err != nil {
		return err
	}
	return nil
}

// Get returns a single user matching a unique ID
func (u *UserSQLStore) Get(id string) (data.UserData, error) {
	row := u.selectStmt.QueryRow(id)
	user := data.UserData{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.Password, &user.Email, &user.Country, &user.Created, &user.Updated)
	if err != nil {
		return data.UserData{}, err
	}
	return user, nil
}

// GetPagedByCountry returns a paginated list of users with a specific Country
// CURRENTLY NOT WORKING
func (u *UserSQLStore) GetPagedByCountry(limit, page, country string) ([]data.UserData, error) {

	li, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	pg, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	offset := 0
	if pg > 0 {
		offset = (pg * li) + 1
	}

	sqlString := "SELECT * FROM users WHERE Country=? LIMIT ?,?"
	rows, err := u.database.Query(sqlString, country, limit, offset)
	if err != nil {
		return nil, err
	}

	var users []data.UserData
	for rows.Next() {
		var user data.UserData
		err = rows.Scan(&user)
		if err != nil {
			users = append(users, user)
		}
	}
	if len(users) > 0 {
		return users, err
	}
	return nil, err
}

// Update updates an existing user, the ID must already exist in the database
func (u *UserSQLStore) Update(user data.UserData) error {
	result, err := u.updateStmt.Exec(user.FirstName, user.LastName, user.Nickname, user.Password, user.Email, user.Country, time.Now(), user.ID)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	} else if n == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

// Delete removes a user from the database by ID
func (u *UserSQLStore) Delete(id string) error {
	_, err := u.deleteStmt.Exec(id)
	return err
}
