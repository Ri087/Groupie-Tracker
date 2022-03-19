package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	Top3Artists            [3]ApiTop3Artist
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

type ApiTop3Artist struct {
	Info    ApiArtistInfo
	Spotify SpotifyArtistInfo
}

type SpotifyArtistInfo struct {
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	ID         string `json:"id"`
	Popularity int    `json:"popularity"`
}

type ApiArtistInfo struct {
	Id    int
	Image string
	Name  string
}

func GenerateTop3Artists(ApiStruct *ApiStructure, ATS *TokenSpotify) {
	// ApiStruct.AllIdArtists
	for ind, id := range ApiStruct.AllIdArtists {
		fmt.Println(ind, "", id)
	}
}
