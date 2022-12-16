package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	apiUrl = "https://groupietrackers.herokuapp.com/api"
)

type MainPageResponse struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

type ApiArtist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type ApiLocation struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type ApiDate struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type ApiRelation struct {
	Id             int                   `json:"id"`
	DatesLocations map[string]([]string) `json:"datesLocations"`
}

func GetApiUrl() []string {
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var mainPageResponse MainPageResponse
	err = json.Unmarshal(body, &mainPageResponse)
	if err != nil {
		log.Fatal(err)
	}
	return []string{mainPageResponse.Artists, mainPageResponse.Locations, mainPageResponse.Dates, mainPageResponse.Relation}
}

func GetAllArtist() []ApiArtist {
	url := GetApiUrl()[0]
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var artists []ApiArtist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Fatal(err)
	}
	return artists
}

func GetArtistInfo(id int) ApiArtist {
	url := GetApiUrl()[0] + "/" + strconv.Itoa(id)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var artist ApiArtist
	err = json.Unmarshal(body, &artist)
	if err != nil {
		log.Fatal(err)
	}
	return artist
}

func GetArtistByName(name string) ApiArtist {
	allArtist := GetAllArtist()
	for _, artist := range allArtist {
		if artist.Name == name {
			return artist
		}
	}
	return ApiArtist{}
}

func GetLocationInfo(id int) ApiLocation {
	url := GetApiUrl()[1] + "/" + strconv.Itoa(id)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var location ApiLocation
	err = json.Unmarshal(body, &location)
	if err != nil {
		log.Fatal(err)
	}
	return location
}

func GetDateInfo(id int) ApiDate {
	url := GetApiUrl()[2] + "/" + strconv.Itoa(id)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var date ApiDate
	err = json.Unmarshal(body, &date)
	if err != nil {
		log.Fatal(err)
	}
	return date
}

func GetRelationInfo(id int) ApiRelation {
	url := GetApiUrl()[3] + "/" + strconv.Itoa(id)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var relation ApiRelation
	err = json.Unmarshal(body, &relation)
	if err != nil {
		log.Fatal(err)
	}
	return relation
}

//Some random stuff here

type GoogleResponse struct {
	Search_metadata    []string
	Search_parameters  []string
	Search_information []string
	Images_results     []GoogleImage
}

type GoogleImage struct {
	Position        int
	Thumbnail       string
	Source          string
	Title           string
	Link            string
	Original        string
	Original_width  int
	Original_height int
	Is_product      bool
}

func GetArtistPictureLink(name string) string {
	name = strings.Replace(name, " ", "%20", -1)
	url := "https://serpapi.com/search.json?q=" + name + "&tbm=isch&api_key=2c1bc58028db937882d64c5c61e3b444aa159eacdda9340b385f33023ebe8a14"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	GResponse := GoogleResponse{}

	json.Unmarshal(body, &GResponse)

	if len(GResponse.Images_results) < 1 {
		return "https://www.google.com/url?sa=i&url=http%3A%2F%2Fpleasepretty.elob.fr%2Fjackie-chan-wtf-meme-face-70958233396%2F&psig=AOvVaw2XVDgb6TVEVbxo_yX_3v_q&ust=1671290961728000&source=images&cd=vfe&ved=0CBAQjRxqFwoTCMjLpZO6_vsCFQAAAAAdAAAAABAE"
	}
	return GResponse.Images_results[0].Thumbnail
}
