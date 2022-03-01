package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Relations struct {
	Index []Relation
}

func ApiRelations(Api *Api) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	RelationsInd, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var RelationsStruct Relations
	json.Unmarshal(RelationsInd, &RelationsStruct)
	Api.ApiRelations = RelationsStruct.Index
}
