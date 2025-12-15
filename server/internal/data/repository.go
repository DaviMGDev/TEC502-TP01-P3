package data

// import "cod-server/internal/utils"

// Repository é uma interface genérica de CRUD para entidades de domínio.
// Implementações podem ser em memória, com SQL ou adaptadores com cache.
type Repository[T any] interface {
	// Create insere uma nova entidade com o id fornecido.
	Create(id string, entity T) error
	// Read busca uma entidade pelo id, retornando erro se não encontrada.
	Read(id string) (T, error)
	// Update substitui a entidade armazenada sob o id fornecido.
	Update(id string, entity T) error
	// Delete remove a entidade com o id fornecido.
	Delete(id string) error
	// List retorna todas as entidades no repositório.
	List() ([]T, error)
	// ListBy retorna entidades que satisfaçam o predicado de filtro fornecido.
	ListBy(filter func(T) bool) ([]T, error)
}
