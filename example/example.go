package main

import (
	"net/http"

	"github.com/gruzovator/zipper/example/assets"
)

// next command packs assets-src excluding bin dir (bin/**) and all txt files (*.txt)
// Include/exclude patterns format: https://github.com/gobwas/glob
//go:generate zipper -src assets-src -dest assets/assets.go -pkg assets -exclude bin/** -exclude **.txt

func main() {
	http.Handle("/", http.FileServer(assets.NewZippedFilesFS()))

	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(err)
	}
}
