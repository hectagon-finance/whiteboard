package validator

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hectagon-finance/whiteboard/types"
)

var Chan_Block = make(chan types.Block)
var mem []string

func Logic() {
	for {
		select {
		case block := <-Chan_Block:
			for _, txn := range block.GetTransactions() {

				mem = append(mem, string(txn.Data))
			}
		}
	}
}

func ClientReadHandler() {
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		// write json to response\
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mem)
	})

	log.Fatal(http.ListenAndServe("localhost:"+"1"+Port, nil))
}
