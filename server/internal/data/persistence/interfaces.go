package persistence

// SQL-backed repository interfaces for domain models with a manager for resource sharing.
// UserRepository defines CRUD operations for users backed by a SQL database.
// CardRepository defines CRUD operations for cards backed by a SQL database.
// MatchRepository defines CRUD operations for matches backed by a SQL database.
// RepositoryManager holds shared SQL DB connections for repository implementations.
// NewRepositoryManager constructs a manager for sharing DB connections across repos.

import (
	"cod-server/internal/domain"
	"database/sql"
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
