package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Date struct {
	Locations string
	DateTime  string
}

var wg sync.WaitGroup
var mutex sync.Mutex
var CacheApi = map[string]string{}

func CallExternalApi[T any](url string, structure *T, sleepTime time.Duration) {
	response, err := http.Get(url)
	if err != nil {
		println("Error when calling " + url)
		return
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		println("Error when reading body of " + url)
		return
	}
	err = json.Unmarshal(content, structure)
	if err != nil {
		println("Error when pasing body of " + url)
		return
	}

	mutex.Lock()
	CacheApi[url] = string(content)
	mutex.Unlock()

	time.Sleep(time.Duration(sleepTime))
}

func CallCacheApi[T any](url string, structure *T) error {
	mutex.Lock()
	val, ok := CacheApi[url]
	mutex.Unlock()
	if ok {
		return json.Unmarshal([]byte(val), structure)
	}
	return errors.New("Unable to find in api cache " + url)
}

func GetFromApi[T any](url string, structure *T, update bool, sleepTime time.Duration) {
	err := CallCacheApi(url, structure)
	if err != nil || update {
		CallExternalApi(url, structure, sleepTime)
	}
}

func SaveApiCacheToFile() {
	mutex.Lock()
	save, err := json.Marshal(CacheApi)
	mutex.Unlock()
	if err != nil {
		println("Error: JSON error")
		return
	}

	file, err := os.Create("data.json")
	if err != nil {
		println("Error 2 with data")
		return
	}

	file.Write(save)
	file.Close()
}

func LoadApiDataFromFile() {
	content, err := os.ReadFile("data.json")
	if err == nil {
		json.Unmarshal(content, &CacheApi)
	} else {
		println("Seams like the first start of this web server.")
		println("Some operations may take several minutes.")
	}
}
