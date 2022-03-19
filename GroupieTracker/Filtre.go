package GroupieTracker

import (
	"net/http"
	"strconv"
	"strings"
)

//filtre des Artsites

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
	CountryValue string
	GenresValue  string
	CountryTab   []string
	GenresTab    []string
	MembersTab   []string
	SearchBar    string
}

func FilterReset(ApiStruct *ApiStructure) {
	ApiStruct.Filtres.ArtD196069 = ""
	ApiStruct.Filtres.ArtD197079 = ""
	ApiStruct.Filtres.ArtD198089 = ""
	ApiStruct.Filtres.ArtD199099 = ""
	ApiStruct.Filtres.ArtD200009 = ""
	ApiStruct.Filtres.ArtD201019 = ""
	ApiStruct.Filtres.AlbD196069 = ""
	ApiStruct.Filtres.AlbD197079 = ""
	ApiStruct.Filtres.AlbD198089 = ""
	ApiStruct.Filtres.AlbD199099 = ""
	ApiStruct.Filtres.AlbD200009 = ""
	ApiStruct.Filtres.AlbD201019 = ""
	ApiStruct.Filtres.NbMember = "0"
	ApiStruct.Filtres.GenresValue = "All"
	ApiStruct.Filtres.CountryValue = "All"
	ApiStruct.Filtres.SearchBar = "artiste"
}

func FLT(filters map[string][]string, ApiStruct *ApiStructure, ATS *TokenSpotify) {
	ApiStruct.TabApiFiltre = []ApiArtiste{}
	FLTCheck(filters, ApiStruct)
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
		filters["Location"] = ApiStruct.Filtres.CountryTab
	}
	for _, i := range ApiStruct.TabApiArtiste {
		TabAppend(filters, ApiStruct, i)
	}
	if len(filters["Location"]) == len(ApiStruct.Filtres.CountryTab) {
		ApiStruct.Filtres.CountryValue = "All"
	}
	ApiStruct.Filtres.GenresValue = filters["genres"][0]
	if ApiStruct.Filtres.GenresValue != "All" {
		TabArtists := ApiStruct.TabApiFiltre
		ApiStruct.TabApiFiltre = []ApiArtiste{}
		for _, i := range TabArtists {
			Genres := RequestGenresById(ApiStruct.AllIdArtists[i.Id], ATS)
			if CheckIfInTab(ApiStruct.Filtres.GenresValue, Genres.Genres) {
				ApiStruct.TabApiFiltre = append(ApiStruct.TabApiFiltre, i)
			}
		}
	}
}

func FLTCheck(filters map[string][]string, ApiStruct *ApiStructure) {
	if filters["art_date"] != nil {
		for _, i := range filters["art_date"] {
			if i == "1960" {
				ApiStruct.Filtres.ArtD196069 = "checked"
			}
			if i == "1970" {
				ApiStruct.Filtres.ArtD197079 = "checked"
			}
			if i == "1980" {
				ApiStruct.Filtres.ArtD198089 = "checked"
			}
			if i == "1990" {
				ApiStruct.Filtres.ArtD199099 = "checked"
			}
			if i == "2000" {
				ApiStruct.Filtres.ArtD200009 = "checked"
			}
			if i == "2010" {
				ApiStruct.Filtres.ArtD201019 = "checked"
			}
		}
	}
	if filters["alb_date"] != nil {
		for _, i := range filters["alb_date"] {
			if i == "1960" {
				ApiStruct.Filtres.AlbD196069 = "checked"
			}
			if i == "1970" {
				ApiStruct.Filtres.AlbD197079 = "checked"
			}
			if i == "1980" {
				ApiStruct.Filtres.AlbD198089 = "checked"
			}
			if i == "1990" {
				ApiStruct.Filtres.AlbD199099 = "checked"
			}
			if i == "2000" {
				ApiStruct.Filtres.AlbD200009 = "checked"
			}
			if i == "2010" {
				ApiStruct.Filtres.AlbD201019 = "checked"
			}
		}
	}

	if filters["nb_member"][0] != "0" {
		ApiStruct.Filtres.NbMember = filters["nb_member"][0]
	}
	if filters["Location"][0] != "All" {
		ApiStruct.Filtres.CountryValue = filters["Location"][0]
	}
}

func TabCountry(ApiStruct *ApiStructure) {
	ApiStruct.Filtres.CountryTab = []string{"All"}
	for _, i := range ApiStruct.TabApiArtisteLocations {
		for _, k := range i.Locations {
			country := strings.Split(k, "-")[1]
			if !CheckIfInTab(country, ApiStruct.Filtres.CountryTab) {
				ApiStruct.Filtres.CountryTab = append(ApiStruct.Filtres.CountryTab, country)
			}
		}
	}
}

func CheckIfInTab(value string, TabValue []string) bool {
	for _, i := range TabValue {
		if i == value {
			return true
		}
	}
	return false
}

func TabAppend(filters map[string][]string, ApiStruct *ApiStructure, i ApiArtiste) {
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
							for _, n := range ApiStruct.TabApiArtisteLocations[i.Id-1].Locations {
								CountryLocation := strings.Split(n, "-")[1]
								for _, q := range filters["Location"] {
									if CountryLocation == q {
										ApiStruct.TabApiFiltre = append(ApiStruct.TabApiFiltre, i)
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

//Filtre de la "search bar"
func FiltreSearchBar(w http.ResponseWriter, r *http.Request, F *Filter) {
	filtre := r.FormValue("select-option")
	F.SearchBar = filtre
	http.Redirect(w, r, "/#second-page", http.StatusFound)
}

func SearchNameArtsit(value string, ApiStruct *ApiStructure) string {
	for _, i := range ApiStruct.TabApiArtiste {
		if i.Name == value {
			return strconv.Itoa(i.Id)
		}
	}
	return ""
}

func ArtisteNotFound(id int, ApiStruct *ApiStructure) bool {
	if id < 1 || id > ApiStruct.TabApiArtiste[len(ApiStruct.TabApiArtiste)-1].Id {
		return true
	}
	for _, i := range ApiStruct.TabApiArtiste {
		if i.Id == id {
			return false
		}
	}
	return true
}
