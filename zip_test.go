package main

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testDataDir = "example/assets-src"

func TestZip_WithoutOptions(t *testing.T) {
	var data bytes.Buffer
	if err := Zip(&data, testDataDir); err != nil {
		t.Fatal(err)
	}

	files := testDataDirFiles()
	zFiles := zippedFiles(data.Bytes())

	if len(files) != len(zFiles) {
		t.Fatal("zipped files number is not equal to test files number")
	}

	for filePath, info := range files {
		zInfo, ok := zFiles[filePath]
		if !ok {
			t.Fatalf("no %q file in zipped data", filePath)
		}
		if info.Size() != zInfo.Size() {
			t.Fatalf("%q file size: orig: %d, zipped: %d", filePath, info.Size(), zInfo.Size())
		}
	}
}

func TestZip_WithModTime(t *testing.T) {
	now := time.Now().Round(time.Second)

	var data bytes.Buffer
	if err := Zip(&data, testDataDir, WithModTime(now)); err != nil {
		t.Fatal(err)
	}

	files := testDataDirFiles()
	zFiles := zippedFiles(data.Bytes())

	if len(files) != len(zFiles) {
		t.Fatal("zipped files number is not equal to test files number")
	}

	for filePath, info := range files {
		zInfo, ok := zFiles[filePath]
		if !ok {
			t.Fatalf("no %q file in zipped data", filePath)
		}
		if info.Size() != zInfo.Size() {
			t.Fatalf("%q file size: orig: %d, zipped: %d", filePath, info.Size(), zInfo.Size())
		}
		if !now.Equal(zInfo.ModTime()) {
			t.Fatalf("%q mod time: want: %s, zipped: %s", filePath, now, zInfo.ModTime())
		}
	}
}

func TestZip_WithFileMod(t *testing.T) {
	fileMode := os.FileMode(0644)

	var data bytes.Buffer
	if err := Zip(&data, testDataDir, WithFileMode(fileMode)); err != nil {
		t.Fatal(err)
	}

	files := testDataDirFiles()
	zFiles := zippedFiles(data.Bytes())

	if len(files) != len(zFiles) {
		t.Fatal("zipped files number is not equal to test files number")
	}

	for filePath, info := range files {
		zInfo, ok := zFiles[filePath]
		if !ok {
			t.Fatalf("no %q file in zipped data", filePath)
		}
		if info.Size() != zInfo.Size() {
			t.Fatalf("%q file size: orig: %d, zipped: %d", filePath, info.Size(), zInfo.Size())
		}
		if fileMode != zInfo.Mode() {
			t.Fatalf("%q file size: orig: 0%o, zipped: 0%o", filePath, fileMode, zInfo.Mode())
		}
	}
}

func TestZip_WithIncludeOnlyPNG(t *testing.T) {
	var data bytes.Buffer
	if err := Zip(&data, testDataDir, WithIncludePatterns("**.png")); err != nil {
		t.Fatal(err)
	}

	files := testDataDirFiles()
	zFiles := zippedFiles(data.Bytes())

	for filePath := range files {
		ext := filepath.Ext(filePath)
		_, isInZip := zFiles[filePath]
		switch {
		case ext == ".png" && !isInZip:
			t.Fatalf("%q file should be in zipped data", filePath)
		case ext != ".png" && isInZip:
			t.Fatalf("%q file should NOT be in zipped data", filePath)
		}
	}
}

func TestZip_WithExclude_bin_Dir(t *testing.T) {
	var data bytes.Buffer
	if err := Zip(&data, testDataDir, WithIncludePatterns("bin/**")); err != nil {
		t.Fatal(err)
	}

	files := testDataDirFiles()
	zFiles := zippedFiles(data.Bytes())

	for filePath := range files {
		dir := filepath.Dir(filePath)
		_, isInZip := zFiles[filePath]
		switch {
		case dir == "bin" && !isInZip:
			t.Fatalf("%q file from bin dir should NOT be in zipped data", filePath)
		case dir != "bin" && isInZip:
			t.Fatalf("%q file from bin dir should be in zipped data", filePath)
		}
	}
}

func testDataDirFiles() map[string]os.FileInfo {
	files := make(map[string]os.FileInfo)
	err := filepath.Walk(testDataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(testDataDir, path)
		files[relPath] = info
		return nil
	})
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		panic("there are no files in the test data dir")
	}
	return files
}

func zippedFiles(data []byte) map[string]os.FileInfo {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		panic(err)
	}
	files := make(map[string]os.FileInfo)
	for _, zippedFile := range zipReader.File {
		files[zippedFile.Name] = zippedFile.FileInfo()
	}
	return files
}
