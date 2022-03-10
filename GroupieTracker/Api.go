package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ApiInit(Api *Api) {

	
	ApiTab := []string{"artists", "dates", "locations", "relation"}
	for _, i := range ApiTab {
		response, err := http.Get("https://groupietrackers.herokuapp.com/api/" + i)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		ApiIndex, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if i == "artists" {
			json.Unmarshal(ApiIndex, &Api.ApiArtist)
		} else if i == "dates" {
			var DatesStruct Dates
			json.Unmarshal(ApiIndex, &DatesStruct)
			Api.ApiDates = DatesStruct.Index
		} else if i == "locations" {
			var LocationsStruct Locations
			json.Unmarshal(ApiIndex, &LocationsStruct)
			Api.ApiLocations = LocationsStruct.Index
		} else {
			var RelationsStruct Relations
			json.Unmarshal(ApiIndex, &RelationsStruct)
			Api.ApiRelations = RelationsStruct.Index
		}

	}
}

type Dates struct {
	Index []Date
}

type Locations struct {
	Index []Location
}

type Relations struct {
	Index []Relation
}
