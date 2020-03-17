package auth

import (
	"fmt"
)

type UsersDataStorage struct {
	data       map[string]*User
	lastUserId uint32
}

func NewUsersDataStorage() *UsersDataStorage {
	return &UsersDataStorage{
		data: make(map[string]*User),
	}
}

func (ud *UsersDataStorage) Login(login string) (*User, error) {
	user, exist := ud.data[login]
	if !exist {
		return nil, fmt.Errorf("login '%s' doesnot exist", login)
	}
	return user, nil
}

func (ud *UsersDataStorage) Register(login string) (*User, error) {
	if _, exist := ud.data[login]; exist {
		return nil, fmt.Errorf("login '%s' already existed", login)
	}
	u := &User{
		Id:    ud.lastUserId,
		Login: login,
	}
	ud.data[login] = u
	ud.lastUserId++

	return u, nil
}

type User struct {
	Id    uint32
	Login string
}
