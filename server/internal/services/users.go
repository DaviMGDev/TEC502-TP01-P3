package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"golang.org/x/crypto/bcrypt"

	uuid "github.com/google/uuid"
)

type UserService struct {
	userRepo data.Repository[domain.UserInterface]
}

func NewUserService(userRepo data.Repository[domain.UserInterface]) UserServiceInterface {
	return &UserService{userRepo: userRepo}
}

func (us *UserService) Register(username, password string) error {
	id := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &domain.User{
		ID:       id,
		Username: username,
		Password: string(hashedPassword),
		Cards:    nil,
	}
	return us.userRepo.Create(id, user)
}

func (us *UserService) Login(username, password string) (*domain.UserInterface, error) {
	users, err := us.userRepo.ListBy(func(u domain.UserInterface) bool {
		return u.GetUsername() == username && u.CheckPassword(password)
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil 
}
