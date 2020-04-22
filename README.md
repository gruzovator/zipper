# Zipper

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/gruzovator/zipper)](https://goreportcard.com/report/github.com/gruzovator/zipper)

Zipper is a simple tool to pack directories into go file.

## Usage Example

Command:
```
zipper -src assets-src -dest assets/assests.go -pkg assets -exclude *.bin
```

Output go file provides:

* const ZippedFiles = `<base64 encoded zip>`
* func NewZippedFilesFS() vfs.FileSystem
* func NewZippedFilesHttpFS() http.FileSystem
