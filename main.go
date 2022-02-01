package main

import (
	"GroupieTracker/GroupieTracker"
	"net/http"
	"text/template"
)

func main() {
	GroupieTracker.ApiDates()
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/ressources/", http.StripPrefix("/ressources/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		err := templateshtml.ExecuteTemplate(w, "index.html", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	//Page game.html
	http.HandleFunc("/artiste", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "artiste.html", "")
	})

	http.ListenAndServe(":8080", nil)
}
