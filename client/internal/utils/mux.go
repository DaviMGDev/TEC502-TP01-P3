package utils

type Mux[fn any] struct {
	table map[string]fn
	defaultHandler fn
}

func NewMux[fn any](defaultHandler fn) *Mux[fn] {
	return &Mux[fn]{
		table: make(map[string]fn),
		defaultHandler: defaultHandler,
	}
}

func (mux *Mux[fn]) Register(command string, handler fn) {
	mux.table[command] = handler
}

func (mux *Mux[fn]) Handle(command string) fn {
	handler, exists := mux.table[command]
	if !exists {
		return mux.defaultHandler
	}
	return handler
}
