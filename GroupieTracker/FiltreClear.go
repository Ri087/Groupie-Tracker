package GroupieTracker

import "net/http"

func FiltreClear(w http.ResponseWriter, r *http.Request, F *Filtre_Artist, Api *Api) {
	Api.ApiFiltre = []Artist{}
	Api.ApiFiltre = append(Api.ApiFiltre, Api.ApiArtist...)
	http.Redirect(w, r, "/artiste", http.StatusFound)
}
