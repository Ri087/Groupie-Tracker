package GroupieTracker

import (
	"fmt"
	"net/http"
	"strconv"
)

type Filtre_Artist struct {
	ModifArtistCrea bool
	ModifDateAlbum  bool
	DateStart       int
	DateEnd         int
	DateAlbumS      int
	DateAlbumE      int
	DateAlbumInt    int
}

func FuncFiltreDate(w http.ResponseWriter, r *http.Request, Fa *Filtre_Artist) {
	date_filtre := r.FormValue("filtre_date")
	if len(date_filtre) < 1 {
		Fa.ModifArtistCrea = false
		Fa.DateStart = 0
		Fa.DateEnd = 2100
	} else {
		Fa.ModifArtistCrea = true
		if date_filtre == "all" {
			Fa.DateStart = 0
			Fa.DateEnd = 2100
		} else {
			date, _ := strconv.Atoi(date_filtre)
			Fa.DateStart = date
			Fa.DateEnd = date + 9
			fmt.Println("* ", Fa.DateStart, "|", Fa.DateEnd)
		}
	}
	http.Redirect(w, r, "/artiste", http.StatusFound)
}

// Date de creation du premier album
// func FuncFiltreAlbumCrea(w http.ResponseWriter, r *http.Request, Fa *Filtre_Artist) {
// 	date_filtre := r.FormValue("filtre-date-album")
// 	if len(date_filtre) < 1 {
// 		Fa.ModifDateAlbum = false
// 		Fa.DateAlbumS = 0
// 		Fa.DateAlbumE = 2100
// 	} else {
// 		Fa.ModifDateAlbum = true
// 		if date_filtre == "all" {
// 			Fa.DateAlbumS = 0
// 			Fa.DateAlbumE = 2100
// 		} else {
// 			date, _ := strconv.Atoi(date_filtre)
// 			Fa.DateAlbumS = date
// 			Fa.DateAlbumE = date + 9
// 			fmt.Println("**", Fa.DateAlbumS, "|", Fa.DateAlbumE)
// 		}
// 	}
// 	http.Redirect(w, r, "/artiste", http.StatusFound)
// }
