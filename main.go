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
	Main.AccStruct = GroupieTracker.AccountStructureInit()
	var s = GroupieTracker.New("6b053d7dfcbe4c69a576561f8c098391", "d00791e8792a4f13bc1bb8b95197505d")
	Main.Token = s.Authorize()
	return Main
}

func main() {
	Main := MainStructureInit()
	go GenerateSpotifyToken(Main)

	// GroupieTracker.TabGenres(Main.ApiStruct, Main.Token)

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
		if cookie, err := r.Cookie("TOKEN"); err == nil && Main.AccStruct.AllToken[cookie.Value] != "" {
			GroupieTracker.LoadUserByToken(cookie.Value, Main.AccStruct)
			if Main.AccStruct.ProfilParameters.Profil.User.ArtistsLiked[IDArtist] {
				Main.AccStruct.ProfilParameters.Artists.ArtistCheck = true
			}
		}
		Main.ApiStruct.SpecificApiPageArtiste = GroupieTracker.ApiArtistsPageArtiste(IDArtist, Main.Token)
		locs := GroupieTracker.Mapapi(Main.ApiStruct, id-1)
		data := struct {
			Main MainStructure
			Locs [][]float64
		}{Main: *Main, Locs: locs}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "pages-artistes.html", data)
		GroupieTracker.ArtistsProfilReset(Main.AccStruct)
	})

	http.HandleFunc("/creation", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err == nil {
			if Main.AccStruct.AllToken[cookie.Value] == "" {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
			} else {
				http.Redirect(w, r, "/profil", http.StatusFound)
				return
			}
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "creation.html", Main)
		GroupieTracker.CreationVerificationReset(Main.AccStruct)
	})

	http.HandleFunc("/checkcreation", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err == nil {
			if Main.AccStruct.AllToken[cookie.Value] == "" {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
			} else {
				http.Redirect(w, r, "/profil", http.StatusFound)
				return
			}
		}
		GroupieTracker.CreationUserInformationFill(r.FormValue("mail"), r.FormValue("pwd"), r.FormValue("pwdc"), Main.AccStruct)
		if GroupieTracker.CreationVerificationEntryUser(Main.AccStruct) {
			GroupieTracker.AccountCreation(Main.AccStruct)
			GroupieTracker.SetAuthentificationToken(w, Main.AccStruct)
			GroupieTracker.CreationUserInformationReset(Main.AccStruct)
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		GroupieTracker.CreationUserInformationReset(Main.AccStruct)
		http.Redirect(w, r, "/creation", http.StatusFound)
	})

	http.HandleFunc("/connection", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err == nil {
			if Main.AccStruct.AllToken[cookie.Value] == "" {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
			} else {
				http.Redirect(w, r, "/profil", http.StatusFound)
				return
			}
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "connection.html", Main)
		GroupieTracker.ConnectionVerificationReset(Main.AccStruct)
	})

	http.HandleFunc("/checkconnection", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err == nil {
			if Main.AccStruct.AllToken[cookie.Value] == "" {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
			} else {
				http.Redirect(w, r, "/profil", http.StatusFound)
				return
			}
		}
		GroupieTracker.ConnectionUserInformationFill(r.FormValue("mail"), r.FormValue("pwd"), Main.AccStruct)
		if GroupieTracker.ConnectionVerificationEntryUser(Main.AccStruct) {
			GroupieTracker.SetAuthentificationToken(w, Main.AccStruct)
			GroupieTracker.ConnectionUserInformationReset(Main.AccStruct)
			http.Redirect(w, r, "/profil", http.StatusFound)
			return
		}
		GroupieTracker.ConnectionUserInformationReset(Main.AccStruct)
		http.Redirect(w, r, "/connection", http.StatusFound)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err == nil {
			var tab []string
			for _, i := range Main.AccStruct.AllTabToken[Main.AccStruct.AllToken[cookie.Value]] {
				if i != cookie.Value {
					tab = append(tab, i)
				}
			}
			if len(tab) == 0 {
				delete(Main.AccStruct.AllTabToken, Main.AccStruct.AllToken[cookie.Value])
			} else {
				Main.AccStruct.AllTabToken[Main.AccStruct.AllToken[cookie.Value]] = tab
			}
			delete(Main.AccStruct.AllToken, cookie.Value)
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			GroupieTracker.SaveAllToken(Main.AccStruct)
			GroupieTracker.SaveAllTabToken(Main.AccStruct)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		} else if Main.AccStruct.AllToken[cookie.Value] == "" {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		GroupieTracker.LoadUserByToken(cookie.Value, Main.AccStruct)
		if Main.AccStruct.ProfilParameters.Profil.User.Name != "" {
			Main.AccStruct.ProfilParameters.Name.NameIsOk = true
			for _, i := range Main.AccStruct.ProfilParameters.Profil.User.Friends {
				Main.AccStruct.ProfilParameters.Friends.Friends = append(Main.AccStruct.ProfilParameters.Friends.Friends, GroupieTracker.FriendsInformations{Main.AccStruct.AllAccount[i].User, i})
			}
		}
		GroupieTracker.ArtistsProfilFill(Main.AccStruct)
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "profil.html", Main)
		GroupieTracker.ProfilAccountReset(Main.AccStruct)
		GroupieTracker.NameParametersReset(Main.AccStruct)
		GroupieTracker.FriendsParametersReset(Main.AccStruct)
		GroupieTracker.ArtistsProfilReset(Main.AccStruct)
	})

	http.HandleFunc("/pseudo", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		} else if Main.AccStruct.AllToken[cookie.Value] == "" {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		GroupieTracker.LoadUserByToken(cookie.Value, Main.AccStruct)
		if Main.AccStruct.ProfilParameters.Profil.User.EntireName == "" {
			GroupieTracker.ProfilUserInformationFill(r.FormValue("name"), r.FormValue("hashtag"), Main.AccStruct)
			if GroupieTracker.NameVerificationEntryUser(Main.AccStruct) {
				Main.AccStruct.AllAccount[Main.AccStruct.AllToken[cookie.Value]] = Main.AccStruct.ProfilParameters.Profil
				Main.AccStruct.AllHashtag[Main.AccStruct.ProfilParameters.Profil.User.EntireName] = Main.AccStruct.AllToken[cookie.Value]
			}
		} else {
			ToSuppr := Main.AccStruct.ProfilParameters.Profil.User.EntireName
			GroupieTracker.ProfilUserInformationFill(r.FormValue("name"), r.FormValue("hashtag"), Main.AccStruct)
			if GroupieTracker.NameVerificationEntryUser(Main.AccStruct) {
				delete(Main.AccStruct.AllHashtag, ToSuppr)
				Main.AccStruct.AllAccount[Main.AccStruct.AllToken[cookie.Value]] = Main.AccStruct.ProfilParameters.Profil
				Main.AccStruct.AllHashtag[Main.AccStruct.ProfilParameters.Profil.User.EntireName] = Main.AccStruct.AllToken[cookie.Value]
			}
		}
		GroupieTracker.SaveAllAccount(Main.AccStruct)
		GroupieTracker.SaveAllHashtag(Main.AccStruct)
		GroupieTracker.ProfilAccountReset(Main.AccStruct)
		http.Redirect(w, r, "/profil", http.StatusFound)
	})

	http.HandleFunc("/addfriend", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		} else if Main.AccStruct.AllToken[cookie.Value] == "" {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		if Main.AccStruct.AllHashtag[r.FormValue("name")] != "" {
			FriendsId := Main.AccStruct.AllHashtag[r.FormValue("name")]
			if Main.AccStruct.AllToken[cookie.Value] == FriendsId {
				Main.AccStruct.ProfilParameters.Friends.ThatsU = true
				http.Redirect(w, r, "/profil", http.StatusFound)
				return
			}
			GroupieTracker.LoadUserByToken(cookie.Value, Main.AccStruct)
			for _, i := range Main.AccStruct.ProfilParameters.Profil.User.Friends {
				if i == FriendsId {
					Main.AccStruct.ProfilParameters.Friends.AlreadyHave = true
					http.Redirect(w, r, "/profil", http.StatusFound)
					return
				}
			}
			Main.AccStruct.ProfilParameters.Profil.User.Friends = append(Main.AccStruct.ProfilParameters.Profil.User.Friends, FriendsId)
			Main.AccStruct.AllAccount[Main.AccStruct.AllToken[cookie.Value]] = Main.AccStruct.ProfilParameters.Profil
			GroupieTracker.SaveAllAccount(Main.AccStruct)
		} else {
			Main.AccStruct.ProfilParameters.Friends.NotExist = true
		}
		GroupieTracker.ProfilAccountReset(Main.AccStruct)
		http.Redirect(w, r, "/profil", http.StatusFound)
	})

	http.HandleFunc("/profil/", func(w http.ResponseWriter, r *http.Request) {
		IdUser := r.URL.Path[8:]
		Main.AccStruct.ProfilParameters.Profil = Main.AccStruct.AllAccount[IdUser]
		if Main.AccStruct.ProfilParameters.Profil.User.EntireName != "" {
			Main.AccStruct.ProfilParameters.Visit.Exist = true
			cookie, err := r.Cookie("TOKEN")
			if err == nil && Main.AccStruct.AllToken[cookie.Value] == "" {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
			} else if err == nil {
				GroupieTracker.ProfilAuthorizeVisit(Main.AccStruct.AllToken[cookie.Value], Main.AccStruct)
			}
			GroupieTracker.ProfilAuthorizeVisit("", Main.AccStruct)
		}
		var templateshtml = template.Must(template.ParseGlob("./static/html/*.html"))
		templateshtml.ExecuteTemplate(w, "profil-visite.html", Main)
		GroupieTracker.ProfilAccountReset(Main.AccStruct)
		GroupieTracker.ProfilVisitReset(Main.AccStruct)
	})

	http.HandleFunc("/showprofil", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		} else if Main.AccStruct.AllToken[cookie.Value] == "" {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		value := r.FormValue("Profil")
		GroupieTracker.LoadUserByToken(cookie.Value, Main.AccStruct)
		GroupieTracker.ProfilSettings(value, Main.AccStruct)
		Main.AccStruct.AllAccount[Main.AccStruct.AllToken[cookie.Value]] = Main.AccStruct.ProfilParameters.Profil
		GroupieTracker.SaveAllAccount(Main.AccStruct)
		http.Redirect(w, r, "/profil", http.StatusFound)
	})

	http.HandleFunc("/liked", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("TOKEN")
		if err != nil {
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		} else if Main.AccStruct.AllToken[cookie.Value] == "" {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/connection", http.StatusFound)
			return
		}
		IdArtist := r.Referer()[30:]
		GroupieTracker.LoadUserByToken(cookie.Value, Main.AccStruct)
		Main.AccStruct.ProfilParameters.Profil.User.ArtistsLiked[IdArtist] = !Main.AccStruct.ProfilParameters.Profil.User.ArtistsLiked[IdArtist]
		Main.AccStruct.AllAccount[Main.AccStruct.AllToken[cookie.Value]] = Main.AccStruct.ProfilParameters.Profil
		GroupieTracker.SaveAllAccount(Main.AccStruct)
		GroupieTracker.ProfilAccountReset(Main.AccStruct)
		http.Redirect(w, r, r.Referer(), http.StatusFound)
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
