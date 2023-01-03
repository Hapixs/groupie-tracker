package main

import (
	"api"
	"handlers"
	"net/http"
	"objects"
	"os"
	"strconv"
	"utils"
)

func main() {

	objects.WebServerConfig.InitConfig()

	println("Reading start arguments")
	objects.GameProcessArguments(os.Args[1:])

	utils.PrepareFolders()

	api.LoadGroups()
	go utils.UpdateAllGroupsPics()

	// Can be moved in a func like Run for some reason
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	handlers.InitHandlers()

	serverport, _, _ := objects.WebServerConfig.GetConfigItem(objects.ServerPort)
	http.ListenAndServe(":"+strconv.Itoa(serverport), nil)
}
