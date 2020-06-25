package main

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const testDataDir = "example/assets-src"

func TestZip(t *testing.T) {
	cases := []struct {
		name           string
		ignoreModTimes bool
		useFileMode    os.FileMode
	}{
		{
			name:           "zip dir",
			ignoreModTimes: false,
		},
		{
			name:           "zip dir and ignore mod times",
			ignoreModTimes: true,
		},
		{
			name:           "zip dir and use special file mode",
			ignoreModTimes: false,
			useFileMode:    os.ModePerm,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			var zippedDir bytes.Buffer
			err := Zip(&zippedDir, testDataDir,
				WithIgonreModTimes(c.ignoreModTimes),
				WithFileMode(c.useFileMode),
			)
			if err != nil {
				t.Fatalf("Zip failed: %s", err)
			}

			// check archive
			zipReader, err := zip.NewReader(bytes.NewReader(zippedDir.Bytes()), int64(len(zippedDir.Bytes())))
			if err != nil {
				panic(err)
			}
			for _, zippedFile := range zipReader.File {
				// check original file exists
				origFileStat, err := os.Stat(filepath.Join(testDataDir, zippedFile.Name))
				if err != nil {
					t.Fatalf("zipped file %s original file stat error: %s", zippedFile.Name, err)
				}
				// check size
				if int64(zippedFile.UncompressedSize64) != origFileStat.Size() {
					t.Fatalf("zipped file %s size: %d,  original file size: %d", zippedFile.Name,
						zippedFile.UncompressedSize64, origFileStat.Size())
				}
				//check modtime
				wantModTime := origFileStat.ModTime()
				if c.ignoreModTimes {
					wantModTime = defaultModTime
				}
				if zippedFile.Modified.Unix() != wantModTime.Unix() {
					t.Fatalf("zipped file %s modtime: %s,  original file modtime: %s", zippedFile.Name,
						zippedFile.Modified, wantModTime)
				}
				//check file mode
				wantMod := origFileStat.Mode()
				if c.useFileMode != 0 {
					wantMod = os.ModePerm
				}

				if zippedFile.Mode() != wantMod {
					t.Fatalf("zipped file %s mode: %s,  original file mode: %s", zippedFile.Name,
						zippedFile.Mode(), wantMod)
				}
			}

		})
	}
}
