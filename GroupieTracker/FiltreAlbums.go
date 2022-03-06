package GroupieTracker

import (
	"strconv"
	"strings"
)

func FiltreAlbums(F *Filtre_Artist, Api *Api) {
	var IntValueAlbum int
	var TempoTab []Artist
	if F.ModifArtistCrea == true || F.ModifNbGroup == true || F.ModifLocation == true {
		TempoTab = append(TempoTab, Api.ApiFiltre...)
		Api.ApiFiltre = []Artist{}
		//Pas encore trouvé un moyen claire de faire un all
		if F.ValueArtist == "all" {
			Api.ApiFiltre = append(Api.ApiFiltre, Api.ApiArtist...)
		} else {
			IntValueAlbum, _ = strconv.Atoi(F.ValueAlbum)
			for _, i := range TempoTab {
				split := strings.Split(i.FirstAlbum, "-")
				splitInt, _ := strconv.Atoi(split[2])
				if splitInt >= IntValueAlbum && splitInt <= IntValueAlbum+9 {
					Api.ApiFiltre = append(Api.ApiFiltre, i)
				}
			}
		}
	} else {
		Api.ApiFiltre = []Artist{}
		if F.ValueAlbum == "all" {
			F.ModifDateAlbum = false
			for _, i := range Api.ApiArtist {
				Api.ApiFiltre = append(Api.ApiFiltre, i)
			}
		} else {
			F.ModifDateAlbum = true
			IntValueAlbum, _ = strconv.Atoi(F.ValueAlbum)
			for _, i := range Api.ApiArtist {
				split := strings.Split(i.FirstAlbum, "-")
				splitInt, _ := strconv.Atoi(split[2])
				if splitInt >= IntValueAlbum && splitInt <= IntValueAlbum+9 {
					Api.ApiFiltre = append(Api.ApiFiltre, i)
				}
			}
		}
	}
}
