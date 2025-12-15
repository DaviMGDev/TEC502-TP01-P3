package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	GetID() string
	GetUsername() string
	CheckPassword(password string) bool
}

type User struct {
	ID       string        `json:"id"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Cards    PackInterface `json:"cards"`
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
