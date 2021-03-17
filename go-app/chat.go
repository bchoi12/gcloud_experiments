package main

import (
	"net/http"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

var upgrader = websocket.Upgrader{}

var clients = make(map[string]map[*websocket.Conn]bool)
var broadcast = make(map[string]chan ChatMessage)
var msgQueue = make(map[string][]ChatMessage)


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

func sendMessage(room string, client *websocket.Conn, msg ChatMessage) {
	cd := ChatData{}
	cd.ChatMessage = &msg

	data, err := proto.Marshal(&cd)
	if err != nil {
		log.Printf("error packaging message : %v", err)
	}

	err = client.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Printf("error writing out message: %v", err)
	}
}

func broadcastMessage(room string, msg ChatMessage) {
	msgQueue[room] = append(msgQueue[room], msg)
	if (len(msgQueue[room]) > 50) {
		msgQueue[room] = msgQueue[room][1:51]
	}
	broadcast[room] <- msg
}

func chatClientHandler(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Path[len("/chatclient/"):]

	log.Printf("room opened/joined: %s", room)

	if len(room) == 0 {
		log.Printf("error: room name is empty")
		return
	}

	if clients[room] == nil {
		clients[room] = make(map[*websocket.Conn]bool)
		broadcast[room] = make(chan ChatMessage)
		if msgQueue[room] == nil {
			msgQueue[room] = make([]ChatMessage, 0)
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
		_, msg, err := ws.ReadMessage()

		if err != nil {
			log.Printf("error reading incoming message: %v", err)

			if e, ok := err.(*websocket.CloseError); ok {
				switch e.Code {
					case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					return
				}
			}
		}

		chatData := ChatData{}
		err = proto.Unmarshal(msg, &chatData)

		if err != nil {
			log.Printf("error parsing incoming message: %v", err)
			continue
		}

		if chatData.GetChatMessage() != nil {
			broadcastMessage(room, *chatData.GetChatMessage())
		}
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	f := r.URL.Path[1:]

	if info, err := os.Stat(f); err == nil && !info.IsDir() {
  		http.ServeFile(w, r, f)
	} else {
		http.ServeFile(w, r, "chat/index.html")
	}
}

func closeConnection(room string, ws *websocket.Conn) {
	ws.Close()
	delete(clients[room], ws)
}