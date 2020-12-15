package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi!")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/sticker/", stickerHandler)
	http.HandleFunc("/bbb/", bbbHandler)
	http.HandleFunc("/chatclient", chatClientHandler)
	http.HandleFunc("/chat/", chatHandler)

	go handleMessages()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal(err)
	}
}