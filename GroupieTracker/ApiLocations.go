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

func ApiLocations(Api *Api) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	LocationsInd, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var LocationsStruct Locations
	json.Unmarshal(LocationsInd, &LocationsStruct)
	Api.ApiLocations = LocationsStruct.Index
}
