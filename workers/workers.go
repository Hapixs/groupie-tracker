package workers

import (
	"api"
	"objects"
	"sync"
)

// Private variables of the package 'workers'
var mutex sync.Mutex

// Publics variables of the package 'workers'
var GroupList = make([]*objects.Group, 0)
var TrackList = make([]*objects.Track, 0)
var GenreList = make([]*objects.MusicGenre, 0)
var ArtistList = make([]*objects.Artist, 0)

var LocationByName = make(map[string]([]*objects.Location))
var TrackById = make(map[int]*objects.Track)

var GroupById = make(map[int]*objects.Group)
var GroupByGenre = map[*objects.MusicGenre]([]*objects.Group){}

var GenreById = make(map[int]*objects.MusicGenre)

func Init() {
	api.LoadApiDataFromFile()
	prepareGroups()
}
