package persistence

import (
	"database/sql"
	"cod-server/internal/domain"
)

type UserRepository interface {
	Create(id string, user *domain.User) error
	Read(id string) (*domain.User, error)
	Update(id string, user *domain.User) error
	Delete(id string) error
	List() ([]*domain.User, error)
	ListBy(filter func(*domain.User) bool) ([]*domain.User, error)
}

type CardRepository interface {
	Create(id string, card *domain.Card) error
	Read(id string) (*domain.Card, error)
	Update(id string, card *domain.Card) error
	Delete(id string) error
	List() ([]*domain.Card, error)
	ListBy(filter func(*domain.Card) bool) ([]*domain.Card, error)
}

type MatchRepository interface {
	Create(id string, match *domain.Match) error
	Read(id string) (*domain.Match, error)
	Update(id string, match *domain.Match) error
	Delete(id string) error
	List() ([]*domain.Match, error)
	ListBy(filter func(*domain.Match) bool) ([]*domain.Match, error)
}

type RepositoryManager struct {
	DB *sql.DB
}

func NewRepositoryManager(db *sql.DB) *RepositoryManager {
	return &RepositoryManager{DB: db}
}