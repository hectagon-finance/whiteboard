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

	// Wait for a short time to allow the server to start
	time.Sleep(2 * time.Second)

	// Stop the validator
	v.Stop()

	// If no errors or panics occurred during the test, the test passed
}

func TestValidatorSendMessage(t *testing.T) {
	// Initialize a new validator
	v1 := NewValidator(8081)

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
	time.Sleep(2 * time.Second)

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