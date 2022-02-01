package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Dates struct {
	Index []Date
}

type Date struct {
	Id    int
	Dates []string
}

func ApiDates() {
	var dateAPi Dates
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	json.Unmarshal(responseData, &dateAPi)
	fmt.Println(dateAPi)
}
