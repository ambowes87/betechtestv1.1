package db

import (
	"github.com/ambowes87/betechtestv1.1/pkg/data"
)

type UserStore interface {
	Open() bool
	Ping() error
	Close()
	Add(user data.UserData) error
	Get(id string) (data.UserData, error)
	Update(user data.UserData) error
	Delete(id string) error
}
