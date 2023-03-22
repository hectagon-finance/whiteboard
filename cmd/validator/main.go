package main

// import (
// 	"os"
// 	"strconv"
// 	"time"

// 	."github.com/hectagon-finance/whiteboard/types"
// )

// func main() {
// 	if len(os.Args) >= 2 {
// 		current_validator_id := os.Args[1]

// 		a, err := strconv.Atoi(current_validator_id)
// 		if err != nil {
// 			// ... handle error
// 			panic(err)
// 		}

// 		// Create two validators
// 		v := types.NewValidator(a)

// 		// Start the validators
// 		v.Start()

// 		// Wait for a few seconds to let the validators establish connections
// 		time.Sleep(100 * time.Second)

// 		// Stop the validators
// 		v.Stop()
// 	}
// }
