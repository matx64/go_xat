package models

import (
	"fmt"
	"time"
)

type Message struct {
	RoomId   string `json:"roomId"`
	Username string `json:"username"`
	Text     string `json:"text"`
	Kind     string `json:"kind"`
	SentAt   int64  `json:"sentAt"`
}

func NewMessage(roomId, username, kind, text string) Message {
	if kind == "join" {
		text = fmt.Sprintf("%s joined the chat.", username)
	} else if kind == "left" {
		text = fmt.Sprintf("%s left the chat.", username)
	}

	return Message{RoomId: roomId, Username: username, Text: text, Kind: kind, SentAt: time.Now().Unix()}
}
