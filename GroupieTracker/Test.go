package GroupieTracker

import (
	"encoding/base64"
	"fmt"
	"net/http"
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

func Test(w http.ResponseWriter, r *http.Request) {
	clientID := "88fe57bfdc1f4fe18473613343bd419c"
	clientSecret := "88fe57bfdc1f4fe18473613343bd419c"
	Stify := Spotify{clientID: clientID, clientSecret: clientSecret}

	data := fmt.Sprintf("%v:%v", Stify.clientID, Stify.clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	auth := fmt.Sprintf("Basic %s", encoded)

	fmt.Println(auth)

	// request := gorequest.New()
	// request := http.

	// request.Post(ACCOUNTS_URL)
	// func (s *SuperAgent) Post(targetUrl string) *SuperAgent
	// s.ClearSuperAgent()
	// s.Method = POST
	// s.Url = targetUrl
	// s.Errors = nil

	// request.Set("Authorization", auth)
	// s.Header.Set(param, value)
	// func (h Header) Set(key, value string)

	// request.Send("grant_type=client_credentials")
}

// client_id := "88fe57bfdc1f4fe18473613343bd419c"
// response_type := "token"
// redirect_uri := "http://localhost:8080/callback"
// scope := "user-read-private user-read-email"
// state := "34fFs29kd09"
// request, err := http.NewRequest("GET", "https://accounts.spotify.com/authorize?", nil)
// if err != nil {
// 	log.Fatal(err)
// 	os.Exit(1)
// }
// q := request.URL.Query()
// q.Add("client_id", client_id)
// q.Add("response_type", response_type)
// q.Add("redirect_uri", redirect_uri)
// q.Add("scope", scope)
// q.Add("state", state)
// request.URL.RawQuery = q.Encode()
// http.Redirect(w, request, request.URL.String(), http.StatusFound)
