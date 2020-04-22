package main

import (
	"net/http"

	"github.com/gruzovator/zipper/example/assets"
)

func main() {
	http.Handle("/", http.FileServer(assets.NewZippedFilesHttpFS()))

	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(err)
	}
}
