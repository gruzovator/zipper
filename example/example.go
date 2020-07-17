package main

import (
	"net/http"

	"github.com/gruzovator/zipper/example/assets"
)

//go:generate zipper -src assets-src -dest assets/assets.go -pkg assets -exclude *.bin

func main() {
	http.Handle("/", http.FileServer(assets.NewZippedFilesFS()))

	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(err)
	}
}
