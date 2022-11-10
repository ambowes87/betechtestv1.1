package db

import (
	"errors"

	"github.com/ambowes87/betechtestv1.1/pkg/data"
)

type UserMemStore struct {
	fakeUsers map[string]data.UserData
}

func NewUserMemStore() *UserMemStore {
	return &UserMemStore{
		fakeUsers: make(map[string]data.UserData),
	}
}

func (u *UserMemStore) Open() bool {
	return u.fakeUsers != nil
}

func (u *UserMemStore) Ping() error {
	return nil
}

func (u *UserMemStore) Close() {
	u.fakeUsers = nil
}

func (u *UserMemStore) Add(user data.UserData) error {
	_, ok := u.fakeUsers[user.ID]
	if ok {
		return errors.New("a user already exists with ID " + user.ID)
	}
	u.fakeUsers[user.ID] = user
	return nil
}

func (u *UserMemStore) Get(id string) (data.UserData, error) {
	user, ok := u.fakeUsers[id]
	if !ok {
		return user, errors.New("could not find user with id " + id)
	}
	return user, nil
}

func (u *UserMemStore) GetPagedByCountry(limit, page, country string) ([]data.UserData, error) {
	return nil, errors.New("not implemented")
}

func (u *UserMemStore) Update(user data.UserData) error {
	_, ok := u.fakeUsers[user.ID]
	if !ok {
		return errors.New("could not find user with id " + user.ID)
	}
	u.fakeUsers[user.ID] = user
	return nil
}

func (u *UserMemStore) Delete(id string) error {
	_, ok := u.fakeUsers[id]
	if !ok {
		return errors.New("could not find user with id " + id)
	}
	delete(u.fakeUsers, id)
	return nil
}
