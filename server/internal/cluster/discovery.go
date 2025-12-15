package cluster

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

const (
	DiscoveryPort     = 9999
	DiscoveryInterval = 5 * time.Second
	DiscoveryMessage  = "COD_SERVER_DISCOVERY"
)

// DiscoveryServiceInterface defines the contract for the peer discovery service.
type DiscoveryServiceInterface interface {
	Start()
}

// DiscoveryService manages automatic peer discovery via UDP broadcast and optional HTTP checks.
type DiscoveryService struct {
	raftAddress      string
	httpAddress      string
	Port             int
	Interval         time.Duration
	Message          string
	knownPeers       []string
	OnPeerDiscovered func(peerRaftAddress string)
	logger           *log.Logger
}

// NewDiscoveryService creates a new DiscoveryService with default intervals and message signature.
func NewDiscoveryService(raftAddress, httpAddress string) *DiscoveryService {
	logger := log.With("component", "discovery")
	return &DiscoveryService{
		raftAddress: raftAddress,
		httpAddress: httpAddress,
		Port:        DiscoveryPort,
		Interval:    DiscoveryInterval,
		Message:     DiscoveryMessage,
		knownPeers:  make([]string, 0),
		logger:      logger,
	}
}

// Start launches background goroutines for listening, broadcasting, and periodic checks.
func (ds *DiscoveryService) Start() {
	go ds.listen()
	go ds.broadcast()
	go ds.periodicPeerCheck() // Verifica periodicamente os nós conhecidos via HTTP
}

// addKnownPeer adds a peer to the known peers list if not already present.
func (ds *DiscoveryService) addKnownPeer(peerRaftAddress string) {
	for _, knownPeer := range ds.knownPeers {
		if knownPeer == peerRaftAddress {
			return // Já conhecido
		}
	}
	ds.knownPeers = append(ds.knownPeers, peerRaftAddress)
}

// discoverViaHTTP attempts to discover peers by probing an HTTP discovery endpoint.
func (ds *DiscoveryService) discoverViaHTTP(targetAddress string) error {
	// Extrair IP do endereço Raft
	parts := strings.Split(targetAddress, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid raft address format: %s", targetAddress)
	}
	ip := parts[0]
	httpAddr := fmt.Sprintf("http://%s:8080", ip) // Assumindo porta 8080 para HTTP

	resp, err := http.Get(httpAddr + "/raft/discovery") // Novo endpoint para descoberta
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Em uma implementação completa, processaríamos a resposta
		// Por enquanto, apenas logamos a descoberta
		ds.logger.Infof("Discovered peer via HTTP: %s", targetAddress)
		ds.addKnownPeer(targetAddress)
		if ds.OnPeerDiscovered != nil {
			go ds.OnPeerDiscovered(targetAddress)
		}
	}

	return nil
}

// periodicPeerCheck periodically verifies known peers; placeholder for future liveness checks.
func (ds *DiscoveryService) periodicPeerCheck() {
	ticker := time.NewTicker(10 * time.Second) // Verifica a cada 10 segundos
	defer ticker.Stop()

	for {
		<-ticker.C
		// Aqui faríamos verificações nos nós conhecidos
		// Por exemplo, chamadas HTTP para confirmar que estão ativos
		ds.logger.Debug("Performing periodic peer check", "known_peers_count", len(ds.knownPeers))
	}
}

func (ds *DiscoveryService) listen() {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", ds.Port))
	if err != nil {
		ds.logger.Fatal("Falha ao resolver endereço UDP para escuta", "err", err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		ds.logger.Fatal("Falha ao escutar na porta UDP", "err", err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			ds.logger.Warn("Erro ao ler do UDP", "err", err)
			continue
		}

		message := string(buf[:n])
		// Extrai o endereço do Raft da mensagem, que é o payload
		if len(message) > len(ds.Message) && message[:len(ds.Message)] == ds.Message {
			peerRaftAddr := message[len(ds.Message):]

			// Não reagir às próprias mensagens
			if peerRaftAddr == ds.raftAddress {
				continue
			}

			ds.logger.Infof("Nó par descoberto com endereço Raft: %s", peerRaftAddr)

			if ds.OnPeerDiscovered != nil {
				go ds.OnPeerDiscovered(peerRaftAddr)
			}
		}
	}
}

func (ds *DiscoveryService) broadcast() {
	broadcastAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", ds.Port))
	if err != nil {
		ds.logger.Fatal("Falha ao resolver endereço de broadcast", "err", err)
	}

	conn, err := net.DialUDP("udp4", nil, broadcastAddr)
	if err != nil {
		ds.logger.Fatal("Falha ao conectar para broadcast", "err", err)
	}
	defer conn.Close()

	ticker := time.NewTicker(ds.Interval)
	defer ticker.Stop()

	// Mensagem a ser enviada = MagicString + nosso endereço Raft
	message := ds.Message + ds.raftAddress

	for {
		<-ticker.C
		_, err := conn.Write([]byte(message))
		if err != nil {
			ds.logger.Warn("Falha ao enviar broadcast", "err", err)
		} else {
			// ds.logger.Debug("Broadcast enviado.")
		}
	}
}
