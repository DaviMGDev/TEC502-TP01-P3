package data

// import "cod-server/internal/utils"

type Repository[T any] interface {
	Create(id string, entity T) error
	Read(id string) (T, error)
	Update(id string, entity T) error
	Delete(id string) error
	List() ([]T, error)
	ListBy(filter func(T) bool) ([]T, error)
}
