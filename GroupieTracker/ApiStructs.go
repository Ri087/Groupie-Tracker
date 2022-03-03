package GroupieTracker

type Api struct {
	ApiArtist    []Artist
	ApiDates     []Date
	ApiLocations []Location
	ApiRelations []Relation
	Id           int
}

type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
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
