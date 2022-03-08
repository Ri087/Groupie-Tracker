package GroupieTracker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Mapapi(popo *Api) {
	url := "https://api.mapbox.com/geocoding/v5/mapbox.places/Paris.json?types=place%2Cpostcode%2Caddress&access_token=pk.eyJ1IjoiYm9udGFheiIsImEiOiJjbDBnajQwazYwMGFqM2VxdnJkdDdpZmgxIn0.Cjpn9M5HiKAvVynLZbPlaQ"
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))

}
