package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/matx64/go_xat/db"
	"github.com/matx64/go_xat/models"
)

var (
	clients     = make(map[string]map[*websocket.Conn]bool)
	broadcaster = make(chan models.Message)
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rdb := db.StartRedis()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		handleConnection(w, r, rdb)
	})

	go handleMessages(rdb)

	fmt.Println("🚀 Server started.")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil))
}

func handleConnection(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		log.Fatal(err)
	}
	defer ws.Close()

	query := r.URL.Query()
	username := query.Get("username")
	roomId := query.Get("roomId")

	if val, ok := clients[roomId]; ok {
		// room already open
		val[ws] = true

		if rdb.Exists("room:"+roomId+":messages").Val() != 0 {
			sendPreviousMessages(roomId, ws, rdb)
		}

	} else {
		clients[roomId] = map[*websocket.Conn]bool{ws: true}
	}

	broadcaster <- models.NewMessage(roomId, username, "join", "")

	for {
		var msg models.Message

		err = ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Reading error: %#v\n", err)
			delete(clients[roomId], ws)
			break
		}

		log.Printf("recv: message %q", msg)

		broadcaster <- msg
	}

	if len(clients[roomId]) == 0 {
		delete(clients, roomId)

		if err := rdb.Del("room:" + roomId + ":messages").Err(); err != nil {
			panic(err)
		}

		return
	}

	broadcaster <- models.NewMessage(roomId, username, "left", "")
}

func sendPreviousMessages(roomId string, ws *websocket.Conn, rdb *redis.Client) {
	messages, err := rdb.LRange("room:"+roomId+":messages", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	for _, message := range messages {
		var msg models.Message
		json.Unmarshal([]byte(message), &msg)
		messageClient(ws, msg)
	}
}

func messageClient(ws *websocket.Conn, msg models.Message) {
	err := ws.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("Write error: %v", err)
		ws.Close()
		delete(clients[msg.RoomId], ws)
	}
}

func messageClients(msg models.Message) {
	for client := range clients[msg.RoomId] {
		messageClient(client, msg)
	}
}

func handleMessages(rdb *redis.Client) {
	for {
		msg := <-broadcaster
		db.StoreMessage(msg, rdb)
		messageClients(msg)
	}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}
