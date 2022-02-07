package main

import (
	"GroupieTracker/GroupieTracker"
	"fmt"
	"net/http"
	"text/template"
)

func main() {
	GroupieTracker.TraitementArtiste()
	GroupieTracker.TraitementDates()
	GroupieTracker.TraitementLocations()
	GroupieTracker.TraitementRealation()
	fmt.Println(GroupieTracker.All)
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/ressources/", http.StripPrefix("/ressources/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		err := templateshtml.ExecuteTemplate(w, "index.html", GroupieTracker.All.Artists)
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
