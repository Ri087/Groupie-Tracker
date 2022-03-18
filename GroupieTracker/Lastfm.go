package GroupieTracker

import (
	"encoding/json"
)

type LastFm struct {
	Artist struct {
		Name       string `json:"name"`
		Mbid       string `json:"mbid"`
		URL        string `json:"url"`
		Streamable string `json:"streamable"`
		Ontour     string `json:"ontour"`
		Stats      struct {
			Listeners string `json:"listeners"`
			Playcount string `json:"playcount"`
		} `json:"stats"`
		Similar struct {
			Artist []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"artist"`
		} `json:"similar"`
		Bio struct {
			Links struct {
				Link struct {
					Text string `json:"#text"`
					Rel  string `json:"rel"`
					Href string `json:"href"`
				} `json:"link"`
			} `json:"links"`
			Published string `json:"published"`
			Summary   string `json:"summary"`
			Content   string `json:"content"`
		} `json:"bio"`
	} `json:"artist"`
}

func LastfmRequest(name string) *LastFm {
	ApiLastFm := LastFm{}
	apiKey := "dbdb51d293e6fdfd0ef702b6dd98a64e"
	url := "http://ws.audioscrobbler.com/2.0/?method=artist.getinfo&artist=" + name + "&api_key=" + apiKey + "&format=json"
	json.Unmarshal(GetReadAll(url), &ApiLastFm)
	return &ApiLastFm
}
