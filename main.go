package main

import (
	"GroupieTracker/GroupieTracker"
	"math/rand"
	"net/http"
	"text/template"
)

func main() {
	CheckCreation := &GroupieTracker.CheckCreation{}
	CheckConnection := &GroupieTracker.CheckCo{}
	Acc := &GroupieTracker.Account{}
	GroupieTracker.ApiArtists()
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
