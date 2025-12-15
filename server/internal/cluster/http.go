package cluster

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/raft"
)

// GinHttpTransport implementa a ClusterTransportInterface
type GinHttpTransport struct {
	bindAddress string
	nodeID      string
	router      *gin.Engine
	client      *resty.Client
	raftNode    *raft.Raft
	timeout     time.Duration
	logger      *log.Logger
}

// NewGinHttpTransport cria a instância do transporte
func NewGinHttpTransport(bindAddress, nodeID string, raftNode *raft.Raft) ClusterTransportInterface {
	logger := log.With("component", "http-transport")

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.StandardLog(log.StandardLogOptions{ForceLevel: log.InfoLevel}).Writer()
	gin.DefaultErrorWriter = log.StandardLog(log.StandardLogOptions{ForceLevel: log.ErrorLevel}).Writer()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))

	transport := &GinHttpTransport{
		bindAddress: bindAddress,
		nodeID:      nodeID,
		router:      router,
		client:      resty.New(),
		raftNode:    raftNode,
		timeout:     10 * time.Second,
		logger:      logger,
	}
	transport.setupRoutes()
	return transport
}

// LoggerMiddleware é um middleware Gin para logging estruturado com charmbracelet/log
func LoggerMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		var logFn func(msg interface{}, keyvals ...interface{})

		if c.Writer.Status() >= http.StatusInternalServerError {
			logFn = logger.Error
		} else if c.Writer.Status() >= http.StatusBadRequest {
			logFn = logger.Warn
		} else {
			logFn = logger.Info
		}

		logFn("Request",
			"status", param.StatusCode,
			"method", param.Method,
			"path", param.Path,
			"latency", param.Latency,
			"ip", param.ClientIP,
			"user-agent", c.Request.UserAgent(),
			"errors", param.ErrorMessage,
		)
	}
}

func (t *GinHttpTransport) setupRoutes() {
	group := t.router.Group("/raft")
	group.POST("/join", t.handleJoin)
	group.POST("/command", t.handleCommand)
}

// Start inicia o servidor HTTP (Gin) em background
func (t *GinHttpTransport) Start() error {
	t.logger.Infof("Iniciando servidor HTTP em %s", t.bindAddress)
	go func() {
		if err := t.router.Run(t.bindAddress); err != nil {
			t.logger.Fatal("Falha ao iniciar o servidor HTTP", "err", err)
		}
	}()
	return nil
}

// JoinCluster é usado por um nó novo para pedir entrada no cluster
func (t *GinHttpTransport) JoinCluster(targetAddress string, myRaftID string, myRaftAddress string) error {
	req := JoinRequest{
		NodeID:      myRaftID,
		NodeAddress: myRaftAddress,
	}

	t.logger.Infof("Enviando requisição de join para %s", targetAddress)
	resp, err := t.client.R().
		SetBody(req).
		Post(fmt.Sprintf("http://%s/raft/join", targetAddress))

	if err != nil {
		return fmt.Errorf("falha ao enviar requisição de join para %s: %w", targetAddress, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("erro ao entrar no cluster. status: %s, body: %s", resp.Status(), resp.String())
	}
	t.logger.Infof("Join bem-sucedido no nó %s", targetAddress)
	return nil
}

// ForwardCommand é usado por um nó seguidor para repassar um evento ao líder
func (t *GinHttpTransport) ForwardCommand(leaderAddress string, eventBytes []byte) error {
	t.logger.Debugf("Encaminhando comando para o líder em %s", leaderAddress)
	resp, err := t.client.R().
		SetBody(bytes.NewReader(eventBytes)).
		SetHeader("Content-Type", "application/json").
		Post(fmt.Sprintf("http://%s/raft/command", leaderAddress))

	if err != nil {
		return fmt.Errorf("falha ao encaminhar comando para o líder %s: %w", leaderAddress, err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("erro do líder ao processar comando. status: %s, body: %s", resp.Status(), resp.String())
	}
	return nil
}

// Handlers do Gin (privados)
func (t *GinHttpTransport) handleJoin(c *gin.Context) {
	var req JoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido: " + err.Error()})
		return
	}

	t.logger.Infof("Recebido pedido de join de nó %s em %s", req.NodeID, req.NodeAddress)

	if t.raftNode.State() != raft.Leader {
		t.logger.Warn("Recebido pedido de join, mas não sou o líder")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "não sou o líder"})
		return
	}

	addVoterFuture := t.raftNode.AddVoter(raft.ServerID(req.NodeID), raft.ServerAddress(req.NodeAddress), 0, 0)
	if err := addVoterFuture.Error(); err != nil {
		t.logger.Error("Falha ao adicionar nó ao cluster", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao adicionar nó ao cluster: " + err.Error()})
		return
	}

	t.logger.Infof("Nó %s em %s adicionado ao cluster com sucesso", req.NodeID, req.NodeAddress)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (t *GinHttpTransport) handleCommand(c *gin.Context) {
	if t.raftNode.State() != raft.Leader {
		t.logger.Warn("Recebido comando para aplicar, mas não sou o líder")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "não sou o líder"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		t.logger.Error("Falha ao ler corpo da requisição de comando", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "falha ao ler corpo da requisição: " + err.Error()})
		return
	}

	t.logger.Debug("Aplicando comando recebido via HTTP")
	applyFuture := t.raftNode.Apply(body, t.timeout)
	if err := applyFuture.Error(); err != nil {
		t.logger.Error("Falha ao aplicar comando no raft", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao aplicar comando no raft: " + err.Error()})
		return
	}

	// Opcional: retornar a resposta da FSM para o nó seguidor
	res := applyFuture.Response()
	if resErr, ok := res.(error); ok {
		t.logger.Error("Comando aplicado com erro na FSM", "err", resErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": resErr.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}