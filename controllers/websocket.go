package controllers

import (
	"encoding/json"
	"io"
	"log"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/matx64/go_xat/models"
)

func MessageClient(msg models.Message, room models.Room, client *websocket.Conn) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("Write error: %v", err)
		models.DisconnectClient(room, client)
	}
}

func MessageClients(msg models.Message, room models.Room) {
	for client := range room.Clients {
		MessageClient(msg, room, client)
	}
}

func SendPreviousMessages(room models.Room, client *websocket.Conn, rdb *redis.Client) {
	messages, err := rdb.LRange("room:"+room.Id+":messages", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	for _, message := range messages {
		var msg models.Message
		json.Unmarshal([]byte(message), &msg)
		MessageClient(msg, room, client)
	}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}
