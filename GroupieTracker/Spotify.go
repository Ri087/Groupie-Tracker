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
type SpotifyStruct struct {
	Artists struct {
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Total int `json:"total"`
			} `json:"followers"`
			Genres     []string `json:"genres"`
			Href       string   `json:"href"`
			ID         string   `json:"id"`
			Name       string   `json:"name"`
			Popularity int      `json:"popularity"`
			URI        string   `json:"uri"`
		} `json:"items"`
	} `json:"artists"`
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

func (spotify *Spotify) Authorize() TokenSpotify {
	ATS := TokenSpotify{}
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

type SpotifyPageArtiste struct {
	Name        string
	Followers   int
	Genres      []string
	Id          string
	ApiHref     string
	SpotifyHref string
	Rank        int
}

func PageArtistSpotify(ID string, nameArtist string, ATS *TokenSpotify) *SpotifyPageArtiste {
	ApiSpotify := SpotifyStruct{}
	Artist := &SpotifyPageArtiste{}
	name := NameNoSpace(nameArtist)
	body := Request(name, ATS)
	json.Unmarshal(body, &ApiSpotify)
	Artist.Name = ApiSpotify.Artists.Items[0].Name
	Artist.Followers = ApiSpotify.Artists.Items[0].Followers.Total
	Artist.Genres = ApiSpotify.Artists.Items[0].Genres
	Artist.Name = ApiSpotify.Artists.Items[0].Name
	Artist.Id = ApiSpotify.Artists.Items[0].ID
	Artist.ApiHref = ApiSpotify.Artists.Items[0].Href
	Artist.Rank = ApiSpotify.Artists.Items[0].Popularity
	Artist.SpotifyHref = "https://open.spotify.com/artist/" + Artist.Id
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
	base_url := "https://api.spotify.com/v1/search?q=" + name + "&type=artist"
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
		for _, k := range ApiSpotify.Artists.Items[0].Genres {
			if l == k {
				ApiStruct.TabApiFiltre = append(ApiStruct.TabApiFiltre, i)
				return
			}
		}
	}
}
