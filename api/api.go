package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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
	Name      string
	ImageLink string
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
		a := LoadArtistWithImage(v.Name, m)
		g.Members = append(g.Members, a)
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
	filters := strings.Split(filter, ",")
	check := make(map[int](int))

	for _, f := range filters {
		for _, k := range GroupMap {
			if strings.Contains(strings.ToUpper(k.Name), f) {
				check[k.Id] = 1
				continue
			}
			for _, date := range k.DateLocations {
				if strings.Contains(strings.ToUpper(date.Locations), f) {
					check[k.Id] = 1
					continue
				}
				d := TransformDateToText(date.DateTime)
				if strings.Contains(strings.ToUpper(d), f) {
					check[k.Id] = 1
					continue
				}
			}
			for _, m := range k.Members {
				if strings.Contains(strings.ToUpper(m.Name), f) {
					check[k.Id] = 1
				}
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

// Thx wikipedia :D

type WikiRequest struct {
	Query WikiQuery `json:"query"`
}

type WikiQuery struct {
	Page map[int](WikiData) `json:"pages"`
}

type WikiData struct {
	Thumbnail WikiThumbnail `json:"thumbnail"`
}

type WikiThumbnail struct {
	Source string `json:"source"`
}

func GetWikipediaImage(artist string) WikiRequest {
	artist = RemoveAccents(artist)
	artist = strings.Join(strings.Split(artist, " "), "%20")
	rand.Seed(time.Now().UnixMilli())
	time.Sleep(time.Duration(rand.Intn(150)))
	url := "https://en.wikipedia.org/w/api.php?action=query&titles=" + artist + "&prop=pageimages&format=json&pithumbsize=100"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var request WikiRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		println("Error when parsing wikipedia api response for" + artist)
		return WikiRequest{}
	}
	return request
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func LoadArtistWithImage(group, m string) Artist {
	artist := Artist{m, ""}
	request := GetWikipediaImage(m)

	for _, k := range request.Query.Page {
		if k.Thumbnail.Source == "" {
			if group == m {
				artist.ImageLink = "https://cdn-icons-png.flaticon.com/512/32/32297.png"
			} else {
				artist = LoadArtistWithImage(group, group)
			}
		} else {
			artist.ImageLink = k.Thumbnail.Source
		}
	}
	return artist
}
