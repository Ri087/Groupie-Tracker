package main

import (
	"GroupieTracker/GroupieTracker"
	"net/http"
	"strconv"
	"text/template"
)

type MainStructure struct {
	ApiStruct *GroupieTracker.ApiStructure
	AccStruct *GroupieTracker.AccountStruct
}

func MainStructureInit() *MainStructure {
	Main := &MainStructure{}
	Main.ApiStruct = GroupieTracker.ApiStructInit()
	Main.AccStruct = GroupieTracker.AccStructureInit()
	return Main
}

func main() {
	Main := MainStructureInit()

	fileServer := http.FileServer(http.Dir("./static"))
	var s GroupieTracker.Spotify = GroupieTracker.New("6b053d7dfcbe4c69a576561f8c098391", "d00791e8792a4f13bc1bb8b95197505d")
	s.Authorize()
	http.Handle("/ressources/", http.StripPrefix("/ressources/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Main.ApiStruct.TabApiArtiste, Main.ApiStruct.TabApiArtisteLocations = GroupieTracker.ApiArtistsArtiste()
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "index.html", Main)
	})

	//Page principal
	http.HandleFunc("/artiste", func(w http.ResponseWriter, r *http.Request) {
		Main.ApiStruct.TabApiArtiste, Main.ApiStruct.TabApiArtisteLocations = GroupieTracker.ApiArtistsArtiste()
		GroupieTracker.TabCountry(Main.ApiStruct)
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "artiste.html", Main)
		GroupieTracker.FilterReset(Main.ApiStruct)
	})

	http.HandleFunc("/artiste2", func(w http.ResponseWriter, r *http.Request) {
		Main.ApiStruct.TabApiArtiste, Main.ApiStruct.TabApiArtisteLocations = GroupieTracker.ApiArtistsArtiste()
		GroupieTracker.TabCountry(Main.ApiStruct)
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "artiste2.html", Main)
		GroupieTracker.FilterReset(Main.ApiStruct)
	})
	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		GroupieTracker.FLT(r.URL.Query(), Main.ApiStruct)
		http.Redirect(w, r, "/artiste2", http.StatusFound)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		id := GroupieTracker.SearchNameArtsit(r.FormValue("search-artist"), Main.ApiStruct)
		if id == "" {
			http.Redirect(w, r, "/#second-page", http.StatusFound)
		} else {
			http.Redirect(w, r, "/artiste/"+id, http.StatusFound)
		}
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
		IDArtist := r.URL.Path[9:]
		id, _ := strconv.Atoi(IDArtist)
		if GroupieTracker.ArtisteNotFound(id, Main.ApiStruct) {
			http.Redirect(w, r, "/artiste", http.StatusFound)
			return
		}
		Main.ApiStruct.SpecificApiPageArtiste = GroupieTracker.ApiArtistsPageArtiste(IDArtist)
		locs := GroupieTracker.Mapapi(Main.ApiStruct, id)
		data := struct {
			Main MainStructure
			Locs [][]float64
		}{Main: *Main, Locs: locs}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "pages-artistes.html", data)
	})

	http.HandleFunc("/creation", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "creation.html", Main)
		GroupieTracker.GoodCreationReset(Main.AccStruct)
	})

	http.HandleFunc("/connection", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "connection.html", Main)
		// 	CheckConnection.Mail, CheckConnection.Pwd = false, false
	})

	http.HandleFunc("/checkcreation", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		mail, pwd := r.FormValue("mail"), r.FormValue("pwd")
		if GroupieTracker.VerifEntryUser(mail, pwd, r.FormValue("pwdc"), Main.AccStruct) {
			GroupieTracker.CreateAccount(mail, pwd, Main.AccStruct)
			GroupieTracker.AuthentificationToken(mail, Main.AccStruct, w)
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/creation", http.StatusFound)
	})

	http.HandleFunc("/checkconnection", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		// Login(w, r, CheckConnection, Main.ACC)
	})

	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("AUTHENTIFICATION_TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		// 	GroupieTracker.LoginAcc(cookie.Value, Acc)
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "profil.html", Main)
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
		// GroupieTracker.Test(w, r)
	})
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
	})
	http.ListenAndServe(":8080", nil)
}

// func Login(w http.ResponseWriter, r *http.Request, CC *GroupieTracker.CheckCo, Acc *GroupieTracker.Account) {
// 	mail, pwd := r.FormValue("mail"), r.FormValue("pwd")
// 	if GroupieTracker.CheckConnection(mail, pwd, CC, Acc) {
// 		SetCookie(w, mail, pwd, Acc)
// 		http.Redirect(w, r, "/profil", http.StatusFound)
// 	} else {
// 		http.Redirect(w, r, "/connection", http.StatusFound)
// 	}
// }
