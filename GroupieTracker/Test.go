package GroupieTracker

import (
	"net/http"
)

type Spotify struct {
	clientID     string
	clientSecret string
	//	redirectURI        string
	accessToken string
}

type TestStruct struct {
	data interface{}
}

func Test(w http.ResponseWriter, r *http.Request) {
	// // BASE_URL := "https://api.spotify.com"
	// ACCOUNTS_URL := "https://accounts.spotify.com/api/token"
	// // API_VERSION := "v1"
	// clientID := "88fe57bfdc1f4fe18473613343bd419c"
	// clientSecret := "88fe57bfdc1f4fe18473613343bd419c"
	// Stify := Spotify{clientID: clientID, clientSecret: clientSecret}

	// data := fmt.Sprintf("%v:%v", Stify.clientID, Stify.clientSecret)
	// encoded := base64.StdEncoding.EncodeToString([]byte(data))

	// auth := fmt.Sprintf("Basic %s", encoded)

	// r.Method = "POST"
	// r.Header.Set("Authorization", auth)
	// r.Url = ACCOUNTS_URL

	// j := TestStruct{}
	// err := json.Unmarshal(body, j)

	// if err != nil {
	// 	fmt.Println("[Authorize] Error parsing Json!")
	// 	errs := []error{err}
	// 	fmt.Println(errs)
	// 	os.Exit(1)
	// }
}
