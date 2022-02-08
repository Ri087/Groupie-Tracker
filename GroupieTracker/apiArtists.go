package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      string
	CreationDate int
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

func ApiArtists() []Artist {
	var artists []Artist
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	json.Unmarshal(responseData, &artists)
	return artists
	// for _, x := range artists {
	// 	fmt.Println(x.Name + "||" + x.Image + "||" + x.Members)
	// }
}
