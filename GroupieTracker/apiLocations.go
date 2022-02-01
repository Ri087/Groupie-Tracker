package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Locations struct {
	Index []Location
}

type Location struct {
	Id        int
	Locations []string
	Dates     string
}

func ApiLocations() {
	var locAPi Locations
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	json.Unmarshal(responseData, &locAPi)
	fmt.Println(locAPi)
}
