package protocol

// Event é uma reexportação do tipo Event do protocolo compartilhado para compatibilidade retroativa dentro do pacote cliente.

import (
	shared_protocol "shared/protocol"
)

// Reexporta o tipo Event compartilhado para compatibilidade retroativa
type Event = shared_protocol.Event
