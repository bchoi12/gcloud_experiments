package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"net/http"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi!")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/view/"):]
		p, err := loadPage(title)
		if err != nil {
				http.Redirect(w, r, "/edit/"+title, http.StatusFound)
				return
		}
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func stickerHandler(w http.ResponseWriter, r *http.Request) {
	img := r.URL.Path[len("/sticker/"):]
	http.ServeFile(w, r, "sticker/" + img)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/edit/"):]
		p, err := loadPage(title)
		if err != nil {
				p = &Page{Title: title}
		}
		t, _ := template.ParseFiles("edit.html")
		t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/save/"):]
		body := r.FormValue("body")
		p := &Page{Title: title, Body: []byte(body)}
		p.save()
		http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", viewHandler)
	// http.HandleFunc("/edit/", editHandler)
	// http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/sticker/", stickerHandler)

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