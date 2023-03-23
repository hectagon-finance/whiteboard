package start

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
	. "github.com/hectagon-finance/whiteboard/validator"
)

func Start(v *Validator) {
	v.Status = "active"
	ready := make(chan bool)
	go StartServer(v, ready)
	<-ready
	go StartClient(v)
}

func StartServer(v *Validator, ready chan bool) {
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
		v.ClientsMutex.Lock()
		v.Clients[conn] = true
		v.ClientsMutex.Unlock()

		// Handle incoming messages from this connection
		handleConnections(v, conn)
	})

	// Create an HTTP server
	v.HttpServer = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", v.Port),
		Handler: mux,
	}

	go func() {
		// Start the HTTP server
		log.Printf("Validator %s starting server on port %d", v.ValidatorId, v.Port)
		ready <- true
		if err := v.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %v", err)
		}
	}()

	v.StopServer = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := v.HttpServer.Shutdown(ctx); err != nil {
			log.Printf("Error stopping server: %v", err)
		}
	}
}

func handleConnections(v *Validator, conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		HandleMessage(v, msg)
	}
}

func StartClient(v *Validator) {
	u := url.URL{Scheme: "ws", Host: "localhost:" + strconv.Itoa(v.Port), Path: "/ws"}

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

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// BroadcastMempool(v)
				BroadcastPeer(v)
				checkMemPool(v)
			}
		}
	}
}

func checkMemPool(v *Validator) {
	if v.MemPool.Size() == 1 {
		v.TempBlock = NewBlock(1, [32]byte{}, v.MemPool.GetTransactions())
		blockHash := v.TempBlock.GetHash()
		fmt.Println("Block hash:", blockHash)
		BroadcastBlockHash(v, blockHash)
		v.MemPool.Clear()
	}
}
