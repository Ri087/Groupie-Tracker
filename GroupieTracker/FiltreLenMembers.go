package GroupieTracker

import (
	"strconv"
)

func FiltreLenMembers(F *Filtre_Artist, Api *Api) {
	var TempoTab []Artist
	var IntValueMembers int
	if F.ModifArtistCrea == true || F.ModifDateAlbum == true || F.ModifLocation == true {
		TempoTab = append(TempoTab, Api.ApiFiltre...)
		Api.ApiFiltre = []Artist{}
		if F.ValueArtist == "all" {
			Api.ApiFiltre = append(Api.ApiFiltre, Api.ApiArtist...)
		} else {
			IntValueMembers, _ = strconv.Atoi(F.ValueMember)
			for _, i := range TempoTab {
				if len(i.Members) == IntValueMembers {
					Api.ApiFiltre = append(Api.ApiFiltre, i)
				}
			}
		}
	} else {
		Api.ApiFiltre = []Artist{}
		if F.ValueAlbum == "all" {
			F.ModifNbGroup = false
			for _, i := range Api.ApiArtist {
				Api.ApiFiltre = append(Api.ApiFiltre, i)
			}
		} else {
			F.ModifNbGroup = true
			IntValueMembers, _ = strconv.Atoi(F.ValueMember)
			for _, i := range Api.ApiArtist {
				if len(i.Members) == IntValueMembers {
					Api.ApiFiltre = append(Api.ApiFiltre, i)
				}
			}
		}
	}
}
