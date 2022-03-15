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

const (
	BASE_URL     = "https://api.spotify.com"
	ACCOUNTS_URL = "https://accounts.spotify.com/api/token"
	API_VERSION  = "v1"
)

func New(clientID, clientSecret string) Spotify {

	return initialize(clientID, clientSecret)
}
func initialize(clientID, clientSecret string) Spotify {
	spot := Spotify{clientID: clientID, clientSecret: clientSecret}
	return spot
}

func (spotify *Spotify) Authorize(w http.ResponseWriter, r *http.Request) {
	IDArtist := r.URL.Path[9:]
	// id, _ := strconv.Atoi(IDArtist)
	// Main.ApiStruct.SpecificApiPageArtiste = ApiArtistsPageArtiste(IDArtist)
	SpecificApiPageArtiste := ApiArtistsPageArtiste(IDArtist)
	nameArtsit := SpecificApiPageArtiste.Artists.Name
	client := http.Client{}
	auth := fmt.Sprintf("Basic %s", spotify.getEncodedKeys())
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, _ := http.NewRequest("POST", ACCOUNTS_URL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", auth)
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	ATS := &TokenSpotify{}
	json.Unmarshal(body, ATS)

	RequestSpotify(&Spotify{}, ATS, &SpotifyStruct{}, nameArtsit)
}
func (spotify *Spotify) getEncodedKeys() string {
	data := fmt.Sprintf("%v:%v", spotify.clientID, spotify.clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return encoded
}

type SpotifyStruct struct {
	Artists struct {
		Href  string `json:"href"`
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Total int `json:"total"`
			} `json:"followers"`
			Genres []string `json:"genres"`
			Href   string   `json:"href"`
			ID     string   `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name       string `json:"name"`
			Popularity int    `json:"popularity"`
			Type       string `json:"type"`
			URI        string `json:"uri"`
		} `json:"items"`
	} `json:"artists"`
}

func RequestSpotify(spotify *Spotify, TS *TokenSpotify, spotifyArtist *SpotifyStruct, nameArtist string) {
	var name string
	data := url.Values{}
	client := http.Client{}
	for _, i := range nameArtist {
		if i != ' ' {
			name += string(i)
		}
	}
	base_url := "https://api.spotify.com/v1/search?q=" + name + "&type=artist"
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	fmt.Println("Access Token = ", TS.Access_token)
	req.Header.Set("Authorization", "Bearer "+TS.Access_token)
	req.Header.Set("Content-Type", "application/json")
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, spotifyArtist)
	fmt.Println(spotifyArtist.Artists.Items[0].Name)
	fmt.Println(spotifyArtist.Artists.Items[0].Followers.Total)
	fmt.Println(spotifyArtist.Artists.Items[0].Genres)
	fmt.Println(spotifyArtist.Artists.Items[0].ID)
	fmt.Println(spotifyArtist.Artists.Items[0].Type)
}
