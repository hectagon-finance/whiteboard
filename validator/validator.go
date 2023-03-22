package validator

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
)

type ValidatorStruct struct {
	ValidatorId  string
	PublicKey    string
	PrivateKey   string
	Blockchain   Blockchain
	MemPool      MemPool
	Balance      int64
	Stake        int64
	Status       string
	LastBlock    Block
	Port         int
	HttpServer   *http.Server
	ClientsMutex sync.Mutex
	Clients      map[*websocket.Conn]bool
	Peers        []string
	StopServer   func()
}

func (v *ValidatorStruct) Id() string {
	return v.ValidatorId
}

func (v *ValidatorStruct) GetPublicKey() string {
	return v.PublicKey
}

func (v *ValidatorStruct) GetPrivateKey() string {
	return v.PrivateKey
}

func (v *ValidatorStruct) GetMemPool() MemPool {
	return v.MemPool
}

func (v *ValidatorStruct) GetBalance() int64 {
	return v.Balance
}

func (v *ValidatorStruct) GetStake() int64 {
	return v.Stake
}

func (v *ValidatorStruct) GetStatus() string {
	return v.Status
}

func (v *ValidatorStruct) GetLastBlock() Block {
	return v.LastBlock
}

func (v *ValidatorStruct) GetPort() int {
	return v.Port
}
