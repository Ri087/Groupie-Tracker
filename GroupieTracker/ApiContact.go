package GroupieTracker

import (
	"encoding/json"
	"fmt"
)

type ApiContacts struct {
	GroupieTracker []struct {
		ID       int    `json:"Id"`
		Name     string `json:"Name"`
		Age      int    `json:"Age"`
		Post     string `json:"Post"`
		Mail     string `json:"Mail"`
		URLImage string `json:"UrlImage"`
	} `json:"GroupieTracker"`
}

func ApiContactRequest() *ApiContacts {
	contact := ApiContacts{}
	url := "https://pacific-plains-95254.herokuapp.com/"
	json.Unmarshal(GetReadAll(url), &contact)
	fmt.Println(contact.GroupieTracker[0].Name)
	return &contact

}
