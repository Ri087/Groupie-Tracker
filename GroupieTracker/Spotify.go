package GroupieTracker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	BASE_URL     = "https://api.spotify.com"
	ACCOUNTS_URL = "https://accounts.spotify.com/api/token"
	API_VERSION  = "v1"
)

type Spotify struct {
	clientID     string
	clientSecret string
	accessToken  string
}

type TokenSpotify struct {
	Access_token string
	Token_type   string
	Expires_in   int
}

type SpotifyPageArtiste struct {
	Name         string
	Followers    int
	Genres       []string
	Id           string
	ApiHref      string
	SpotifyHref  string
	Rank         int
	TrackHref    string
	TrackName    string
	TrackId      string
	AlbumHref    string
	AlbumName    string
	AlbumRealase string
	AlbumNbTrack int
	AlbumId      string
	BioPublished string
	ExtraitBio   string
	FullBio      string
}
type SpotifyStruct struct {
	Artists struct {
		Items []struct {
			Followers struct {
				Total int `json:"total"`
			} `json:"followers"`
			Genres     []string `json:"genres"`
			Href       string   `json:"href"`
			ID         string   `json:"id"`
			Name       string   `json:"name"`
			Popularity int      `json:"popularity"`
			Type       string   `json:"type"`
			URI        string   `json:"uri"`
		} `json:"items"`
	} `json:"artists"`
	Tracks struct {
		Items []struct {
			Album struct {
				AlbumType    string `json:"album_type"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href                 string `json:"href"`
				ID                   string `json:"id"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				TotalTracks          int    `json:"total_tracks"`
				Type                 string `json:"type"`
				URI                  string `json:"uri"`
			} `json:"album"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Name string `json:"name"`
			URI  string `json:"uri"`
		} `json:"items"`
	} `json:"tracks"`
}

func New(clientID, clientSecret string) Spotify {

	return initialize(clientID, clientSecret)
}
func initialize(clientID, clientSecret string) Spotify {
	spot := Spotify{clientID: clientID, clientSecret: clientSecret}
	return spot
}

func (spotify *Spotify) Authorize() *TokenSpotify {
	ATS := &TokenSpotify{}
	client := http.Client{}
	data := url.Values{}
	auth := fmt.Sprintf("Basic %s", spotify.getEncodedKeys())
	data.Set("grant_type", "client_credentials")
	req, _ := http.NewRequest("POST", ACCOUNTS_URL, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", auth)
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &ATS)
	return ATS
}

func (spotify *Spotify) getEncodedKeys() string {
	data := fmt.Sprintf("%v:%v", spotify.clientID, spotify.clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return encoded
}

func PageArtistSpotify(ID string, nameArtist string, ATS *TokenSpotify) *SpotifyPageArtiste {
	ApiSpotify := SpotifyStruct{}
	Artist := &SpotifyPageArtiste{}
	name := NameNoSpace(nameArtist)
	body := Request(name, ATS)
	json.Unmarshal(body, &ApiSpotify)
	LastFmApi := LastfmRequest(name)
	Artist.Name = ApiSpotify.Artists.Items[0].Name
	Artist.Followers = ApiSpotify.Artists.Items[0].Followers.Total
	Artist.Genres = ApiSpotify.Artists.Items[0].Genres
	Artist.Id = ApiSpotify.Artists.Items[0].ID
	Artist.ApiHref = ApiSpotify.Artists.Items[0].Href
	Artist.Rank = ApiSpotify.Artists.Items[0].Popularity
	Artist.SpotifyHref = "https://open.spotify.com/artist/" + Artist.Id
	Artist.TrackName = ApiSpotify.Tracks.Items[0].Name
	Artist.TrackId = ApiSpotify.Tracks.Items[0].ID
	Artist.TrackHref = "https://open.spotify.com/embed/track/" + Artist.TrackId + "?utm_source=generator&theme=0"
	Artist.AlbumName = ApiSpotify.Tracks.Items[0].Album.Name
	Artist.AlbumId = ApiSpotify.Tracks.Items[0].Album.ID
	Artist.AlbumHref = "https://open.spotify.com/embed/album/" + Artist.AlbumId + "?utm_source=generator&theme=0"
	Artist.AlbumRealase = ApiSpotify.Tracks.Items[0].Album.ReleaseDate
	Artist.AlbumNbTrack = ApiSpotify.Tracks.Items[0].Album.TotalTracks
	Artist.BioPublished = LastFmApi.Artist.Bio.Published
	Artist.ExtraitBio = LastFmApi.Artist.Bio.Summary
	Artist.FullBio = LastFmApi.Artist.Bio.Content
	return Artist
}
func NameNoSpace(nameArtist string) string {
	var name string
	for _, i := range nameArtist {
		if i != ' ' {
			name += string(i)
		} else {
			name += "+"
		}
	}
	return name
}
func Request(name string, ATS *TokenSpotify) []byte {
	data := url.Values{}
	client := http.Client{}
	base_url := "https://api.spotify.com/v1/search?query=" + name + "&type=track%2Cartist%2Cartist&market=FR&limit=1"
	req, _ := http.NewRequest("GET", base_url, strings.NewReader(data.Encode()))
	req.Header.Set("Authorization", "Bearer "+ATS.Access_token)
	req.Header.Set("Content-Type", "application/json")
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func FiltreArtsitSpotify(ApiStruct *ApiStructure, ATS *TokenSpotify, filters map[string][]string) {
	tempoTab := ApiStruct.TabApiFiltre
	ApiStruct.TabApiFiltre = []ApiAccueil{}
	for _, i := range tempoTab {
		ApiSpotify := SpotifyStruct{}
		name := NameNoSpace(i.Name)
		body := Request(name, ATS)
		json.Unmarshal(body, &ApiSpotify)

		AppendTabSpotify(i, filters, ApiSpotify, ApiStruct)

	}
}

func AppendTabSpotify(i ApiAccueil, filters map[string][]string, ApiSpotify SpotifyStruct, ApiStruct *ApiStructure) {
	for _, l := range filters["genres"] {
		for _, k := range ApiSpotify.Artists.Items[0].Genres {
			if l == k {
				ApiStruct.TabApiFiltre = append(ApiStruct.TabApiFiltre, i)
				return
			}
		}
	}
}
