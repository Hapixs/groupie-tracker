package workers

import (
	"api"
	"objects"
	"sync"
)

// Private variables of the package 'workers'
var waitgroup sync.WaitGroup
var mutex sync.Mutex

// Publics variables of the package 'workers'
var GroupMap = map[int](objects.Group){}
var GroupByGenreMap = map[api.DeezerGenre]([]objects.Group){}
var trackById map[int]objects.Track = make(map[int]objects.Track)
var Locations = map[string]([]api.Date){}
var artistsList = []objects.Artist{}
