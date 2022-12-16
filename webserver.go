package main

import (
	"handlers"
	"net/http"
	"utils"
)

func main() {

	utils.PrepareFolders()

	println("Updating all assets...")

	go utils.UpdateAllAlbumPics()
	go utils.UpdateAllArtistsPics()

	// Can be moved in a func like Run for some reason
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	handlers.InitHandlers()

	http.ListenAndServe(port, nil)
}

// Maybee create a config structur to save the config with custom flags etc ?
const (
	port = ":8080"
)
