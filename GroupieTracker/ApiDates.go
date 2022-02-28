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

func ApiDates(Api *Api) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	DatesInd, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var DatesStruct Dates
	json.Unmarshal(DatesInd, &DatesStruct)
	Api.ApiDates = DatesStruct.Index
}
