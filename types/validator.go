package types

import (
	"context"
	"encoding/json"
	// "flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"


	"github.com/gorilla/websocket"
)

type Validator interface {
	// Get the validator's id
	Id() string

	// Get the validator's public key
	PublicKey() string

	// Get the validator's private key
	PrivateKey() string

	// Get the validator's mempool
	MemPool() MemPool

	// Get the validator's balance
	Balance() int64

	// Get the validator's stake
	Stake() int64


	// Get the validator's status
	Status() string

	// Get the validator's last block
	LastBlock() Block

	// Get the validator's last block hash
	LastBlockHash() string

	// Start the validator
	Start()

	// Stop the validator
	Stop()

	// Start validator's client
	StartClient()

	// Stop validator's client
	StopClient()

	// Start validator's server
	StartServer(ready chan bool)

	// Stop validator's server
	StopServer()

	// Add a new peer to the validator
	AddPeer(peer string)

	// Get port
	Port() int
}

type validator struct {
	validatorId  string
	publicKey    string
	privateKey   string
	memPool      MemPool
	balance      int64
	stake        int64
	status       string
	lastBlock    Block
	port         int
	httpServer   *http.Server
	clientsMutex sync.Mutex
	clients      map[*websocket.Conn]bool
	peers 	  	 []string
	stopServer   func()
}


func (v *validator) Id() string {
	return v.validatorId
}

func (v *validator) PublicKey() string {
	return v.publicKey
}

func (v *validator) PrivateKey() string {
	return v.privateKey
}

func (v *validator) MemPool() MemPool {
	return v.memPool
}

func (v *validator) Balance() int64 {
	return v.balance
}

func (v *validator) Stake() int64 {
	return v.stake
}

func (v *validator) Status() string {
	return v.status
}

func (v *validator) LastBlock() Block {
	return v.lastBlock
}

func (v *validator) LastBlockHash() string {
	return v.lastBlock.Hash()
}

func (v *validator) Start() {
	v.status = "active"
	ready := make(chan bool)
	go v.StartServer(ready)
	<-ready
	go v.StartClient()
}

func (v *validator) Stop() {
	v.StopServer()
	v.StopClient()
	v.status = "inactive"
}

func (v *validator) StartClient() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:"+strconv.Itoa(v.port) +"/ws", nil)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	log.Println("Connected to the server")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			v.checkForAvailableValidators()
		}
	}
}

func (v *validator) StopClient() {
	v.clientsMutex.Lock()
	for conn := range v.clients {
		conn.Close()
		delete(v.clients, conn)
	}
	v.clientsMutex.Unlock()
}

func (v *validator) StartServer(ready chan bool) {
	// Create a new Gorilla WebSocket upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Create a new http.ServeMux for this validator
	mux := http.NewServeMux()

	// Set up the WebSocket endpoint for this validator
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading HTTP to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		// Add the connection to the clients map
		v.clientsMutex.Lock()
		v.clients[conn] = true
		v.clientsMutex.Unlock()

		// Handle incoming messages from this connection
		v.handleConnections(conn)
	})

	// Create an HTTP server
	v.httpServer = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", v.port),
		Handler: mux,
	}

	go func() {
		// Start the HTTP server
		log.Printf("Validator %s starting server on port %d", v.validatorId, v.port)
		ready <- true
		if err := v.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %v", err)
		}
	}()

	v.stopServer = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := v.httpServer.Shutdown(ctx); err != nil {
			log.Printf("Error stopping server: %v", err)
		}
	}
}
func (v *validator) handleConnections(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		v.handleMessage(msg)
	}
}


func (v *validator) StopServer() {
	if v.stopServer != nil {
		log.Printf("Validator %s stopping server on port %d", v.validatorId, v.port)
		v.stopServer()
	}
}

func (v *validator) checkForAvailableValidators() {
	for _, peer := range v.peers {
		conn, _, err := websocket.DefaultDialer.Dial("ws://"+peer, nil)
		if err != nil {
			fmt.Println("Error connecting to the peer:", err)
			continue
		}
		defer conn.Close()

		v.clientsMutex.Lock()
		v.clients[conn] = true
		v.clientsMutex.Unlock()

		message := map[string]interface{}{
			"validatorId": v.validatorId,
			"message":     "Hello, I'm validator " + v.validatorId,
		}
		v.sendMessage(conn, message)
	}
}

func (v *validator) sendMessage(conn *websocket.Conn, message map[string]interface{}) {
	msg, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling the message:", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, msg)
}

func (v *validator) handleMessage(msg []byte) {
	var message map[string]interface{}
	err := json.Unmarshal(msg, &message)
	if err != nil {
		fmt.Println("Error unmarshaling the message:", err)
		return
	}
	fmt.Println("Validator",v.validatorId,": Received message from validator", message["validatorId"].(string), ":", message["message"].(string))
}

func NewValidator(port int) Validator {
	id := rand.Intn(1000)
	validatorId := strconv.Itoa(id)
	publicKey := "public-key"
	privateKey := "private-key"
	memPool := NewMemPool()

	return &validator{
		validatorId: validatorId,
		publicKey:   publicKey,
		privateKey:  privateKey,
		memPool:     memPool,
		balance:     0,
		stake:       0,
		status:      "inactive",
		port:        port,
		clients:     make(map[*websocket.Conn]bool),
		peers:       []string{},
	}
}

func (v *validator) AddPeer(peer string) {
	v.peers = append(v.peers, peer)
}

func (v *validator) Port() int {
	return v.port
}