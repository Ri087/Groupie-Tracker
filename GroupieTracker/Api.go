package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ApiStructure struct {
	TabApiArtiste          []ApiArtiste
	TabApiArtisteLocations []ApiArtisteLocations
	TabApiFiltre           []ApiArtiste
	AllIdArtists           map[int]string
	AllArtistsId           map[string]int
	ContactApi             ApiContacts
	Filtres                Filter
	SpecificApiPageArtiste ArtistsApiPageArtiste
	Top3Artists            [3]ApiTop3Artists
}

func ApiStructInit() *ApiStructure {
	ApiStruct := &ApiStructure{}
	ApiStruct.TabApiArtiste, ApiStruct.TabApiArtisteLocations = ApiArtistsArtiste()
	ApiStruct.TabApiFiltre, _ = ApiArtistsArtiste()
	ApiStruct.Filtres = Filter{}
	FilterReset(ApiStruct)
	ApiStruct.SpecificApiPageArtiste = ArtistsApiPageArtiste{}
	return ApiStruct
}

func GetReadAll(ApiLink string) []byte {
	response, err := http.Get(ApiLink)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	ApiResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return ApiResponse
}

type ApiWApi map[string]string

func LinkApi() ApiWApi {
	ApiLink := ApiWApi{}
	json.Unmarshal(GetReadAll("https://groupietrackers.herokuapp.com/api"), &ApiLink)
	return ApiLink
}

// Artiste

type StructApiArtiste struct {
	TabApiArtiste          []ApiArtiste
	TabApiArtisteLocations []ApiArtisteLocations
}

type ApiArtiste struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
}

type ApiArtisteLocations struct {
	Id        int
	Locations []string
}

type ApiArtisteLocationsIndex struct {
	Index []ApiArtisteLocations
}

func ApiArtistsArtiste() ([]ApiArtiste, []ApiArtisteLocations) {
	ApiArtistsIndex := ApiArtisteLocationsIndex{}
	ApiArtists := StructApiArtiste{}

	json.Unmarshal(GetReadAll((LinkApi()["artists"])), &ApiArtists.TabApiArtiste)
	json.Unmarshal(GetReadAll((LinkApi()["locations"])), &ApiArtistsIndex)

	return ApiArtists.TabApiArtiste, ApiArtistsIndex.Index
}

// Page Artiste

type ArtistsApiPageArtiste struct {
	Artists   ApiPageArtiste
	Locations ApiPageArtisteLocations
	Dates     ApiPageArtisteDates
	Relations ApiPageArtisteRelations
	Spotify   SpotifyPageArtiste
}

type ApiPageArtiste struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

type ApiPageArtisteLocations struct {
	Id        int
	Locations []string
}

type ApiPageArtisteDates struct {
	Id    int
	Dates []string
}

type ApiPageArtisteRelations struct {
	Id             int
	DatesLocations map[string][]string
}

func ApiArtistsPageArtiste(AllId map[int]string, id string, idint int, Token *TokenSpotify) ArtistsApiPageArtiste {
	ApiArtists := ArtistsApiPageArtiste{}
	json.Unmarshal(GetReadAll(LinkApi()["artists"]+"/"+id), &ApiArtists.Artists)
	json.Unmarshal(GetReadAll(ApiArtists.Artists.Locations), &ApiArtists.Locations)
	json.Unmarshal(GetReadAll(ApiArtists.Artists.ConcertDates), &ApiArtists.Dates)
	json.Unmarshal(GetReadAll(ApiArtists.Artists.Relations), &ApiArtists.Relations)
	ApiArtists.Spotify = PageArtistSpotify(AllId, ApiArtists, idint, Token)

	return ApiArtists
}

// Top 3

type ApiTop3Artists struct {
	Info    ApiArtistInfo
	Spotify SpotifyArtistInfo
}

type SpotifyArtistInfo struct {
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	ID string `json:"id"`
}

type ApiArtistInfo struct {
	Id    int
	Image string
	Name  string
}

func GenerateTop3Artists(ApiStruct *ApiStructure, ATS *TokenSpotify) {
	SpotInfoTab := []SpotifyArtistInfo{}
	for _, k := range ApiStruct.AllIdArtists {
		SpotInfo := SpotifyArtistInfo{}
		for SpotInfo.ID == "" {
			json.Unmarshal(RequestArtistById(k, ATS), &SpotInfo)
		}
		SpotInfoTab = append(SpotInfoTab, SpotInfo)
	}
	Top3Temp := [len(ApiStruct.Top3Artists)]ApiTop3Artists{}
	for i := 0; i < len(ApiStruct.Top3Artists); i++ {
		SpotInfo := SpotifyArtistInfo{
			Followers: struct {
				Total int "json:\"total\""
			}{Total: 0},
		}
		for _, k := range SpotInfoTab {
			if k.Followers.Total > SpotInfo.Followers.Total {
				if !CheckArtistInTabArtists(Top3Temp, k) {
					SpotInfo = k
				}
			}
		}
		ArtistInfo := ApiArtistInfo{}
		id := strconv.Itoa(ApiStruct.AllArtistsId[SpotInfo.ID])
		json.Unmarshal(GetReadAll(LinkApi()["artists"]+"/"+id), &ArtistInfo)
		Top3Temp[i].Info = ArtistInfo
		Top3Temp[i].Spotify = SpotInfo
	}
	ApiStruct.Top3Artists = Top3Temp
}

func CheckArtistInTabArtists(Top3Temp [3]ApiTop3Artists, k SpotifyArtistInfo) bool {
	for _, l := range Top3Temp {
		if l.Spotify == k {
			return true
		}
	}
	return false
}
