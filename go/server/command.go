package server

import (
	"encoding/json"
	"fmt"
	"log"
)

type Command struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (c Command) ToSting() string {
	return fmt.Sprintf("type: %s, payload: %s", c.Type, c.Payload)
}

func PackStructToCommand(name string, payload interface{}) *Command {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	return &Command{
		Type:    name,
		Payload: string(jsonBytes),
	}
}
