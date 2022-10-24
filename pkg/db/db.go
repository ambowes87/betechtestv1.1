package db

import (
	"errors"

	"github.com/ambowes87/betechtestv1.1/pkg/data"
)

func AddUser(user data.UserData) error {
	_, ok := fakeUsers[user.ID]
	if ok {
		return errors.New("a user already exists with ID " + user.ID)
	}
	fakeUsers[user.ID] = user
	return nil
}

func GetUser(id string) (data.UserData, error) {
	u, ok := fakeUsers[id]
	if !ok {
		return u, errors.New("could not find user with id " + id)
	}
	return u, nil
}

func UpdateUser(user data.UserData) error {
	_, ok := fakeUsers[user.ID]
	if !ok {
		return errors.New("could not find user with id " + user.ID)
	}
	fakeUsers[user.ID] = user
	return nil
}

func DeleteUser(id string) error {
	_, ok := fakeUsers[id]
	if !ok {
		return errors.New("could not find user with id " + id)
	}
	delete(fakeUsers, id)
	return nil
}
