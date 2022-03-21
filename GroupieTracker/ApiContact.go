package GroupieTracker

import (
	"encoding/json"
)

type ApiContacts struct {
	GroupieTracker []struct {
		ID          int      `json:"Id"`
		Name        string   `json:"Name"`
		Age         int      `json:"Age"`
		Post        string   `json:"Post"`
		Mail        string   `json:"Mail"`
		URLImage    string   `json:"UrlImage"`
		Cities      string   `json:"Cities"`
		Competences []string `json:"Competences"`
	} `json:"GroupieTracker"`
}

func ApiContactRequest() *ApiContacts {
	contact := ApiContacts{}
	url := "https://cezgindustries-api-contact.herokuapp.com/"
	json.Unmarshal(GetReadAll(url), &contact)
	return &contact

}
