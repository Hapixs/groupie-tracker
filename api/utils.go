package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

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

func GetAllArtist() []Artist {
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
	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Fatal(err)
	}
	return artists
}

func GetArtistInfo(id int) Artist {
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
	var artist Artist
	err = json.Unmarshal(body, &artist)
	if err != nil {
		log.Fatal(err)
	}
	return artist
}

func GetArtistByName(name string) Artist {
	allArtist := GetAllArtist()
	for _, artist := range allArtist {
		if artist.Name == name {
			return artist
		}
	}
	return Artist{}
}

func GetLocationInfo(id int) Location {
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
	var location Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		log.Fatal(err)
	}
	return location
}

func GetDateInfo(id int) Date {
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
	var date Date
	err = json.Unmarshal(body, &date)
	if err != nil {
		log.Fatal(err)
	}
	return date
}

func GetRelationInfo(id int) Relation {
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
	var relation Relation
	err = json.Unmarshal(body, &relation)
	if err != nil {
		log.Fatal(err)
	}
	return relation
}
