package models

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

type Room struct {
	Id      string
	Clients map[*websocket.Conn]bool
}

func NewRoom(id string, firstClient *websocket.Conn) Room {
	return Room{Id: id, Clients: map[*websocket.Conn]bool{firstClient: true}}
}

func CloseRoom(roomId string, rdb *redis.Client, rooms map[string]Room) {
	delete(rooms, roomId)

	if err := rdb.Del("room:" + roomId + ":messages").Err(); err != nil {
		panic(err)
	}
}

func DisconnectClient(room Room, client *websocket.Conn) {
	client.Close()
	delete(room.Clients, client)
}
