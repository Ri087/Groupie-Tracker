package main

import (
	"GroupieTracker/GroupieTracker"
	"encoding/base32"
	"net/http"
	"strconv"
	"text/template"
)

type Main_struct struct {
	A    *GroupieTracker.Api
	ADF  *GroupieTracker.Filter
	Bool bool
}

func ApiInit() *GroupieTracker.Api {
	Apis := &GroupieTracker.Api{}
	GroupieTracker.ApiInit(Apis)
	Apis.ApiFiltre = Apis.ApiArtist
	return Apis
}

func main() {
	Apis := ApiInit()
	Acc := &GroupieTracker.Account{}
	CheckCreation := &GroupieTracker.CheckCreation{}
	CheckConnection := &GroupieTracker.CheckCo{}
	Main := Main_struct{A: Apis, ADF: &GroupieTracker.Filter{}, Bool: false}
	GroupieTracker.FilterReset(Main.ADF)
	GroupieTracker.CountryTab(Apis, Main.ADF)
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/ressources/", http.StripPrefix("/ressources/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "index.html", Main)

	})

	//Page principal
	http.HandleFunc("/artiste", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "artiste.html", Main)
		GroupieTracker.FilterReset(Main.ADF)
	})
	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		GroupieTracker.FLT(r.URL.Query(), Apis, Main.ADF)
		http.Redirect(w, r, "/artiste", http.StatusFound)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if len(r.FormValue("search-artist")) > 0 {
			GroupieTracker.SearchNameArtsit(w, r, *Apis)
		}
		if len(r.FormValue("search-crea-artist")) > 0 {
			GroupieTracker.SearchDateArtsit(w, r, *Apis)
		}
		if len(r.FormValue("search-crea-album")) > 0 {
			GroupieTracker.SearchDateAlbum(w, r, *Apis)
		}
		if len(r.FormValue("search-membre")) > 0 {
			GroupieTracker.SearchMemberArtsit(w, r, *Apis)
		}
		if len(r.FormValue("search-location")) > 0 {
			GroupieTracker.SearchMemberArtsit(w, r, *Apis)
		}

	})
	http.HandleFunc("/filtre-search-bar", func(w http.ResponseWriter, r *http.Request) {
		GroupieTracker.FiltreSearchBar(w, r, Main.ADF)
	})

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "event.html", Main)
	})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "contact.html", Main)
	})

	http.HandleFunc("/artiste/", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		Id_Api_page, _ := strconv.Atoi(r.URL.Path[9:])
		Apis.Id = Id_Api_page - 1
		templateshtml.ExecuteTemplate(w, "pages-artistes.html", Main)
	})
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

	// NE PAS SUPPR CEST DES TESTS
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
	})
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
	})
	http.ListenAndServe(":8080", nil)
}

func Creation(w http.ResponseWriter, r *http.Request, CC *GroupieTracker.CheckCreation, Acc *GroupieTracker.Account) {
	name, pwd, pwdc, mail := r.FormValue("name"), r.FormValue("pwd"), r.FormValue("pwdc"), r.FormValue("mail")
	if GroupieTracker.CheckGoodCreation(name, pwd, pwdc, mail, CC, Acc) {
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
}

func Logout(Acc *GroupieTracker.Account) {
	Acc.Mail, Acc.Password, Acc.Name = "", []byte{}, ""
}
