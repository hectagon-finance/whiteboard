package main

import (
	"github.com/onspeedhp/project/types"
	"time"
)

func main() {
	// Initialize a new validator
	v := types.NewValidator(8080)

	// Start the validator and wait for a short duration
	v.Start()

	// // Wait for a short time to allow the server to start
	time.Sleep(10 * time.Second)

	// Stop the validator
	v.Stop()

	// If no errors or panics occurred during the test, the test passed
}