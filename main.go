package main

import (
	"GroupieTracker/GroupieTracker"
	"math/rand"
	"net/http"
	"text/template"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/ressources/", http.StripPrefix("/ressources/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Artist := GroupieTracker.ApiArtists()
		N := rand.Intn(len(Artist) - 3)
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		err := templateshtml.ExecuteTemplate(w, "index.html", Artist[N:N+3])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	//Page principal
	http.HandleFunc("/artiste", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "artiste.html", "")
	})
	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "event.html", "")
	})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "contact.html", "")
	})
	http.HandleFunc("/connection", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "connection.html", "")
	})

	http.ListenAndServe(":8080", nil)
}
