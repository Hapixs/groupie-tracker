package workers

import (
	"api"
	"objects"
	"strconv"
	"time"
)

var fontMap = map[string]string{
	"Classique":        "Nicoone",
	"Dance":            "Orbitron",
	"Electro":          "Orbitron",
	"Films/Jeux vidéo": "Orbitron",
	"Rap/Hip Hop":      "Lacquer",
	"Pop":              "Reem_Kufi_Ink",
	"Reggae":           "Shadows_Into_Light",
	"Metal":            "Metal_Mania",
	"Rock":             "Rock_Salt",
	"Alternative":      "Unbounded"}

func buildGenderFromDeezerId(deezerGenreId int) *objects.MusicGenre {
	type deezerGenre struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Picture   string `json:"picture_medium"`
		PictureXl string `json:"picture_xl"`
	}
	var deezerGenreRequest deezerGenre

	api.GetFromApi(
		"https://api.deezer.com/genre/"+strconv.Itoa(deezerGenreId),
		&deezerGenreRequest,
		false,
		time.Second/9,
		nil)

	musicGenre := new(objects.MusicGenre)
	musicGenre.Id = deezerGenreRequest.Id
	musicGenre.Name = deezerGenreRequest.Name
	musicGenre.PictureXl = deezerGenreRequest.PictureXl
	musicGenre.Picture = deezerGenreRequest.Picture
	musicGenre.FontName = fontMap[musicGenre.Name]

	mutex.Lock()
	GenreList = append(GenreList, musicGenre)
	GenreById[musicGenre.Id] = musicGenre
	mutex.Unlock()

	return musicGenre
}
