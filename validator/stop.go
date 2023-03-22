package validator

import (
	"log"
)

func Stop(v *ValidatorStruct) {
	StopServer(v)
	StopClient(v)
	v.Status = "inactive"
}

func  StopClient(v *ValidatorStruct) {
	v.ClientsMutex.Lock()
	for conn := range v.Clients {
		conn.Close()
		delete(v.Clients, conn)
	}
	v.ClientsMutex.Unlock()
}

func StopServer(v *ValidatorStruct) {
	if v.StopServer != nil {
		log.Printf("Validator %s stopping server on port %d", v.ValidatorId, v.Port)
		v.StopServer()
	}
}