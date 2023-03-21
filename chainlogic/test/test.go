package main

import (
	"encoding/json"
	"fmt"
)

type Command string

type Instruction struct {
	Id   string
	C    Command
	Data []byte
}

type Block struct {
	BlockHash    string
	Instructions []Instruction
}

func (i *Instruction) UnmarshalJSON(data []byte) error {
	var raw map[string]*json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*raw["Id"], &i.Id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*raw["C"], &i.C)
	if err != nil {
		return err
	}

	var rawData map[string]interface{}
	err = json.Unmarshal(*raw["Data"], &rawData)
	if err != nil {
		return err
	}

	i.Data, err = json.Marshal(rawData)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	data := `{
      "BlockHash":"dsahdasdhkadkjasdjsjkdhdkaskdjasjkd",
      "Instructions":[
         {
            "Id":"2",
            "C":"Start",
            "Data":{
               "Id":"TEpWCQms",
               "EstDayToFinish": 2
            }
         }
      ]
   }`

	var block Block
	err := json.Unmarshal([]byte(data), &block)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	encodedJSON, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling modified JSON:", err)
		return
	}

	fmt.Println("Modified JSON with base64 encoded Data fields:")
	fmt.Println(string(encodedJSON))
}

