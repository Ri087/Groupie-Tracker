package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func TraitementRealation() {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	json.Unmarshal(body, &All.Relation)
}
