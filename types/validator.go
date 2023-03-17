package types

import (
	"context"
	"encoding/json"

	"fmt"
	"log"

	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hectagon-finance/whiteboard/crypto"
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
	AddPeer(peers []interface{})

	// Get port
	Port() int

	// Validate a transaction
	CheckTransaction(tx Transaction) bool

	// broadcast a transaction to all peers
	broadcastTransaction(tx Transaction)
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
	peers        []string
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
	u := url.URL{Scheme: "ws", Host: "localhost:" + strconv.Itoa(v.port), Path: "/ws"}

	// Retry connecting to the server with a delay
	for {
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println("Error connecting to the server, retrying in 1 second:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		defer conn.Close()
		log.Println("Connected to the server")

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				v.checkForAvailableValidators()
				v.findPeer()
			}
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

func (v *validator) ConnectAndSendMessage(message map[string]interface{}) {
	v.clientsMutex.Lock()
	defer v.clientsMutex.Unlock()

	for _, peer := range v.peers {
		if peer != strconv.Itoa(v.port) {
			// Check if the peer is already connected
			isConnected := false
			for conn := range v.clients {
				if conn.RemoteAddr().String() == "localhost:"+peer {
					isConnected = true
					break
				}
			}
			if isConnected {
				continue
			}

			u := url.URL{Scheme: "ws", Host: "localhost:" + peer, Path: "/ws"}
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				fmt.Println("Error connecting to the peer:", err)
				continue
			}
			defer conn.Close()

			v.clients[conn] = true

			v.sendMessage(conn, message)
		}
	}
}

func (v *validator) checkForAvailableValidators() {

	message := map[string]interface{}{
		"type":        "hello",
		"validatorId": v.validatorId,
		"message":     "Hello, I'm validator " + v.validatorId,
	}

	v.ConnectAndSendMessage(message)
}

func (v *validator) findPeer() {

	message := map[string]interface{}{
		"type":        "peer",
		"from":        v.validatorId,
		"validatorId": v.validatorId,
		"message":     v.peers,
	}

	v.ConnectAndSendMessage(message)
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

	switch message["type"].(string) {

	case "hello":
		fmt.Println("Validator", v.validatorId, ": Received message from validator", message["validatorId"].(string), ":", message["message"].(string))
	case "transaction":
		tx := &transaction{
			transactionId: message["transactionId"].(string),
			publicKey:     message["publicKey"].(*crypto.PublicKey),
			timestamp:     int64(message["timestamp"].(float64)),
			signature:     message["signature"].(*crypto.Signature),
			data:          message["data"].([]byte),
		}
		if tx.signature.Verify(*tx.publicKey, tx.data) {
			if v.CheckTransaction(tx) {
				fmt.Printf("Validator %s: Valid transaction received from %s %s: %s\n", v.validatorId, message["from"].(string), message["validatorId"].(string), tx.Id())
				v.broadcastTransaction(tx)
				v.MemPool().AddTransaction(tx)
				fmt.Println(v.MemPool().Size())
			} else {
				fmt.Printf("Validator %s: Already have that transaction\n", v.validatorId)
			}
		} else {
			fmt.Printf("Validator %s: Invalid transaction received from %s %s: %s\n", v.validatorId, message["from"].(string), message["validatorId"].(string), tx.Id())
		}
	case "peer":
		fmt.Printf("Validator %s: Valid peers array received from %s: %s\n", v.validatorId, message["validatorId"].(string), message["message"].([]interface{}))
		v.AddPeer(message["message"].([]interface{}))
	default:
		fmt.Println("Default")
	}
}

func (v *validator) broadcastTransaction(tx Transaction) {
	message := map[string]interface{}{
		"type":          "transaction",
		"from":          "validator",
		"validatorId":   v.validatorId,
		"transactionId": tx.Id(),
		"publicKey":     tx.PublicKey(),
		"timestamp":     tx.Timestamp(),
		"signature":     tx.Signature(),
		"data":          tx.Data(),
	}

	v.ConnectAndSendMessage(message)
}

func (v *validator) CheckTransaction(tx Transaction) bool {
	for i := range v.memPool.GetTransactions() {
		if v.memPool.GetTransactions()[i].Id() == tx.Id() {
			return false
		}
	}
	return true
}

func (v *validator) AddPeer(peers []interface{}) {
	for _, peer := range peers {
		alreadyhave := false
		for _, p := range v.peers {
			if p == peer.(string) {
				alreadyhave = true
				break
			}
		}

		if !alreadyhave {
			fmt.Println("Validator", v.validatorId, ": Adding new peer", peer.(string))
			v.peers = append(v.peers, peer.(string))
		}

	}
}

func (v *validator) Port() int {
	return v.port
}

func NewValidator(port int) Validator {
	// id := rand.Intn(100000000)
	validatorId := strconv.Itoa(port)
	publicKey := "public-key"
	privateKey := "private-key"
	memPool := NewMemPool()

	peers := []string{"8080"}

	if port != 8080 {
		peers = append(peers, strconv.Itoa(port))
	}

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
		peers:       peers,
	}
}
