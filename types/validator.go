package types

import (
	"math/rand"
	"encoding/json"
	"fmt"
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
	StartServer()

	// Stop validator's server
	StopServer()
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
	go v.StartServer()
	go v.StartClient()
}

func (v *validator) Stop() {
	v.StopServer()
	v.StopClient()
	v.status = "inactive"
}

func (v *validator) StartClient() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:"+strconv.Itoa(v.port), nil)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

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

func (v *validator) StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", v.handleConnections)
	v.httpServer = &http.Server{Addr: ":" + strconv.Itoa(v.port), Handler: mux}
	err := v.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println("Error starting the server:", err)
	}
}

func (v *validator) StopServer() {
	// write code in here
	if err := v.httpServer.Close(); err != nil {
		fmt.Println("Error stopping the server:", err)
	}
}

func (v *validator) handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading the connection:", err)
		return
	}
	defer conn.Close()

	v.clientsMutex.Lock()
	v.clients[conn] = true
	v.clientsMutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(v.clients, conn)
			break
		}
		v.handleMessage(msg)
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
	fmt.Println(v.validatorId,": Received message from validator", message["validatorId"].(string), ":", message["message"].(string))
}

func NewValidator(port int) Validator {
	id := rand.Intn(1000)
	validatorId := "test-validator-"+strconv.Itoa(id)
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