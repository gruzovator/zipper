# Zipper

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/gruzovator/zipper)](https://goreportcard.com/report/github.com/gruzovator/zipper)

Zipper is a simple tool to pack directories into go file.

## Installation

```
go get -u github.com/gruzovator/zipper
```

## Usage Example

Command:
```
zipper -src assets-src -dest assets/assets.go -pkg assets -exclude *.bin
```

Output go file provides:

* var ZippedFiles []byte
* func NewZippedFilesFS() http.FileSystem
