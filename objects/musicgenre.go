package objects

type MusicGenre struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Picture   string `json:"picture_medium"`
	PictureXl string `json:"picture_xl"`
	FontName  string `json:"font_name"`
}
