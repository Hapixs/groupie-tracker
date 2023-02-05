package objects

type Track struct {
	GroupId     int
	GroupName   string     `json:"group_name"`
	Id          int        `json:"id"`
	ReleaseDate string     `json:"release_date"`
	Title       string     `json:"title_short"`
	Preview     string     `json:"preview"`
	Genre       MusicGenre `json:"genre"`
	Cover       string     `json:"Cover"`
}
