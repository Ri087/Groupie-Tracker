package GroupieTracker

import (
	"log"
	"net/http"
	"os"
)

func Test(w http.ResponseWriter, r *http.Request) {
	client_id := "88fe57bfdc1f4fe18473613343bd419c"
	response_type := "token"
	redirect_uri := "http://localhost:8080/callback"
	scope := "user-read-private user-read-email"
	state := "34fFs29kd09"
	request, err := http.NewRequest("GET", "https://accounts.spotify.com/authorize?", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	q := request.URL.Query()
	q.Add("client_id", client_id)
	q.Add("response_type", response_type)
	q.Add("redirect_uri", redirect_uri)
	q.Add("scope", scope)
	q.Add("state", state)
	request.URL.RawQuery = q.Encode()
	http.Redirect(w, request, request.URL.String(), http.StatusFound)
}
