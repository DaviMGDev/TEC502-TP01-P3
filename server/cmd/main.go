package main

import (
	"cod-server/internal/api"
	"cod-server/internal/api/mqtt"
	"cod-server/internal/auth"
	"cod-server/internal/cluster"
	"cod-server/internal/data/cache"
	"cod-server/internal/data/persistence"
	"cod-server/internal/services"
	"database/sql"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

// getEnv carrega uma variável de ambiente ou retorna um valor padrão se não for encontrada.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// Carrega arquivo .env se existir; avisa mas continua se falhar
	err := godotenv.Load()
	if err != nil {
		log.Warnf("Aviso: Não foi possível carregar o arquivo .env: %v. Usando variáveis de ambiente existentes ou padrões.", err)
	}

	// Carrega configuração de variáveis de ambiente com valores padrão sensatos
	raftDataDir := getEnv("COD_RAFT_DATA_DIR", "./raft-data")
	raftBindAddr := getEnv("COD_RAFT_BIND_ADDR", "127.0.0.1:10000")
	httpBindAddr := getEnv("COD_HTTP_BIND_ADDR", "127.0.0.1:8080")
	nodeID := getEnv("COD_NODE_ID", "node-1")
	mqttBrokerAddr := getEnv("COD_MQTT_BROKER_ADDR", "tcp://localhost:1883")
	isFirstNodeStr := getEnv("COD_IS_FIRST_NODE", "false")
	isFirstNode, _ := strconv.ParseBool(isFirstNodeStr)

	log.Info("Iniciando servidor COD...")

	// Inicializa repositórios de dados e serviços com pool de conexões
	db, err := sql.Open("sqlite3", "./game_data.db")
	if err != nil {
		log.Fatal("falha ao abrir banco de dados SQLite: %v", err)
	}
	defer db.Close()

	// Configura pool de conexões SQLite para acesso concorrente
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	sqlUserRepo := persistence.NewSqlUserRepository(db)
	sqlCardRepo := persistence.NewSqlCardRepository(db)
	sqlMatchRepo := persistence.NewSqlMatchRepository(db)

	userRepo := persistence.NewUserRepoAdapter(sqlUserRepo)
	cardRepo := persistence.NewCardRepoAdapter(sqlCardRepo)
	matchRepo := persistence.NewMatchRepoAdapter(sqlMatchRepo)

	// Envolve repositórios com camada de cache para otimização de desempenho
	userRepo = cache.NewCachedUserRepository(userRepo)
	cardRepo = cache.NewCachedCardRepository(cardRepo)
	matchRepo = cache.NewCachedMatchRepository(matchRepo)

	userService := services.NewUserService(userRepo)
	cardsService := services.NewCardsService(cardRepo, userRepo)
	matchService := services.NewMatchService(matchRepo, cardRepo, userRepo)

	// Inicializa manipulador de eventos da API com serviços e autenticação
	authService := auth.NewAuthService("") // Usa a chave padrão
	eventHandler := api.NewEventHandler(userService, cardsService, matchService, authService)

	// Cria Máquina de Estados Finitos do Raft para gerenciamento de estado distribuído
	fsm := cluster.NewClusterFSM(eventHandler)

	// Configura e inicializa consenso Raft com transporte TCP
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	log.SetLevel(log.DebugLevel)
	var raftLogger io.Writer = log.StandardLog(log.StandardLogOptions{ForceLevel: log.DebugLevel}).Writer()

	addr, err := net.ResolveTCPAddr("tcp", raftBindAddr)
	if err != nil {
		log.Fatal("falha ao resolver endereço TCP do Raft: %v", err)
	}
	transport, err := raft.NewTCPTransport(raftBindAddr, addr, 3, 10*time.Second, raftLogger)
	if err != nil {
		log.Fatal("falha ao criar transporte TCP do Raft: %v", err)
	}

	if err := os.MkdirAll(raftDataDir, 0755); err != nil {
		log.Fatal("falha ao criar diretório de dados do Raft: %v", err)
	}

	logStore, err := raftboltdb.NewBoltStore(raftDataDir + "/logs.db")
	if err != nil {
		log.Fatal("falha ao criar log store: %v", err)
	}
	stableStore, err := raftboltdb.NewBoltStore(raftDataDir + "/stable.db")
	if err != nil {
		log.Fatal("falha ao criar stable store: %v", err)
	}
	snapshotStore, err := raft.NewFileSnapshotStore(raftDataDir, 1, raftLogger)
	if err != nil {
		log.Fatal("falha ao criar snapshot store: %v", err)
	}

	raftNode, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		log.Fatal("falha ao criar nó Raft: %v", err)
	}

	// Bootstrap (apenas para o primeiro nó)
	if isFirstNode {
		log.Info("Realizando bootstrap do cluster...")
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		if err := raftNode.BootstrapCluster(configuration).Error(); err != nil {
			log.Fatal("falha ao realizar bootstrap do cluster: %v", err)
		}
	}

	// Inicializa transporte HTTP da API para comunicação entre nós
	httpTransport := cluster.NewGinHttpTransport(httpBindAddr, nodeID, raftNode)
	if err := httpTransport.Start(); err != nil {
		log.Fatal("Falha ao iniciar transporte HTTP: %v", err)
	}

	// Configura adaptador MQTT para comunicação de eventos do cliente
	mqttAdapter, err := mqtt.NewMQTTAdapter(mqttBrokerAddr, nodeID)
	if err != nil {
		log.Fatal("Falha ao criar adaptador MQTT: %v", err)
	}
	if err := mqttAdapter.Connect(); err != nil {
		log.Fatal("Falha ao conectar ao broker MQTT: %v", err)
	}
	defer mqttAdapter.Disconnect()

	// Cria coordenador Raft para gerenciar roteamento de eventos e consenso
	coordinator := cluster.NewRaftCoordinator(raftNode, httpTransport, mqttAdapter)

	// Inscreve-se em todos os tópicos de eventos do cliente
	// Tópicos correspondem aos definidos no EventService do cliente
	tópicos := []string{
		"user/register",
		"user/login",
		"chat/room/+", // Usando + como wildcard para qualquer room
		"game/start_game",
		"game/+/play_card", // Wildcard para room específico
		"game/+/surrender", // Wildcard para room específico
		"game/join_game",
		"store/buy",
		"cards/+/exchange+", // Wildcard para room e user
		"game/actions",      // Tópico original
	}

	for _, tópico := range tópicos {
		mqttAdapter.Subscribe(tópico, func(client paho.Client, msg paho.Message) {
			event, err := api.FromJson(msg.Payload())
			if err != nil {
				log.Errorf("Erro ao desserializar evento MQTT: %v", err)
				return
			}
			log.Infof("Evento MQTT recebido no tópico %s: %+v", msg.Topic(), event)
			if err := coordinator.Handle(*event); err != nil {
				log.Errorf("Erro ao processar evento via coordenador: %v", err)
			}
		})
	}

	// Inicializa serviço de descoberta de pares para associação automática ao cluster
	discovery := cluster.NewDiscoveryService(string(transport.LocalAddr()), httpBindAddr)
	discovery.OnPeerDiscovered = func(peerIp string) {
		targetAddr := fmt.Sprintf("%s:%s", peerIp, "8080") // Assumindo porta 8080
		log.Infof("Nó par descoberto em %s. Tentando adicionar ao cluster...", targetAddr)

		// Lógica para evitar adicionar nós duplicados ou a si mesmo
		if raftNode.State() != raft.Leader {
			log.Info("Não sou o líder, não posso adicionar um novo nó.")
			return
		}

		future := raftNode.GetConfiguration()
		if err := future.Error(); err != nil {
			log.Errorf("Falha ao obter configuração do cluster: %v", err)
			return
		}
		for _, srv := range future.Configuration().Servers {
			if srv.Address == raft.ServerAddress(targetAddr) {
				log.Infof("Nó %s já faz parte do cluster.", targetAddr)
				return
			}
		}

		// Adicionando o nó descoberto
		addVoterFuture := raftNode.AddVoter(raft.ServerID(peerIp), raft.ServerAddress(targetAddr), 0, 0)
		if err := addVoterFuture.Error(); err != nil {
			log.Errorf("Falha ao adicionar novo nó ao cluster: %v", err)
		} else {
			log.Infof("Nó %s adicionado ao cluster.", targetAddr)
		}
	}
	discovery.Start()

	// Bloqueia até que o sinal de desligamento (Ctrl-C) seja recebido, então desliga graciosamente
	log.Info("Servidor COD rodando. Pressione CTRL-C para sair.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Info("Desligando o servidor...")
}
