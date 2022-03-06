package GroupieTracker

import (
	"strconv"
)

type Filtre_Artist struct {
	ModifArtistCrea bool
	ModifDateAlbum  bool
	ModifNbGroup    bool
	ModifLocation   bool
	ValueArtist     string
	ValueAlbum      string
	ValueMember     string
	ValueLocation   string
	Test            string
}

func FiltreArtsites(F *Filtre_Artist, Api *Api) {
	var IntValueArtist int
	var TempoTab []Artist
	if F.ModifDateAlbum == true || F.ModifNbGroup == true || F.ModifLocation == true {
		F.ModifArtistCrea = true
		TempoTab = append(TempoTab, Api.ApiFiltre...)
		Api.ApiFiltre = []Artist{}
		//Pas encore trouvé un moyen claire de faire un all
		if F.ValueArtist == "all" {
			Api.ApiFiltre = append(Api.ApiFiltre, Api.ApiArtist...)
		} else {
			IntValueArtist, _ = strconv.Atoi(F.ValueArtist)
			for _, i := range TempoTab {
				if i.CreationDate >= IntValueArtist && i.CreationDate <= IntValueArtist+9 {
					Api.ApiFiltre = append(Api.ApiFiltre, i)
				}
			}
		}
	} else {
		Api.ApiFiltre = []Artist{}
		if F.ValueArtist == "all" {
			F.ModifArtistCrea = false
			for _, i := range Api.ApiArtist {
				Api.ApiFiltre = append(Api.ApiFiltre, i)
			}
		} else {
			F.ModifArtistCrea = true
			IntValueArtist, _ = strconv.Atoi(F.ValueArtist)
			for _, i := range Api.ApiArtist {
				if i.CreationDate >= IntValueArtist && i.CreationDate <= IntValueArtist+9 {
					Api.ApiFiltre = append(Api.ApiFiltre, i)
				}
			}
		}
	}
}
