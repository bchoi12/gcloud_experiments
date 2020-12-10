package main

import (
	"net/http"
)

func bbbHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/bbb/"):]
	http.ServeFile(w, r, "bbb/" + file)
}