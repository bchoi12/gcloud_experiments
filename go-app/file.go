package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func serveFiles(dir string) {
	http.HandleFunc(dir, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprintf(w, "Hi!")
			return
		}
		url := r.URL.Path[1:]

		if file, err := os.Stat(url); os.IsNotExist(err) {
			log.Printf("404 for url %s", url)
			notFound(w)
		} else if file.IsDir() {
			url += "/index.html"
			if _, err := os.Stat(url); err == nil {
				http.FileServer(http.Dir(r.URL.Path))
			} else {
				log.Printf("404 for dir %s", url)
				notFound(w)
			}
		} else {
			http.ServeFile(w, r, url)
		}
	})
}

func notFound(w http.ResponseWriter) {
	fmt.Fprintf(w, "404!")
}