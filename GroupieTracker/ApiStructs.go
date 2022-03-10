package GroupieTracker

type Api struct {
	ApiArtist    []Artist
	ApiDates     []Date
	ApiLocations []Location
	ApiRelations []Relation
	ApiFiltre    []Artist
	Id           int
}

type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	TLocations   []Locations
}

type Date struct {
	Id    int
	Dates []string
}

type Location struct {
	Id        int
	Locations []string
}

type Relation struct {
	Id             int
	DatesLocations map[string][]string
}

type MapStruct struct {
	Features []struct {
		Center []float64 `json:"center"`
	} `json:"features"`
}
