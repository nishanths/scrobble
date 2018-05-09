package main

import (
	"github.com/nishanths/scrobble/server"
	"google.golang.org/appengine"
)

func main() {
	server.RegisterHandlers()
	appengine.Main()
}
