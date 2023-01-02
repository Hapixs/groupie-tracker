package main

import (
	"api"
	"handlers"
	"net/http"
	"objects"
	"os"
	"utils"
)

func main() {

	println("Reading start arguments")
	objects.GameProcessArguments(os.Args[1:])

	utils.PrepareFolders()

	api.LoadGroups()
	go utils.UpdateAllGroupsPics()

	// Can be moved in a func like Run for some reason
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	handlers.InitHandlers()

	println("WebServer ready to use !")
	http.ListenAndServe(port, nil)
}

// Maybee create a config structur to save the config with custom flags etc ?
const (
	port = ":8080"
)
