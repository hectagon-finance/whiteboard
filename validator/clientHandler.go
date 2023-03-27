package validator

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	. "github.com/hectagon-finance/whiteboard/types"
)

type Msg struct {
	memPool MemPool
}

var Chan_1 = make(chan Msg, 100)

var DraftBlock Block

func ClientHandler(v *Validator, is_genesis string) {
	if is_genesis != "genesis" {
		Sync(v)
	}
	fmt.Println("Server is running on port", v.Addr)
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", v.Serve)
	log.Fatal(http.ListenAndServe(v.Addr, nil))
}

func broadcastBlockHash() {
	for {
		i := 0
		msg := <- Chan_1
		if (i == 5) || (i == 23) {
			log.Println(msg)
		}
	}
	// if v.MemPool.Size() == 3  || counter == 10 {
	// 	v.TempBlock = NewBlock(1, [32]byte{}, v.MemPool.GetTransactions())
	// 	blockHash := v.TempBlock.GetHash()
	// 	BroadcastBlockHash(v, blockHash)
	// 	v.MemPool.Clear()
	// 	return 0
	// }

	// counter++
	// return counter
}
