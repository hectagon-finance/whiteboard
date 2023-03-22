package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
)

/*** Code for lib ***/
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 52 letter, 52 = 110100b => keep 6 smallest bit, zer0 the rest
	letterIdxMask = 1<<letterIdxBits - 1 // 1b move right 6 pos -> 1000000, -1 and turn into 0111111
)
const (
	letterIdxMax = 63 / letterIdxBits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 { // remain = 0 -> init new random number and assign to cache and restart remain
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		// shift cache 6 bit to the right
		cache = cache >> letterIdxBits
		remain--
	}
	return string(b)
}

/*** End code for lib ***/

/** Code for Blockchain Logic **/
// NOTE: do not use json in real life example
type Command string
type Status string

const (
	Create Command = "Create"
	Start  Command = "Start"
	Stop   Command = "Stop"
	Pause  Command = "Pause"
	Finish Command = "Finish"
)

const (
	JustCreated Status = "JustCreated"
	Doing       Status = "Doing"
	Paused      Status = "Paused"
	Stopped     Status = "Stopped"
	Finished    Status = "Finished"
)

type CreateInstruction struct {
	Title string
	Desc  string
}
type StartInstruction struct {
	Id             string
	EstDayToFinish int
}
type StopInstruction struct {
	Id     string
	Reason string
}
type PauseInstruction struct {
	Id         string
	EstWaitDay int
}
type FinishInstruction struct {
	Id             string
	CongratMessage string
}

type Instruction struct {
	Id   string
	C    Command
	Data []byte
}

type Block struct {
	BlockHash    string
	Instructions []Instruction
}

type Task struct {
	Id     string
	Title  string
	Desc   string
	Status Status
}

type Event struct {
	Id            string
	BlockHash     string
	InstructionId string
	Msg           string
}

var events = make([]Event, 0)

func logic(mem []byte, block Block) []byte {
	var tasks []Task
	newMem := mem
	err := json.Unmarshal(mem, &tasks)
	if err != nil {
		return mem
	}
	instructions := block.Instructions
	blockHash := block.BlockHash
	for _, ins := range instructions {
		switch ins.C {
		case Create:
			fmt.Println("\ncreate")
			var createInstruction *CreateInstruction
			err = json.Unmarshal(ins.Data, &createInstruction)
			if err == nil {
				tasks = append(tasks, Task{
					Id:     RandString(8),
					Title:  createInstruction.Title,
					Desc:   createInstruction.Desc,
					Status: JustCreated,
				})
				emitEvent(blockHash, ins.Id, fmt.Sprintf("Create Task{%s, %s}", createInstruction.Title, createInstruction.Desc))
				newMem, _ = json.Marshal(tasks)
			}
		case Start:
			fmt.Println("\nstart")
			var startInstrucion *StartInstruction
			err = json.Unmarshal(ins.Data, &startInstrucion)
			if err == nil {
				t := findTask(tasks, startInstrucion.Id)
				fmt.Println(t)
				if t != nil && (t.Status == JustCreated || t.Status == Paused) {
					t.Status = Doing
					newMem, _ = json.Marshal(tasks)
					emitEvent(blockHash, ins.Id, fmt.Sprintf("Start Task #%s(%s), est to finish in %d", t.Id, t.Title, startInstrucion.EstDayToFinish))
				}
				fmt.Println(string(newMem))
			}
		case Stop:
			fmt.Println("\nstop")
			var stopInstrucion *StopInstruction
			err = json.Unmarshal(ins.Data, &stopInstrucion)
			if err == nil {
				t := findTask(tasks, stopInstrucion.Id)
				fmt.Println(t)
				if t != nil && t.Status != Finished {
					t.Status = Stopped
					emitEvent(blockHash, ins.Id, fmt.Sprintf("Stop Task #%s(%s), because of %s", t.Id, t.Title, stopInstrucion.Reason))
					newMem, _ = json.Marshal(tasks)
				}
				fmt.Println(string(newMem))
			}
		case Pause:
			fmt.Println("\npause")
			var pauseInstrucion *PauseInstruction
			err = json.Unmarshal(ins.Data, &pauseInstrucion)
			if err == nil {
				t := findTask(tasks, pauseInstrucion.Id)
				fmt.Println(t)
				if t != nil && (t.Status == JustCreated || t.Status == Doing) {
					t.Status = Paused
					emitEvent(blockHash, ins.Id, fmt.Sprintf("Pause Task #%s(%s), est to wait %d day(s)", t.Id, t.Title, pauseInstrucion.EstWaitDay))
					newMem, _ = json.Marshal(tasks)
				}
			}
		case Finish:
			fmt.Println("\nfinish")
			var finishInstrucion *FinishInstruction
			err = json.Unmarshal(ins.Data, &finishInstrucion)
			if err == nil {
				t := findTask(tasks, finishInstrucion.Id)
				fmt.Println(t)
				if t != nil && t.Status != Stopped {
					t.Status = Finished
					emitEvent(blockHash, ins.Id, fmt.Sprintf("Finish Task #%s(%s). %s", t.Id, t.Title, finishInstrucion.CongratMessage))
					newMem, _ = json.Marshal(tasks)
				}
			}
		}
	}
	return newMem
}

func emitEvent(blockHash string, instructionId string, message string) {
	id := RandString(8)
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

func main() {
	// test your understanding of my code here hehe!
	fileContent, err := ioutil.ReadFile("block_data.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	var block Block
	err = json.Unmarshal(fileContent, &block)
	if err != nil {
		fmt.Println("Error unmarshaling file content:", err)
		os.Exit(1)
	}

	// printBlock(&block)

	// Đọc dữ liệu từ file mem.txt
	oldMem, err := ioutil.ReadFile("mem.txt")
	if err != nil {
		fmt.Println("Error reading mem.txt:", err)
		os.Exit(1)
	}
	// ...
	newMem := logic(oldMem, block)

	// Lưu newMem vào file mem.txt
	err = ioutil.WriteFile("mem.txt", newMem, 0644)
	if err != nil {
		fmt.Println("Error writing to mem.txt:", err)
		os.Exit(1)
	}

	fmt.Println("\n", string(newMem))
}
