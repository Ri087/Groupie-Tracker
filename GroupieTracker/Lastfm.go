package GroupieTracker

import (
	"encoding/json"
)

type LastFm struct {
	Artist struct {
		Bio struct {
			Published string `json:"published"`
			Summary   string `json:"summary"`
			Content   string `json:"content"`
		} `json:"bio"`
	} `json:"artist"`
}

func LastfmRequest(name string) LastFm {
	apiKey := "dbdb51d293e6fdfd0ef702b6dd98a64e"
	url := "http://ws.audioscrobbler.com/2.0/?method=artist.getinfo&artist=" + name + "&api_key=" + apiKey + "&format=json"

	ApiLastFm := LastFm{}
	json.Unmarshal(GetReadAll(url), &ApiLastFm)
	return ApiLastFm
}
