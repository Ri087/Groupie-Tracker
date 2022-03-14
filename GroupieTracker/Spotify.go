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

func (spotify *Spotify) Authorize() {

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
	// fmt.Println(ATS.Access_token)
	// fmt.Println(ATS.Expires_in)
	// fmt.Println(ATS.Token_type)
	RequestSpotify(&Spotify{}, ATS, &Spo{})

}
func (spotify *Spotify) getEncodedKeys() string {
	data := fmt.Sprintf("%v:%v", spotify.clientID, spotify.clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return encoded
}

// func (spotify *Spotify) createTargetURL(endpoint string) string {
// 	result := fmt.Sprintf("%s/%s/%s", BASE_URL, API_VERSION, endpoint)
// 	fmt.Println(spotify)
// 	return result
// }
type Spo struct {
	Artists struct {
		Href  string `json:"href"`
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Href  interface{} `json:"href"`
				Total int         `json:"total"`
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
		Limit    int         `json:"limit"`
		Next     string      `json:"next"`
		Offset   int         `json:"offset"`
		Previous interface{} `json:"previous"`
		Total    int         `json:"total"`
	} `json:"artists"`
}

func RequestSpotify(spotify *Spotify, TS *TokenSpotify, test *Spo) {
	data := url.Values{}
	client := http.Client{}
	base_url := "https://api.spotify.com/v1/search?q=" + "Queen" + "&type=artist"
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	fmt.Println("Access Token = ", TS.Access_token)
	req.Header.Set("Authorization", "Bearer "+TS.Access_token)
	req.Header.Set("Content-Type", "application/json")
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, test)
	fmt.Println(test)

}
