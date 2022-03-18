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

type SpotifyTopTrack struct {
	Tracks []struct {
		Href string `json:"href"`
		Name string `json:"name"`
	} `json:"tracks"`
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

func GetAllIdArtistsJson(ApiStruct *ApiStructure) {
	save, err := ioutil.ReadFile("./GroupieTracker/Account/AllIdArtists.json")
	if err != nil {
		ioutil.WriteFile("./GroupieTracker/Account/AllIdArtists.json", []byte("{}"), 0644)
		save, _ = ioutil.ReadFile("./GroupieTracker/Account/AllIdArtists.json")
	}
	json.Unmarshal(save, &ApiStruct.AllIdArtists)
}

func SaveAllIdArtists(ApiStruct *ApiStructure) {
	AllIdArtists, _ := json.Marshal(ApiStruct.AllIdArtists)
	ioutil.WriteFile("./GroupieTracker/Account/AllIdArtists.json", AllIdArtists, 0644)
}

type ToGetId struct {
	Artists struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	} `json:"artists"`
}

func GetEveryId(ApiStruct *ApiStructure, ATS *TokenSpotify) {
	GetAllIdArtistsJson(ApiStruct)
	for _, i := range ApiStruct.TabApiArtiste {
		var Id ToGetId
		name := NameNoSpace(i.Name)
		body := Request(name, ATS)
		json.Unmarshal(body, &Id)
		ApiStruct.AllIdArtists[i.Id] = Id.Artists.Items[0].ID
	}
	SaveAllIdArtists(ApiStruct)
}

func PageArtistSpotify(AllId map[int]string, ApiArtists ArtistsApiPageArtiste, id int, ATS *TokenSpotify) SpotifyPageArtiste {
	ApiSpotify := SpotifyStruct{}
	SpotTrack := SpotifyTopTrack{}
	Artist := SpotifyPageArtiste{}
	name := NameNoSpace(ApiArtists.Artists.Name)
	body := RequestByIdArtist(id, AllId, ATS)
	json.Unmarshal(body, &ApiSpotify)
	LastFmApi := LastfmRequest(name)
	Artist.Name = ApiSpotify.Name
	Artist.Followers = ApiSpotify.Followers.Total
	Artist.Genres = ApiSpotify.Genres
	Artist.ApiHref = ApiSpotify.Href
	Artist.Rank = ApiSpotify.Popularity
	Artist.SpotifyHref = ApiSpotify.Href
	body = RequestByIdTopTrack(id, AllId, ATS)
	json.Unmarshal(body, &SpotTrack)

	Artist.TrackName = SpotTrack.Tracks[0].Name
	Artist.TrackHref = SpotTrack.Tracks[0].Href

	Artist.BioPublished = LastFmApi.Artist.Bio.Published
	Artist.ExtraitBio = LastFmApi.Artist.Bio.Summary
	Artist.FullBio = LastFmApi.Artist.Bio.Content
	return Artist
}
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
func Request(name string, ATS *TokenSpotify) []byte {
	data := url.Values{}
	client := http.Client{}
	base_url := "https://api.spotify.com/v1/search?query=" + name + "&type=track%2Cartist%2Cartist&market=FR&limit=1"
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	req.Header.Set("Authorization", "Bearer "+ATS.Access_token)
	req.Header.Set("Content-Type", "application/json")
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func RequestByIdArtist(id int, AllId map[int]string, ATS *TokenSpotify) []byte {
	data := url.Values{}
	client := http.Client{}
	base_url := "https://api.spotify.com/v1/artists/" + AllId[id]
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	req.Header.Set("Authorization", "Bearer "+ATS.Access_token)
	req.Header.Set("Content-Type", "application/json")
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func RequestByIdTopTrack(id int, AllId map[int]string, ATS *TokenSpotify) []byte {
	data := url.Values{}
	client := http.Client{}
	base_url := "https://api.spotify.com/v1/artists/" + AllId[id] + "/top-tracks?market=FR"
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	req.Header.Set("Authorization", "Bearer "+ATS.Access_token)
	req.Header.Set("Content-Type", "application/json")
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func FiltreArtsitSpotify(ApiStruct *ApiStructure, ATS *TokenSpotify, filters map[string][]string) {
	tempoTab := ApiStruct.TabApiFiltre
	ApiStruct.TabApiFiltre = []ApiArtiste{}
	for _, i := range tempoTab {
		ApiSpotify := SpotifyStruct{}
		name := NameNoSpace(i.Name)
		body := Request(name, ATS)
		json.Unmarshal(body, &ApiSpotify)

		AppendTabSpotify(i, filters, ApiSpotify, ApiStruct)

	}
}

func AppendTabSpotify(i ApiArtiste, filters map[string][]string, ApiSpotify SpotifyStruct, ApiStruct *ApiStructure) {
	for _, l := range filters["genres"] {
		for _, k := range ApiSpotify.Genres {
			if l == k {
				ApiStruct.TabApiFiltre = append(ApiStruct.TabApiFiltre, i)
				return
			}
		}
	}
}

type StructArtistTop3 struct {
	Followers int
}
