package objects

type Track struct {
	GroupId     int
	Id          int    `json:"id"`
	ReleaseDate string `json:"release_date"`
	Title       string `json:"title_short"`
	Preview     string `json:"preview"`
	Album       struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Cover string `json:"cover_medium"`
		Genre MusicGenre
	} `json:"album"`
}
