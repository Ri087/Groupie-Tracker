package main

import (
	"GroupieTracker/GroupieTracker"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type MainStructure struct {
	ApiStruct *GroupieTracker.ApiStructure
	AccStruct *GroupieTracker.AccountStruct
	Token     *GroupieTracker.TokenSpotify
}

func MainStructureInit() *MainStructure {
	Main := &MainStructure{}
	Main.ApiStruct = GroupieTracker.ApiStructInit()
	Main.AccStruct = GroupieTracker.AccStructureInit()
	var s = GroupieTracker.New("6b053d7dfcbe4c69a576561f8c098391", "d00791e8792a4f13bc1bb8b95197505d")
	Main.Token = s.Authorize()
	return Main
}

func main() {
	Main := MainStructureInit()
	go GenerateSpotifyToken(Main)

	GroupieTracker.TabGenres(Main.ApiStruct, Main.Token)

	fileServer := http.FileServer(http.Dir("./static"))
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
	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		GroupieTracker.FLT(r.URL.Query(), Main.ApiStruct, Main.Token)
		http.Redirect(w, r, "/artiste", http.StatusFound)

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
		Main.ApiStruct.SpecificApiPageArtiste = GroupieTracker.ApiArtistsPageArtiste(IDArtist, Main.Token)
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
		GroupieTracker.GoodConnectionReset(Main.AccStruct)
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
		mail, pwd := r.FormValue("mail"), r.FormValue("pwd")
		if GroupieTracker.VerifConnectionUser(mail, pwd, Main.AccStruct) {
			GroupieTracker.AuthentificationToken(mail, Main.AccStruct, w)
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/connection", http.StatusFound)

	})

	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("AUTHENTIFICATION_TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		if Main.AccStruct.AuthToken[cookie.Value] == "" {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		GroupieTracker.GetUserInfos(cookie.Value, Main.AccStruct)
		if Main.AccStruct.User.Pseudo == "" {
			Main.AccStruct.PseudoCheck.PseudoNotOk = true
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "profil.html", Main)
		GroupieTracker.PseudoAndFriendReset(Main.AccStruct)
	})

	http.HandleFunc("/pseudo", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("AUTHENTIFICATION_TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		if Main.AccStruct.AuthToken[cookie.Value] == "" {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		TempPseudo := r.FormValue("pseudo")
		if len(TempPseudo) < 3 {
			Main.AccStruct.PseudoCheck.WrongPseudo = true
		} else {
			GroupieTracker.SavePseudo(cookie.Value, TempPseudo, Main.AccStruct)
		}
		http.Redirect(w, r, "/profil", http.StatusFound)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
			cookie, _ := r.Cookie("AUTHENTIFICATION_TOKEN")
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("/profil/", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Path[8:])
		Main.AccStruct.AuthorizeVisit.User = Main.AccStruct.EveryUserInfos[id]
		if Main.AccStruct.AuthorizeVisit.User.Mail != "" {
			Main.AccStruct.AuthorizeVisit.Existant = true
			if cookie, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil && Main.AccStruct.AuthToken[cookie.Value] == "" {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
			}
			if cookie, err := r.Cookie("AUTHENTIFICATION_TOKEN"); err == nil {
				GroupieTracker.GetUserInfos(cookie.Value, Main.AccStruct)
			} else {
				Main.AccStruct.User = GroupieTracker.InfosUser{}
			}
			if Main.AccStruct.AuthorizeVisit.User.Mail == Main.AccStruct.User.Mail {
				http.Redirect(w, r, "/profil", http.StatusFound)
				return
			}
			GroupieTracker.VisitProfil(Main.AccStruct)
			if Main.AccStruct.AuthorizeVisit.Authorize {
				GroupieTracker.ShowedFriends(Main.AccStruct)
			}
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "profil-visite.html", Main)
		GroupieTracker.VisitAuthorizeReset(Main.AccStruct)
	})

	http.HandleFunc("/addfriend", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("AUTHENTIFICATION_TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		if Main.AccStruct.AuthToken[cookie.Value] == "" {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		id := Main.AccStruct.EveryId[r.FormValue("mail")]
		if id == 0 {
			Main.AccStruct.FriendCheck.WrongFriend = true
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		GroupieTracker.GetFriendById(id, Main.AccStruct)
		GroupieTracker.GetUserInfos(cookie.Value, Main.AccStruct)
		if Main.AccStruct.User.Mail == Main.AccStruct.Friend.Mail {
			Main.AccStruct.FriendCheck.ThatsU = true
		} else if Main.AccStruct.Friend.Pseudo != "" {
			GroupieTracker.AddFriend(id, Main.AccStruct)
		}
		http.Redirect(w, r, "/profil", http.StatusFound)
	})

	http.HandleFunc("/showprofil", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("AUTHENTIFICATION_TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		if Main.AccStruct.AuthToken[cookie.Value] == "" {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		GroupieTracker.GetUserInfos(cookie.Value, Main.AccStruct)
		GroupieTracker.ParametersProfil(r.FormValue("showprofil"), r.FormValue("showprofilfriend"), Main.AccStruct)
		http.Redirect(w, r, "/profil", http.StatusFound)
	})
	http.ListenAndServe(":8080", nil)
}

func GenerateSpotifyToken(Main *MainStructure) {
	for {
		time.Sleep(time.Duration(Main.Token.Expires_in) * time.Second)
		var s = GroupieTracker.New("6b053d7dfcbe4c69a576561f8c098391", "d00791e8792a4f13bc1bb8b95197505d")
		Main.Token = s.Authorize()
	}
}
