package utils

// Mux é um roteador/multiplexador genérico de comandos que mapeia nomes de comandos para seus manipuladores.
// Usa um mapa para busca de comando em O(1) e fornece um manipulador padrão para comandos desconhecidos.
type Mux[fn any] struct {
	table          map[string]fn
	defaultHandler fn
}

// NewMux cria uma nova instância de Mux com o manipulador padrão fornecido para comandos não mapeados.
func NewMux[fn any](defaultHandler fn) *Mux[fn] {
	return &Mux[fn]{
		table:          make(map[string]fn),
		defaultHandler: defaultHandler,
	}
}

// Register mapeia um nome de comando para sua função manipuladora.
func (mux *Mux[fn]) Register(command string, handler fn) {
	mux.table[command] = handler
}

// Handle retorna o manipulador para o comando fornecido, ou o manipulador padrão se não encontrado.
func (mux *Mux[fn]) Handle(command string) fn {
	handler, exists := mux.table[command]
	if !exists {
		return mux.defaultHandler
	}
	return handler
}
