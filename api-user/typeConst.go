package tracker

type MainPageResponse struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int   `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string `json:"locations"`
	ConcertDates string `json:"concertDates"`
	Relation     string `json:"relation"`
}

type Location struct {
	Id        int    `json:"id"`
	Locations []string `json:"locations"`
	Dates    string `json:"dates"`
}

type Date struct {
	Id        int    `json:"id"`
	Dates	[]string `json:"dates"`
}

type Relation struct {
	Id        int    `json:"id"`
	DatesLocations []string `json:"datesLocations"`
}
