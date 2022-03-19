package GroupieTracker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// TOKEN REQUEST

const (
	BASE_URL     = "https://api.spotify.com"
	ACCOUNTS_URL = "https://accounts.spotify.com/api/token"
	API_VERSION  = "v1"
)

type Spotify struct {
	clientID     string
	clientSecret string
	accessToken  string
}

type TokenSpotify struct {
	Access_token string
	Token_type   string
	Expires_in   int
}

func New(clientID, clientSecret string) Spotify {
	return initialize(clientID, clientSecret)
}

func initialize(clientID, clientSecret string) Spotify {
	spot := Spotify{clientID: clientID, clientSecret: clientSecret}
	return spot
}

func (spotify *Spotify) Authorize() *TokenSpotify {
	ATS := &TokenSpotify{}
	client := http.Client{}
	data := url.Values{}
	auth := fmt.Sprintf("Basic %s", spotify.getEncodedKeys())
	data.Set("grant_type", "client_credentials")
	req, _ := http.NewRequest("POST", ACCOUNTS_URL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", auth)
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &ATS)
	return ATS
}

func (spotify *Spotify) getEncodedKeys() string {
	data := fmt.Sprintf("%v:%v", spotify.clientID, spotify.clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return encoded
}

// END REQUEST TOKEN -----------------------------------------

func NameNoSpace(nameArtist string) string {
	var name string
	for _, i := range nameArtist {
		if i != ' ' {
			name += string(i)
		} else {
			name += "+"
		}
	}
	return name
}

type ToGetId struct {
	Artists struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	} `json:"artists"`
}

func RequestSearchId(name string, ATS *TokenSpotify) ToGetId {
	data := url.Values{}
	base_url := "https://api.spotify.com/v1/search?query=" + name + "&type=track%2Cartist%2Cartist&market=FR&limit=1"
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	req.Header.Set("Authorization", "Bearer "+ATS.Access_token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)

	var Id ToGetId
	json.Unmarshal(body, &Id)
	return Id
}

func GetEveryId(ApiStruct *ApiStructure, ATS *TokenSpotify) {
	var AllIdArtists = make(map[int]string)
	var AllArtistsId = make(map[string]int)
	for _, i := range ApiStruct.TabApiArtiste {
		name := NameNoSpace(i.Name)
		Id := RequestSearchId(name, ATS)
		for len(Id.Artists.Items) == 0 {
			Id = RequestSearchId(name, ATS)
		}
		AllIdArtists[i.Id] = Id.Artists.Items[0].ID
		AllArtistsId[Id.Artists.Items[0].ID] = i.Id
	}
	ApiStruct.AllIdArtists = AllIdArtists
	ApiStruct.AllArtistsId = AllArtistsId
}

func RequestArtistById(id string, ATS *TokenSpotify) []byte {
	data := url.Values{}
	base_url := "https://api.spotify.com/v1/artists/" + id
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	req.Header.Set("Authorization", "Bearer "+ATS.Access_token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func RequestByIdTopTrack(id string, ATS *TokenSpotify) SpotifyTopTrack {
	data := url.Values{}
	base_url := "https://api.spotify.com/v1/artists/" + id + "/top-tracks?market=FR"
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	req.Header.Set("Authorization", "Bearer "+ATS.Access_token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)

	SpotTrack := SpotifyTopTrack{}
	json.Unmarshal(body, &SpotTrack)
	return SpotTrack
}

type SpotifyPageArtiste struct {
	Name         string
	Followers    int
	Genres       []string
	ApiHref      string
	SpotifyHref  string
	Rank         int
	TrackHref    string
	TrackName    string
	BioPublished string
	ExtraitBio   string
	FullBio      string
}

type SpotifyStruct struct {
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Genres     []string `json:"genres"`
	Href       string   `json:"href"`
	Name       string   `json:"name"`
	Popularity int      `json:"popularity"`
}

type SpotifyTopTrack struct {
	Tracks []struct {
		ID   string
		Name string `json:"name"`
	} `json:"tracks"`
}

func PageArtistSpotify(AllId map[int]string, ApiArtists ArtistsApiPageArtiste, id int, ATS *TokenSpotify) SpotifyPageArtiste {
	ApiSpotify := SpotifyStruct{}
	json.Unmarshal(RequestArtistById(AllId[id], ATS), &ApiSpotify)
	SpotTrack := RequestByIdTopTrack(AllId[id], ATS)
	for len(SpotTrack.Tracks) == 0 {
		SpotTrack = RequestByIdTopTrack(AllId[id], ATS)
	}
	LastFmApi := LastfmRequest(NameNoSpace(ApiArtists.Artists.Name))

	Artist := SpotifyPageArtiste{
		Name:         ApiSpotify.Name,
		Followers:    ApiSpotify.Followers.Total,
		Genres:       ApiSpotify.Genres,
		ApiHref:      ApiSpotify.Href,
		Rank:         ApiSpotify.Popularity,
		SpotifyHref:  ApiSpotify.Href,
		TrackHref:    "https://open.spotify.com/embed/track/" + SpotTrack.Tracks[0].ID,
		TrackName:    SpotTrack.Tracks[0].Name,
		BioPublished: LastFmApi.Artist.Bio.Published,
		ExtraitBio:   LastFmApi.Artist.Bio.Summary,
		FullBio:      LastFmApi.Artist.Bio.Content,
	}
	return Artist
}

type Genres struct {
	Genres []string `json:"genres"`
}

func RequestGenresById(id string, ATS *TokenSpotify) Genres {
	data := url.Values{}
	base_url := "https://api.spotify.com/v1/artists/" + id
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	req.Header.Set("Authorization", "Bearer "+ATS.Access_token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)

	Info := Genres{}
	json.Unmarshal(body, &Info)
	return Info
}

func TabGenres(ApiStruct *ApiStructure, ATS *TokenSpotify) {
	GenresTab := []string{"All"}
	for _, i := range ApiStruct.AllIdArtists {
		Info := RequestGenresById(i, ATS)
		for _, k := range Info.Genres {
			if !CheckIfInTab(k, GenresTab) {
				GenresTab = append(GenresTab, k)
			}
		}
	}
	ApiStruct.Filtres.GenresTab = GenresTab
}
