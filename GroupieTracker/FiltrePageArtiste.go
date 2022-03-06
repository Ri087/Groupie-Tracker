package GroupieTracker

import (
	"net/http"
)

func Filtre(w http.ResponseWriter, r *http.Request, F *Filtre_Artist, Api *Api) {
	rFormArtist := r.FormValue("filtre_date_artiste")
	rFormAlbum := r.FormValue("filtre_date_album")
	rFormMembers := r.FormValue("filtre_groupe")
	// F.ValueLocation = r.FormValue("filtre_location")
	if len(rFormArtist) > 0 {
		F.ValueArtist = rFormArtist
	}
	if len(rFormAlbum) > 0 {
		F.ValueAlbum = rFormAlbum
	}
	if len(rFormMembers) > 0 {
		F.ValueMember = rFormMembers
		F.Test = rFormMembers
	}
	if len(F.ValueArtist) != 0 {
		FiltreArtsites(F, Api)
	}
	if len(F.ValueAlbum) > 0 {
		FiltreAlbums(F, Api)
	}
	if len(F.ValueMember) > 0 {
		FiltreLenMembers(F, Api)
	}
	// if len(F.ValueLocation) != 0 {

	// }
	http.Redirect(w, r, "/artiste", http.StatusFound)
}
