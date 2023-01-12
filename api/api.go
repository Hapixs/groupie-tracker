package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
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

type ApiRelation struct {
	Id             int                   `json:"id"`
	DatesLocations map[string]([]string) `json:"datesLocations"`
}

func getApiUrl() []string {
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

func getAllArtist() []ApiArtist {
	url := getApiUrl()[0]
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

func GetRelationInfo(id int) ApiRelation {
	url := getApiUrl()[3] + "/" + strconv.Itoa(id)
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

// redo the api for local cache

type Group struct {
	Id             int
	ImageLink      string
	Name           string
	Members        []Artist
	CreationYear   int
	FirstAlbumDate string
	DateLocations  []Date
}

type Artist struct {
	Name string
}

type Date struct {
	Locations string
	DateTime  string
}

var GroupMap = map[int](Group){}

var wg sync.WaitGroup

func LoadGroups() {
	println("Loading groups in cache for better performances")
	groups := getAllArtist()
	for _, v := range groups {
		go transformAndCacheGroup(v)
	}
	wg.Wait()
	println(strconv.Itoa(len(GroupMap)) + " groups have been loaded in cache!")
}

var mutex sync.Mutex

func transformAndCacheGroup(v ApiArtist) {
	defer wg.Done()
	wg.Add(1)
	g := Group{
		Id:             v.Id,
		ImageLink:      v.Image,
		Name:           v.Name,
		CreationYear:   v.CreationDate,
		FirstAlbumDate: v.FirstAlbum,
	}
	for _, m := range v.Members {
		g.Members = append(g.Members, Artist{m})
	}
	relations := GetRelationInfo(v.Id)
	for key, value := range relations.DatesLocations {
		for _, date := range value {
			g.DateLocations = append(g.DateLocations, Date{
				Locations: key,
				DateTime:  date,
			})
		}
	}
	mutex.Lock()
	GroupMap[v.Id] = g
	mutex.Unlock()
}

func GetCachedGroups() []Group {
	gs := []Group{}
	for _, value := range GroupMap {
		gs = append(gs, value)
	}
	return gs
}

func GetGroupFromName(name string) *Group {
	for _, v := range GroupMap {
		if v.Name == name {
			return &v
		}
	}
	return &Group{}
}

func GetGroupFromId(id int) Group {
	return GroupMap[id]
}

func EditGroupImageLink(id int, url string) {
	g := GetGroupFromId(id)
	g.ImageLink = url
	GroupMap[id] = g
}

func GetGroupListFiltredByName(filter string) []Group {
	sortedGroups := []Group{}
	filter = strings.ToUpper(filter)
	for _, k := range GroupMap {
		if strings.Contains(strings.ToUpper(k.Name), strings.ToUpper(filter)) {
			sortedGroups = append(sortedGroups, k)
		}
	}

	for i := 0; i < len(sortedGroups); i++ {
		for j := i + 1; j < len(sortedGroups); j++ {
			if strings.Index(strings.ToUpper(sortedGroups[i].Name), filter) > strings.Index(strings.ToUpper(sortedGroups[j].Name), filter) {
				sortedGroups[i], sortedGroups[j] = sortedGroups[j], sortedGroups[i]
			}
		}
	}

	return sortedGroups
}

func GetGroupListFiltredByLocation(filter string) []Group {
	sortedGroups := []Group{}
	filter = strings.ToUpper(filter)
	for _, k := range GroupMap {
		for _, date := range k.DateLocations {
			if strings.Contains(strings.ToUpper(date.Locations), strings.ToUpper(filter)) {
				sortedGroups = append(sortedGroups, k)
				break
			}
		}
	}

	return sortedGroups
}

func GetGroupListFiltredByDate(filter string) []Group {
	sortedGroups := []Group{}
	for _, v := range GroupMap {
		for _, date := range v.DateLocations {
			d := TransformDateToText(date.DateTime)
			if strings.Contains(strings.ToUpper(d), strings.ToUpper(filter)) {
				sortedGroups = append(sortedGroups, v)
				break
			}
		}
	}
	return sortedGroups
}

func IsNumeric(s string) bool {
	for _, c := range s {
		if !(c >= 48 && c <= 57) {
			return false
		}
	}
	return true
}

func TransformDateToText(dateTime string) string {
	d := strings.Split(dateTime, "-")
	day, _ := strconv.Atoi(d[0])
	month, _ := strconv.Atoi(d[1])
	year, _ := strconv.Atoi(d[2])

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location())

	str := date.Weekday().String() + "-" + date.Month().String() + "-" + strconv.Itoa(date.Year())
	return str
}

func GetGroupListFiltredByAll(filter string) []Group {
	sortedGroups := []Group{}
	filter = strings.ToUpper(filter)
	check := make(map[int](int))

	for _, k := range GroupMap {
		if strings.Contains(strings.ToUpper(k.Name), filter) {
			check[k.Id] = 1
			continue
		}
		for _, date := range k.DateLocations {
			if strings.Contains(strings.ToUpper(date.Locations), filter) {
				check[k.Id] = 1
				continue
			}
			d := TransformDateToText(date.DateTime)
			if strings.Contains(strings.ToUpper(d), filter) {
				check[k.Id] = 1
				continue
			}
		}
	}

	for k := range check {
		sortedGroups = append(sortedGroups, GetGroupFromId(k))
	}

	for i := 0; i < len(sortedGroups); i++ {
		for j := i + 1; j < len(sortedGroups); j++ {
			if strings.Contains(strings.ToUpper(sortedGroups[i].Name), filter) && strings.Contains(strings.ToUpper(sortedGroups[j].Name), filter) {
				if strings.Index(strings.ToUpper(sortedGroups[i].Name), filter) > strings.Index(strings.ToUpper(sortedGroups[j].Name), filter) {
					sortedGroups[i], sortedGroups[j] = sortedGroups[j], sortedGroups[i]
				}
			} else if !strings.Contains(strings.ToUpper(sortedGroups[i].Name), filter) {
				sortedGroups[i], sortedGroups[j] = sortedGroups[j], sortedGroups[i]
			}

		}
	}

	return sortedGroups
}
