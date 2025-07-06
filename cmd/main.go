package main

import "github.com/gilwong00/file-streamer/internal/server"

func main() {
	// TODO: will need to add logging and maybe a better implementation
	if err := server.StartServer(); err != nil {
		panic(err)
	}
}
