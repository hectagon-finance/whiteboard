package types

import (
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestValidatorStartStop(t *testing.T) {

	// Initialize a new validator
	v := NewValidator(8080)

	// Start the validator and wait for a short duration
	v.Start()
	time.Sleep(3 * time.Second)

	// Stop the validator
	v.Stop()

	// If no errors or panics occurred during the test, the test passed
}

func TestValidatorSendMessage(t *testing.T) {
	// Create a dummy MemPool and Block for testing purposes
	memPool := NewMemPool()
	// block := NewBlock()

	// Initialize two new validators
	v1 := &validator{
		validatorId: "validator-1",
		publicKey:   "public-key-1",
		privateKey:  "private-key-1",
		memPool:     memPool,
		balance:     1000,
		stake:       100,
		status:      "inactive",
		port:        8081,
		clients:     make(map[*websocket.Conn]bool),
		peers:       []string{},
	}

	// Use a wait group to make sure the validator is started before connecting the test client
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the validator and wait for it to be ready
	go func() {
		v1.Start()
		wg.Done()
	}()
	wg.Wait()

	// Wait for v1 to start the server
	time.Sleep(3 * time.Second)

	// Connect a test client to the validator
	u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("Error connecting to the server: %v", err)
	}
	defer c.Close()

	// Send a message from the test client to the validator
	err = c.WriteJSON(map[string]interface{}{
		"validatorId": "test-client",
		"message":     "Test message from test-client",
	})
	if err != nil {
		t.Fatalf("Error sending message: %v", err)
	}

	// Wait for the message to be processed
	time.Sleep(1 * time.Second)

	// Stop the validator
	v1.Stop()
}