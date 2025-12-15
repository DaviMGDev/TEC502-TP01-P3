package utils

// List é uma interface genérica que define operações de coleção ordenada.
// Implementações devem suportar adicionar, remover e buscar elementos.
type List[T comparable] interface {
	// Append adiciona um item ao final da lista.
	Append(item T)
	// Prepend adiciona um item ao início da lista.
	Prepend(item T)
	// Insert adiciona um item no índice especificado.
	Insert(index int, item T) error
	// Remove exclui a primeira ocorrência do item especificado.
	Remove(item T) error
	// Pop remove e retorna o último item da lista.
	Pop() (T, error)
	// Shift remove e retorna o primeiro item da lista.
	Shift() (T, error)
	// FindIndex retorna o índice do primeiro item que corresponde ao predicado.
	FindIndex(predicate func(T) bool) (int, bool)
	// Find retorna o primeiro item que corresponde ao predicado.
	Find(predicate func(T) bool) (T, bool)
	// Size retorna o número de items na lista.
	Size() int
	// Items retorna uma cópia de todos os items como um slice.
	Items() []T
	// Get recupera o item no índice especificado.
	Get(index int) (T, error)
	// Set substitui o item no índice especificado.
	Set(index int, item T) error
}
