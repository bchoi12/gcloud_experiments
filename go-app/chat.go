package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

var upgrader = websocket.Upgrader{}

var clients = make(map[string]map[*websocket.Conn]bool)
var broadcast = make(map[string]chan Message)

type Message struct {
	Username string `json:"username"`
	Message string `json:"message"`
}
var msgQueue = make(map[string][]Message)


func handleMessages(room string) {
	for {
		msg := <-broadcast[room]

		if (len(clients[room]) == 0) {
			delete(clients, room)
			delete(broadcast, room)
			break
		}

		for client := range clients[room] {
			sendMessage(room, client, msg)
		}
	}
}

func sendMessage(room string, client *websocket.Conn, msg Message) {
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error writing out message: %v", err)
		client.Close()
		delete(clients[room], client)
	}
}

func broadcastMessage(room string, msg Message) {
	msgQueue[room] = append(msgQueue[room], msg)
	if (len(msgQueue[room]) > 50) {
		msgQueue[room] = msgQueue[room][1:50]
	}
	broadcast[room] <- msg
}

func chatClientHandler(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Path[len("/chatclient/"):]

	log.Printf("room is %s", room)

	if len(room) == 0 {
		log.Printf("error: room name is empty")
		return
	}

	if clients[room] == nil {
		clients[room] = make(map[*websocket.Conn]bool)
		broadcast[room] = make(chan Message)
		if msgQueue[room] == nil {
			msgQueue[room] = make([]Message, 0)
		}
	}

	if len(clients[room]) == 0 {
		go handleMessages(room)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range msgQueue[room] {
		sendMessage(room, ws, msg)
	}

	defer closeConnection(room, ws)

	clients[room][ws] = true

	for {
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error reading incoming message: %v", err)
			delete(clients[room], ws)
			break
		}

		msgQueue[room] = append(msgQueue[room], msg)
		if (len(msgQueue[room]) > 50) {
			msgQueue[room] = msgQueue[room][1:50]
		}
		broadcast[room] <- msg
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	ext := r.URL.Path[len("/chat/"):]
	file := "chat/" + ext
	_, err := os.Stat(file)
	if len(ext) == 0 || os.IsNotExist(err) {
		http.ServeFile(w, r, "chat/")
	} else {
  		http.ServeFile(w, r, file)
	}
}

func closeConnection(room string, ws *websocket.Conn) {
	ws.Close()
	delete(clients[room], ws)
}