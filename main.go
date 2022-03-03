package main

import (
	"GroupieTracker/GroupieTracker"
	"encoding/base32"
	"net/http"
	"strconv"
	"text/template"
)

func ApiInit() *GroupieTracker.Api {
	Apis := &GroupieTracker.Api{}
	GroupieTracker.ApiArtists(Apis)
	GroupieTracker.ApiDates(Apis)
	GroupieTracker.ApiLocations(Apis)
	GroupieTracker.ApiRelations(Apis)
	return Apis
}

func main() {
	Acc := &GroupieTracker.Account{}
	CheckCreation := &GroupieTracker.CheckCreation{}
	CheckConnection := &GroupieTracker.CheckCo{}

	Apis := ApiInit()

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/ressources/", http.StripPrefix("/ressources/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "index.html", Apis)

	})

	//Page principal
	http.HandleFunc("/artiste", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "artiste.html", Apis)
	})
	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "event.html", "")
	})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "contact.html", "")
	})

	http.HandleFunc("/artiste/", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		Id_Api_page, _ := strconv.Atoi(r.URL.Path[9:])
		Apis.Id = Id_Api_page - 1
		templateshtml.ExecuteTemplate(w, "pages-artistes.html", Apis)
	})

	// Profil pages
	http.HandleFunc("/connection", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "connection.html", CheckConnection)
		CheckConnection.Mail, CheckConnection.Pwd = false, false
	})
	http.HandleFunc("/creation", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "creation.html", CheckCreation)
		CheckCreation.Name, CheckCreation.Pwd, CheckCreation.Pwdc, CheckCreation.Mail, CheckCreation.Exist = false, false, false, false, false
	})
	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("AUTHENTIFICATION_TOKEN")
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
		if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
			cookie, _ := r.Cookie("AUTHENTIFICATION_TOKEN")
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
		SetCookie(w, mail, pwd, Acc)
		http.Redirect(w, r, "/profil", http.StatusFound)
	}
}

func Login(w http.ResponseWriter, r *http.Request, CC *GroupieTracker.CheckCo, Acc *GroupieTracker.Account) {
	mail, pwd := r.FormValue("mail"), r.FormValue("pwd")
	if GroupieTracker.CheckConnection(mail, pwd, CC, Acc) {
		SetCookie(w, mail, pwd, Acc)
		http.Redirect(w, r, "/profil", http.StatusFound)
	} else {
		http.Redirect(w, r, "/connection", http.StatusFound)
	}
}

func SetCookie(w http.ResponseWriter, mail, pwd string, Acc *GroupieTracker.Account) {
	http.SetCookie(w, &http.Cookie{Name: "AUTHENTIFICATION_TOKEN", Value: base32.StdEncoding.EncodeToString(GroupieTracker.Cryptage(mail))})
	Logout(Acc)
}

func Logout(Acc *GroupieTracker.Account) {
	Acc.Mail, Acc.Password, Acc.Name = "", []byte{}, ""
}
