package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Message struct {
	RoomId   string `json:"roomId"`
	Username string `json:"username"`
	Text     string `json:"text"`
	Type     string `json:"type"`
	SentAt   int64  `json:"sentAt"`
}

// var (
// 	db *redis.Client
// )

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[string]map[*websocket.Conn]bool)
var broadcaster = make(chan Message)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// db = redis.NewClient(&redis.Options{
	// 	Addr:     os.Getenv("REDIS_ADDR"),
	// 	Password: os.Getenv("REDIS_PASSWORD"),
	// 	DB:       0,
	// })

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/connect", handleConnection)

	go handleMessages()

	fmt.Println("Server started.")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil))
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
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
		val[ws] = true
	} else {
		clients[roomId] = map[*websocket.Conn]bool{ws: true}
	}

	broadcaster <- Message{RoomId: roomId, Username: username, Text: fmt.Sprintf("%s joined the chat.", username), Type: "join", SentAt: time.Now().Unix()}

	for {
		var msg Message

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
		return
	}

	broadcaster <- Message{RoomId: roomId, Username: username, Text: fmt.Sprintf("%s left the chat.", username), Type: "left", SentAt: time.Now().Unix()}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func handleMessages() {
	for {
		msg := <-broadcaster
		roomId := msg.RoomId

		for client := range clients[roomId] {
			err := client.WriteJSON(msg)
			if err != nil && unsafeError(err) {
				log.Printf("Write error: %v", err)
				delete(clients[roomId], client)
				client.Close()
			}
		}
	}
}
