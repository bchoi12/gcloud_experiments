package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

type Message struct {
	Username string `json:"username"`
	Message string `json:"message"`
}
var msgQueue = make([]Message, 0)


func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			sendMessage(client, msg)
		}
	}
}

func sendMessage(client *websocket.Conn, msg Message) {
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error writing out message: %v", err)
		client.Close()
		delete(clients, client)
	}
}

func broadcastMessage(msg Message) {
	msgQueue = append(msgQueue, msg)
	if (len(msgQueue) > 50) {
		msgQueue = msgQueue[1:50]
	}
	broadcast <- msg
}

func chatClientHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range msgQueue {
		sendMessage(ws, msg)
	}

	defer closeConnection(ws)

	clients[ws] = true

	for {
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error reading incoming message: %v", err)
			delete(clients, ws)
			break
		}

		msgQueue = append(msgQueue, msg)
		if (len(msgQueue) > 50) {
			msgQueue = msgQueue[1:50]
		}
		broadcast <- msg
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/chat/"):]
	http.ServeFile(w, r, "chat/" + file)
}

func closeConnection(ws *websocket.Conn) {
	ws.Close()
	delete(clients, ws)
}

/*
func bbbHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/bbb/"):]
	http.ServeFile(w, r, "bbb/" + file)
}
*/