package utils

// Map é uma interface genérica para operações de armazenamento chave-valor.
// K deve ser comparável; V pode ser qualquer tipo.
type Map[K comparable, V any] interface {
	// Get recupera um valor pela chave, retornando false se a chave não existir.
	Get(key K) (V, bool)
	// Set armazena ou atualiza um valor para uma chave.
	Set(key K, value V) error
	// Delete remove um par chave-valor do mapa.
	Delete(key K) error
	// Keys retorna um slice de todas as chaves no mapa.
	Keys() []K
	// Values retorna um slice de todos os valores no mapa.
	Values() []V
	// Size retorna o número de pares chave-valor no mapa.
	Size() int
}
