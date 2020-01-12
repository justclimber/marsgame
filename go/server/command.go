package server

import "fmt"

type Command struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (c Command) ToSting() string {
	return fmt.Sprintf("type: %s, payload: %s", c.Type, c.Payload)
}
