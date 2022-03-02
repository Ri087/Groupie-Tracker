package main

import (
	"GroupieTracker/GroupieTracker"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

type Main_struct struct {
	A    *GroupieTracker.Api
	F    *GroupieTracker.Filtre_Artist
	Bool bool
}

func FiltreInit(F *GroupieTracker.Filtre_Artist) {
	F.Modif = false
	F.DateStart = 0
	F.DateEnd = 0
}

func ApiInit() *GroupieTracker.Api {
	Apis := &GroupieTracker.Api{}
	GroupieTracker.ApiArtists(Apis)
	GroupieTracker.ApiDates(Apis)
	GroupieTracker.ApiLocations(Apis)
	GroupieTracker.ApiRelations(Apis)
	return Apis
}

func main() {

	Fil := &GroupieTracker.Filtre_Artist{}
	Acc := &GroupieTracker.Account{}
	CheckCreation := &GroupieTracker.CheckCreation{}
	CheckConnection := &GroupieTracker.CheckCo{}
	Apis := ApiInit()
	Main := Main_struct{A: Apis, F: Fil, Bool: false}
	FiltreInit(Fil)

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
	})
	http.HandleFunc("/filtre", func(w http.ResponseWriter, r *http.Request) {
		FuncFiltre(w, r, Fil)
	})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if !Searchbool(Apis, r.FormValue("search")) {
			http.Redirect(w, r, "#second-page", http.StatusFound)
		}
		id := NametoId(Apis, r.FormValue("search"))
		http.Redirect(w, r, "/artiste/"+id, http.StatusFound)
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
		fmt.Println(Apis.Id)
		templateshtml.ExecuteTemplate(w, "pages-artistes.html", Main)
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

func NametoId(api *GroupieTracker.Api, name string) string {
	var id_of_artist string
	for _, i := range api.ApiArtist {
		if i.Name == name {
			id_of_artist = strconv.Itoa(i.Id)
			return id_of_artist
		}
	}
	return ""
}

func Searchbool(api *GroupieTracker.Api, name string) bool {
	for _, i := range api.ApiArtist {
		if i.Name == name {
			return true
		}
	}
	return false
}

func SetCookie(w http.ResponseWriter, mail string, Acc *GroupieTracker.Account) {
	http.SetCookie(w, &http.Cookie{Name: "Token", Value: GroupieTracker.IDMail(mail)})
	Logout(Acc)
}

func Logout(Acc *GroupieTracker.Account) {
	Acc.Mail, Acc.Password, Acc.Name = "", "", ""
}

func FuncFiltre(w http.ResponseWriter, r *http.Request, Fa *GroupieTracker.Filtre_Artist) {
	date_filtre := r.FormValue("filtre_date")
	fmt.Println(date_filtre)
	if len(date_filtre) < 1 {
		Fa.Modif = false
		Fa.DateStart = 0
		Fa.DateEnd = 2100
	} else {
		Fa.Modif = true
		if date_filtre == "all" {
			Fa.DateStart = 0
			Fa.DateEnd = 2100
		} else {
			date, _ := strconv.Atoi(date_filtre)
			Fa.DateStart = date
			Fa.DateEnd = date + 9
			fmt.Println(Fa.DateStart, "|", Fa.DateEnd)
		}

	}

	http.Redirect(w, r, "/artiste", http.StatusFound)

}
