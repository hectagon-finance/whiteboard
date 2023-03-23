package validator

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	. "github.com/hectagon-finance/whiteboard/types"
)

type Validator struct {
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

func (v *Validator) Id() string {
	return v.ValidatorId
}

func (v *Validator) GetPublicKey() string {
	return v.PublicKey
}

func (v *Validator) GetPrivateKey() string {
	return v.PrivateKey
}

func (v *Validator) GetMemPool() MemPool {
	return v.MemPool
}

func (v *Validator) GetBalance() int64 {
	return v.Balance
}

func (v *Validator) GetStake() int64 {
	return v.Stake
}

func (v *Validator) GetStatus() string {
	return v.Status
}

func (v *Validator) GetLastBlock() Block {
	return v.LastBlock
}

func (v *Validator) GetPort() int {
	return v.Port
}
