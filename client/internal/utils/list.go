package utils 

type List[T comparable] interface {
	Append(item T)
	Prepend(item T)
	Insert(index int, item T) error
	Remove(item T) error
	Pop() (T, error)
	Shift() (T, error)
	FindIndex(predicate func(T) bool) (int, bool)
	Find(predicate func(T) bool) (T, bool)
	Size() int
	Items() []T
	Get(index int) (T, error)
	Set(index int, item T) error
}
