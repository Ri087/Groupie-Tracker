package GroupieTracker

import (
	"strconv"
	"strings"
)

type Filter struct {
	ArtD196069   string
	ArtD197079   string
	ArtD198089   string
	ArtD199099   string
	ArtD200009   string
	ArtD201019   string
	AlbD196069   string
	AlbD197079   string
	AlbD198089   string
	AlbD199099   string
	AlbD200009   string
	AlbD201019   string
	NbMember     string
	Country_tab  []string
	CountryValue string
}

func FilterReset(ADF *Filter) {
	ADF.ArtD196069 = ""
	ADF.ArtD197079 = ""
	ADF.ArtD198089 = ""
	ADF.ArtD199099 = ""
	ADF.ArtD200009 = ""
	ADF.ArtD201019 = ""
	ADF.AlbD196069 = ""
	ADF.AlbD197079 = ""
	ADF.AlbD198089 = ""
	ADF.AlbD199099 = ""
	ADF.AlbD200009 = ""
	ADF.AlbD201019 = ""
	ADF.NbMember = "0"
	ADF.CountryValue = "All"
}

func FLT(filters map[string][]string, Apis *Api, ADF *Filter) {
	Apis.ApiFiltre = []Artist{}
	FLTCheck(filters, ADF)
	if filters["art_date"] == nil {
		filters["art_date"] = []string{"1960", "1970", "1980", "1990", "2000", "2010"}
	}
	if filters["alb_date"] == nil {
		filters["alb_date"] = []string{"1960", "1970", "1980", "1990", "2000", "2010"}
	}
	if filters["nb_member"][0] == "0" {
		filters["nb_member"] = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	}
	if filters["Location"][0] == "All" {
		filters["Location"] = ADF.Country_tab
	}
	for _, i := range Apis.ApiArtist {
		ntm(filters, Apis, ADF, i)
	}
	if len(filters["Location"]) == len(ADF.Country_tab) {
		ADF.CountryValue = "All"
	}
}

func FLTCheck(filters map[string][]string, ADF *Filter) {
	if filters["art_date"] != nil {
		for _, i := range filters["art_date"] {
			if i == "1960" {
				ADF.ArtD196069 = "checked"
			}
			if i == "1970" {
				ADF.ArtD197079 = "checked"
			}
			if i == "1980" {
				ADF.ArtD198089 = "checked"
			}
			if i == "1990" {
				ADF.ArtD199099 = "checked"
			}
			if i == "2000" {
				ADF.ArtD200009 = "checked"
			}
			if i == "2010" {
				ADF.ArtD201019 = "checked"
			}
		}
	}
	if filters["alb_date"] != nil {
		for _, i := range filters["alb_date"] {
			if i == "1960" {
				ADF.AlbD196069 = "checked"
			}
			if i == "1970" {
				ADF.AlbD197079 = "checked"
			}
			if i == "1980" {
				ADF.AlbD198089 = "checked"
			}
			if i == "1990" {
				ADF.AlbD199099 = "checked"
			}
			if i == "2000" {
				ADF.AlbD200009 = "checked"
			}
			if i == "2010" {
				ADF.AlbD201019 = "checked"
			}
		}
	}
	if filters["nb_member"][0] != "0" {
		ADF.NbMember = filters["nb_member"][0]
	}
	if filters["Location"][0] != "All" {
		ADF.CountryValue = filters["Location"][0]
	}
}

func CountryTab(Api *Api, F *Filter) {
	F.Country_tab = append(F.Country_tab, "All")
	for _, i := range Api.ApiLocations {
		for _, o := range i.Locations {
			country := strings.Split(o, "-")[1]
			if !CheckIfInTab(country, F) {
				F.Country_tab = append(F.Country_tab, country)
			}
		}
	}
}

func CheckIfInTab(country string, F *Filter) bool {
	for _, y := range F.Country_tab {
		if country == y {
			return true
		}
	}
	return false
}

func ntm(filters map[string][]string, Apis *Api, ADF *Filter, i Artist) {
	for _, k := range filters["art_date"] {
		artiste_date, _ := strconv.Atoi(k)
		if artiste_date <= i.CreationDate && i.CreationDate <= artiste_date+9 {
			for _, l := range filters["alb_date"] {
				albumCreationDate, _ := strconv.Atoi(strings.Split(i.FirstAlbum, "-")[2])
				album_date, _ := strconv.Atoi(l)
				if album_date <= albumCreationDate && albumCreationDate <= album_date+9 {
					for _, m := range filters["nb_member"] {
						nb_member, _ := strconv.Atoi(m)
						if len(i.Members) == nb_member {
							for _, n := range Apis.ApiLocations[i.Id-1].Locations {
								CountryLocation := strings.Split(n, "-")[1]
								for _, q := range filters["Location"] {
									if CountryLocation == q {
										Apis.ApiFiltre = append(Apis.ApiFiltre, i)
										return
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
