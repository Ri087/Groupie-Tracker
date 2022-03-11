package GroupieTracker

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Spotify struct {
	clientID     string
	clientSecret string
	//	redirectURI        string
	accessToken string
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
	bosy, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(bosy))

}
func (spotify *Spotify) getEncodedKeys() string {

	data := fmt.Sprintf("%v:%v", spotify.clientID, spotify.clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	return encoded
}
func (spotify *Spotify) createTargetURL(endpoint string) string {
	result := fmt.Sprintf("%s/%s/%s", BASE_URL, API_VERSION, endpoint)
	return result
}
