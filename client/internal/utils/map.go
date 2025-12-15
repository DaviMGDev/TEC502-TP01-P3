package utils

type Map[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V) error
	Delete(key K) error
	Keys() []K
	Values() []V
	Size() int
}
