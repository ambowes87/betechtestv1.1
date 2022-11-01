package db

import (
	"database/sql"
	"errors"
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

	database   *sql.DB
	insertStmt *sql.Stmt
	updateStmt *sql.Stmt
	selectStmt *sql.Stmt
	deleteStmt *sql.Stmt
}

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

func (u *UserSQLStore) Ping() error {
	if u.database != nil {
		return u.database.Ping()
	}
	return errors.New("no database initialised")
}

func (u *UserSQLStore) Close() {
	if u.database != nil {
		u.insertStmt.Close()
		u.updateStmt.Close()
		u.selectStmt.Close()
		u.deleteStmt.Close()
		u.database.Close()
	}
}

func (u *UserSQLStore) Add(user data.UserData) error {
	now := time.Now()
	_, err := u.insertStmt.Exec(user.ID, user.FirstName, user.LastName, user.Nickname, user.Password, user.Email, user.Country, now, now)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserSQLStore) Get(id string) (data.UserData, error) {
	row := u.selectStmt.QueryRow(id)
	user := data.UserData{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.Password, &user.Email, &user.Country, &user.Created, &user.Updated)
	if err != nil {
		return data.UserData{}, err
	}
	return user, nil
}

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

func (u *UserSQLStore) Delete(id string) error {
	_, err := u.deleteStmt.Exec(id)
	return err
}
