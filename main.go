package main

import (
	"GroupieTracker/GroupieTracker"
	"fmt"

	// "math/rand"
	"net/http"
	"text/template"
)

type api struct {
	ApiArtist    []GroupieTracker.Artist
	ApiDates     []GroupieTracker.Dates
	ApiLocations []GroupieTracker.Locations
	ApiRelations []GroupieTracker.Relations
}

func main() {
	CheckCreation := &GroupieTracker.CheckCreation{}
	CheckConnection := &GroupieTracker.CheckCo{}
	Acc := &GroupieTracker.Account{}
	// Art := GroupieTracker.Artist{}
	// GroupieTracker.ApiArtists()
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/ressources/", http.StripPrefix("/ressources/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Apis := api{}
		Apis.ApiArtist = GroupieTracker.ApiArtists()
		Date := GroupieTracker.ApiDates()
		Apis.ApiDates = append(Apis.ApiDates, Date)
		Location := GroupieTracker.ApiLocations()
		Apis.ApiLocations = append(Apis.ApiLocations, Location)
		Relation := GroupieTracker.ApiRelations()
		Apis.ApiRelations = append(Apis.ApiRelations, Relation)
		fmt.Println(Apis)
		// N := rand.Intn(len(Artist) - 3)
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "index.html", Apis)

	})
	//Page principal
	http.HandleFunc("/artiste", func(w http.ResponseWriter, r *http.Request) {
		artist := GroupieTracker.ApiArtists()
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "artiste.html", artist)
	})
	// http.HandleFunc("/actionfiltre", func(w http.ResponseWriter, r *http.Request) {
	// 	ActionFiltre(w, r, Art)
	// })
	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "event.html", "")
	})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "contact.html", "")
	})

	// Profil pages

	http.HandleFunc("/connection", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("Token"); err == nil {
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "connection.html", CheckConnection)
		CheckConnection.Mail, CheckConnection.Pwd = false, false
	})
	http.HandleFunc("/creation", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("Token"); err == nil {
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "creation.html", CheckCreation)
		CheckCreation.Name, CheckCreation.Pwd, CheckCreation.Pwdc, CheckCreation.Mail, CheckCreation.Exist = false, false, false, false, false
	})
	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Token")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		GroupieTracker.LoginAcc(cookie.Value, Acc)
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "profil.html", Acc)
		Logout(Acc)
	})

	http.HandleFunc("/checkcreation", func(w http.ResponseWriter, r *http.Request) {
		Creation(w, r, CheckCreation, Acc)
	})

	http.HandleFunc("/checkconnection", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, CheckConnection, Acc)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("Token"); err == nil {
			cookie, _ := r.Cookie("Token")
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// End of profil pages

	http.ListenAndServe(":8080", nil)
}

func Creation(w http.ResponseWriter, r *http.Request, CC *GroupieTracker.CheckCreation, Acc *GroupieTracker.Account) {
	name, pwd, pwdc, mail := r.FormValue("name"), r.FormValue("pwd"), r.FormValue("pwdc"), r.FormValue("mail")
	if GroupieTracker.CheckCrea(name, pwd, pwdc, mail, CC, Acc) {
		http.Redirect(w, r, "/creation", http.StatusFound)
	} else {
		SetCookie(w, mail, Acc)
		http.Redirect(w, r, "/profil", http.StatusFound)
	}
}

func Login(w http.ResponseWriter, r *http.Request, CC *GroupieTracker.CheckCo, Acc *GroupieTracker.Account) {
	mail, pwd := r.FormValue("mail"), r.FormValue("pwd")
	if GroupieTracker.CheckConnection(mail, pwd, CC, Acc) {
		SetCookie(w, mail, Acc)
		http.Redirect(w, r, "/profil", http.StatusFound)
	} else {
		http.Redirect(w, r, "/connection", http.StatusFound)
	}
}

func SetCookie(w http.ResponseWriter, mail string, Acc *GroupieTracker.Account) {
	http.SetCookie(w, &http.Cookie{Name: "Token", Value: GroupieTracker.IDMail(mail)})
	Logout(Acc)
}

func Logout(Acc *GroupieTracker.Account) {
	Acc.Mail, Acc.Password, Acc.Name = "", "", ""
}

// func ActionFiltre(w http.ResponseWriter, r *http.Request, a GroupieTracker.Artist) {
// 	filtre := r.FormValue("filtre")
// 	var NewArtistPrint []int
// 	for _, x := range filtre {
// 		for _, i := range string(a.CreationDate) {
// 			if i > x {
// 				NewArtistPrint = append(NewArtistPrint, a.Id)
// 			}
// 		}
// 	}
// 	Artist := GroupieTracker.ApiArtists()
// 	for _, l := range NewArtistPrint {
// 		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
// 		err := templateshtml.ExecuteTemplate(w, "index.html", Artist[l])
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// 	http.Redirect(w, r, "/artiste", http.StatusFound)

// }
