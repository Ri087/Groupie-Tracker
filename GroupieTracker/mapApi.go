package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type MapStruct struct {
	Features []struct {
		Center []float64 `json:"center"`
	} `json:"features"`
}

func Mapapi(ApiStruct *ApiStructure, id int) [][]float64 {
	var geolocationCities [][]float64
	for _, i := range ApiStruct.TabApiArtisteLocations[id].Locations {
		cities := strings.Split(i, "-")[0]
		url := "https://api.mapbox.com/geocoding/v5/mapbox.places/" + cities + ".json?types=place%2Cpostcode%2Caddress&access_token=pk.eyJ1IjoiYm9udGFheiIsImEiOiJjbDBnajQwazYwMGFqM2VxdnJkdDdpZmgxIn0.Cjpn9M5HiKAvVynLZbPlaQ"
		response, err := http.Get(url)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var MapStruct MapStruct
		json.Unmarshal(responseData, &MapStruct)
		geolocationCities = append(geolocationCities, MapStruct.Features[0].Center)
	}
	return geolocationCities
}
