package main

import (
	"GroupieTracker/GroupieTracker"
	"encoding/base32"
	"net/http"
	"strconv"
	"text/template"
)

type MainStructure struct {
	ApiStruct *GroupieTracker.ApiStructure
}

func MainStructureInit() *MainStructure {
	Main := &MainStructure{}
	Main.ApiStruct = GroupieTracker.ApiStructInit()

	return Main
}

func main() {
	MainStructureMain := MainStructureInit()
	// Acc := &GroupieTracker.Account{}
	// CheckCreation := &GroupieTracker.CheckCreation{}
	// CheckConnection := &GroupieTracker.CheckCo{}
	fileServer := http.FileServer(http.Dir("./static"))
	var s GroupieTracker.Spotify = GroupieTracker.New("6b053d7dfcbe4c69a576561f8c098391", "d00791e8792a4f13bc1bb8b95197505d")
	s.Authorize()
	http.Handle("/ressources/", http.StripPrefix("/ressources/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		MainStructureMain.ApiStruct.TabApiArtiste, MainStructureMain.ApiStruct.TabApiArtisteLocations = GroupieTracker.ApiArtistsArtiste()
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "index.html", MainStructureMain)
	})

	//Page principal
	http.HandleFunc("/artiste", func(w http.ResponseWriter, r *http.Request) {
		MainStructureMain.ApiStruct.TabApiArtiste, MainStructureMain.ApiStruct.TabApiArtisteLocations = GroupieTracker.ApiArtistsArtiste()
		GroupieTracker.TabCountry(MainStructureMain.ApiStruct)
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "artiste.html", MainStructureMain)
		GroupieTracker.FilterReset(MainStructureMain.ApiStruct)
	})
	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		GroupieTracker.FLT(r.URL.Query(), MainStructureMain.ApiStruct)
		http.Redirect(w, r, "/artiste", http.StatusFound)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		id := GroupieTracker.SearchNameArtsit(r.FormValue("search-artist"), MainStructureMain.ApiStruct)
		if id == "" {
			http.Redirect(w, r, "/#second-page", http.StatusFound)
		} else {
			http.Redirect(w, r, "/artiste/"+id, http.StatusFound)
		}
	})

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "event.html", MainStructureMain)
	})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "contact.html", MainStructureMain)
	})

	http.HandleFunc("/artiste/", func(w http.ResponseWriter, r *http.Request) {
		IDArtist := r.URL.Path[9:]
		MainStructureMain.ApiStruct.SpecificApiPageArtiste = GroupieTracker.ApiArtistsPageArtiste(IDArtist)
		id, _ := strconv.Atoi(r.URL.Path[len(r.URL.Path)-1:])
		locs := GroupieTracker.Mapapi(MainStructureMain.ApiStruct, id)
		data := struct {
			Main MainStructure
			Locs [][]float64
		}{Main: *MainStructureMain, Locs: locs}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "pages-artistes.html", data)
	})
	// http.HandleFunc("/connection", func(w http.ResponseWriter, r *http.Request) {
	// 	if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
	// 		http.Redirect(w, r, "/profil", http.StatusFound)
	// 		return
	// 	}
	// 	var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
	// 	templateshtml.ExecuteTemplate(w, "connection.html", CheckConnection)
	// 	CheckConnection.Mail, CheckConnection.Pwd = false, false
	// })
	// http.HandleFunc("/creation", func(w http.ResponseWriter, r *http.Request) {
	// 	if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
	// 		http.Redirect(w, r, "/profil", http.StatusFound)
	// 		return
	// 	}
	// 	var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
	// 	templateshtml.ExecuteTemplate(w, "creation.html", CheckCreation)
	// 	CheckCreation.Name, CheckCreation.Pwd, CheckCreation.Pwdc, CheckCreation.Mail, CheckCreation.Exist = false, false, false, false, false
	// })
	// http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
	// 	cookie, err := r.Cookie("AUTHENTIFICATION_TOKEN")
	// 	if err != nil {
	// 		http.Redirect(w, r, "/connection", http.StatusFound)
	// 		return
	// 	}
	// 	GroupieTracker.LoginAcc(cookie.Value, Acc)
	// 	var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
	// 	templateshtml.ExecuteTemplate(w, "profil.html", Main.ACC)
	// 	Logout(Acc)
	// })

	// http.HandleFunc("/checkcreation", func(w http.ResponseWriter, r *http.Request) {
	// 	Creation(w, r, CheckCreation, Main.ACC)
	// })

	// http.HandleFunc("/checkconnection", func(w http.ResponseWriter, r *http.Request) {
	// 	Login(w, r, CheckConnection, Main.ACC)
	// })

	// http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
	// 	if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
	// 		cookie, _ := r.Cookie("AUTHENTIFICATION_TOKEN")
	// 		cookie.MaxAge = -1
	// 		http.SetCookie(w, cookie)
	// 	}
	// 	http.Redirect(w, r, "/", http.StatusFound)
	// })

	// NE PAS SUPPR CEST DES TESTS
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// GroupieTracker.Test(w, r)
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
