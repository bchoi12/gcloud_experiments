package main

import (
	"net/http"
)

func stickerHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/sticker/"):]
	http.ServeFile(w, r, "sticker/" + file)
}