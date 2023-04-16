package validator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hectagon-finance/whiteboard/types"
	"github.com/hectagon-finance/whiteboard/utils"
)

type Command string
type Status string

const (
	Create Command = "Create"
	Start  Command = "Start"
	Stop   Command = "Stop"
	Pause  Command = "Pause"
	Finish Command = "Finish"
	Assign Command = "Assign"
)

const (
	JustCreated Status = "JustCreated"
	Doing       Status = "Doing"
	Paused      Status = "Paused"
	Stopped     Status = "Stopped"
	Finished    Status = "Finished"
)

type Task struct {
	Id      string
	Title   string
	Desc    string
	Status  Status
	Owner   string
	Handler []string
}

type Event struct {
	Id            string
	BlockHash     string
	InstructionId string
	Msg           string
}

type CreateInstruction struct {
	Id    string
	Title string
	Desc  string
	From  string
}
type StartInstruction struct {
	Id             string
	EstDayToFinish int
	From           string
}
type StopInstruction struct {
	Id     string
	Reason string
	From   string
}
type PauseInstruction struct {
	Id         string
	EstWaitDay int
	From       string
}
type FinishInstruction struct {
	Id             string
	CongratMessage string
	From           string
}
type AssignInstruction struct {
	Id       string
	AssignTo string
	From     string
}

var events = make([]Event, 0)

var Chan_Block = make(chan types.Block)
var mem []byte

func Logic() {
	haha := []Task{}
	mem, _ = json.Marshal(haha)

	for {
		block := <-Chan_Block
		// fmt.Println("block:", block)
		mem = logic(block)
		log.Println("mem:", string(mem))
	}
}

type Instruction struct {
	C    Command
	Data []byte
}

func logic(block types.Block) []byte {
	var tasks []Task
	newMem := mem
	err := json.Unmarshal(mem, &tasks)
	if err != nil {
		fmt.Println("err:", err)
		return mem
	}
	transactions := block.GetTransactions()

	blockHash := utils.Byte32toStr(block.GetHash())

	for _, trans := range transactions {
		data := trans.Data
		var ins Instruction
		err := json.Unmarshal(data, &ins)
		if err != nil {
			log.Println(err)
		}
		// log.Println("ins:", ins)
		switch ins.C {
		case Create:
			fmt.Println("\ncreate")
			var createInstruction *CreateInstruction
			err = json.Unmarshal(ins.Data, &createInstruction)
			if err == nil {
				tasks = append(tasks, Task{
					Id:     createInstruction.Id,
					Title:  createInstruction.Title,
					Desc:   createInstruction.Desc,
					Status: JustCreated,
					Owner:  createInstruction.From,
				})
				emitEvent(blockHash, trans.TransactionId, fmt.Sprintf("Create Task{%s, %s}", createInstruction.Title, createInstruction.Desc))
				newMem, _ = json.Marshal(tasks)
			}
		case Assign:
			fmt.Println("\nassign")
			var assignInstruction *AssignInstruction
			err = json.Unmarshal(ins.Data, &assignInstruction)
			if err == nil {
				fmt.Println("Running")
				t := findTask(tasks, assignInstruction.Id)
				fmt.Println(t)
				fmt.Println("owner", t.Owner)
				fmt.Println("from", assignInstruction.From)
				fmt.Println("Running2")
				if t != nil && t.Owner == assignInstruction.From {
					fmt.Println("Running3")
					t.Handler = append(t.Handler, assignInstruction.AssignTo)
					emitEvent(blockHash, trans.TransactionId, fmt.Sprintf("Assign Task %s For{%s}", assignInstruction.Id, assignInstruction.AssignTo))
					newMem, _ = json.Marshal(tasks)
				}
			}

		case Start:
			fmt.Println("\nstart")
			var startInstrucion *StartInstruction
			err = json.Unmarshal(ins.Data, &startInstrucion)
			if err == nil {
				t := findTask(tasks, startInstrucion.Id)
				fmt.Println(t)
				if t != nil && (t.Status == JustCreated || t.Status == Paused) {
					if checkHandler(startInstrucion.From, t.Handler) || startInstrucion.From == t.Owner {
						t.Status = Doing
						newMem, _ = json.Marshal(tasks)
						emitEvent(blockHash, trans.TransactionId, fmt.Sprintf("Start Task #%s(%s), est to finish in %d", t.Id, t.Title, startInstrucion.EstDayToFinish))
					}
				}
			}
		case Stop:
			fmt.Println("\nstop")
			var stopInstrucion *StopInstruction
			err = json.Unmarshal(ins.Data, &stopInstrucion)
			if err == nil {
				t := findTask(tasks, stopInstrucion.Id)
				fmt.Println(t)
				if t != nil && t.Status != Finished && stopInstrucion.From == t.Owner {
					t.Status = Stopped
					emitEvent(blockHash, trans.TransactionId, fmt.Sprintf("Stop Task #%s(%s), because of %s", t.Id, t.Title, stopInstrucion.Reason))
					newMem, _ = json.Marshal(tasks)
				}
			}
		case Pause:
			fmt.Println("\npause")
			var pauseInstrucion *PauseInstruction
			err = json.Unmarshal(ins.Data, &pauseInstrucion)
			if err == nil {
				t := findTask(tasks, pauseInstrucion.Id)
				fmt.Println(t)
				if t != nil && (t.Status == JustCreated || t.Status == Doing) {
					if checkHandler(pauseInstrucion.From, t.Handler) || pauseInstrucion.From == t.Owner {
						t.Status = Paused
						emitEvent(blockHash, trans.TransactionId, fmt.Sprintf("Pause Task #%s(%s), est to wait %d day(s)", t.Id, t.Title, pauseInstrucion.EstWaitDay))
						newMem, _ = json.Marshal(tasks)
					}
				}
			}
		case Finish:
			fmt.Println("\nfinish")
			var finishInstrucion *FinishInstruction
			err = json.Unmarshal(ins.Data, &finishInstrucion)
			if err == nil {
				t := findTask(tasks, finishInstrucion.Id)
				fmt.Println(t)
				if t != nil && t.Status != Stopped && finishInstrucion.From == t.Owner {
					t.Status = Finished
					emitEvent(blockHash, trans.TransactionId, fmt.Sprintf("Finish Task #%s(%s). %s", t.Id, t.Title, finishInstrucion.CongratMessage))
					newMem, _ = json.Marshal(tasks)
				}
			}
		}
	}
	return newMem

}

func checkHandler(from string, handlers []string) bool {
	for _, handler := range handlers {
		if handler == from {
			return true
		}
	}
	return false
}

func emitEvent(blockHash string, instructionId string, message string) {
	id := utils.RandString(8)
	ev := Event{
		Id:            id,
		BlockHash:     blockHash,
		InstructionId: instructionId,
		Msg:           message,
	}
	fmt.Println("Emit: ", ev) // event will be sent to clients
	events = append(events, ev)
}

func findTask(tasks []Task, Id string) *Task {
	for i := range tasks {
		if tasks[i].Id == Id {
			return &tasks[i]
		}
	}
	return nil
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
