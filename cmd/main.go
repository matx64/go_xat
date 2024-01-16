package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/matx64/go_xat/controllers"
	"github.com/matx64/go_xat/db"
	"github.com/matx64/go_xat/models"
)

var (
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

	rooms := make(map[string]models.Room)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		handleConnection(w, r, rdb, rooms)
	})

	go handleMessages(rdb, rooms)

	addr := os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")
	fmt.Print("🚀 Server started at " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleConnection(w http.ResponseWriter, r *http.Request, rdb *redis.Client, rooms map[string]models.Room) {
	client, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Upgrade error: %v", err)
	}
	defer client.Close()

	query := r.URL.Query()
	username := query.Get("username")
	roomId := query.Get("roomId")

	insertClientInRoom(roomId, username, client, rdb, rooms)

	listenToClientMessages(rooms[roomId], client)

	if len(rooms[roomId].Clients) == 0 {
		models.CloseRoom(roomId, rdb, rooms)
		return
	}

	broadcaster <- models.NewMessage(roomId, username, "left", "")
}

func insertClientInRoom(roomId, username string, client *websocket.Conn, rdb *redis.Client, rooms map[string]models.Room) {
	if room, ok := rooms[roomId]; ok {
		// room already open
		room.Clients[client] = true

		if rdb.Exists("room:"+roomId+":messages").Val() != 0 {
			controllers.SendPreviousMessages(room, client, rdb)
		}
	} else {
		rooms[roomId] = models.NewRoom(roomId, client)
	}

	broadcaster <- models.NewMessage(roomId, username, "join", "")
}

func listenToClientMessages(room models.Room, client *websocket.Conn) {
	for {
		msg := models.Message{}

		err := client.ReadJSON(&msg)
		if err != nil {
			log.Printf("Reading error: %#v\n", err)
			delete(room.Clients, client)
			break
		}

		log.Printf("recv: message %q", msg)

		broadcaster <- msg
	}
}

func handleMessages(rdb *redis.Client, rooms map[string]models.Room) {
	for {
		msg := <-broadcaster
		db.StoreMessage(msg, rdb)
		controllers.MessageClients(msg, rooms[msg.RoomId])
	}
}
