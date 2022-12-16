package main

import (
	"fmt"
	"tracker/api-user"
)

func main() {
	fmt.Println(tracker.GetApiUrl())
	fmt.Println(tracker.GetArtistInfo(1))
	fmt.Println(tracker.GetDateInfo(1))
	fmt.Println(tracker.GetLocationInfo(1))
}
