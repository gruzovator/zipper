# Zipper

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/gruzovator/zipper)](https://goreportcard.com/report/github.com/gruzovator/zipper)

**Deprecated, use //go:embed feature**

Zipper is a tool to pack directories into go file.

## Installation

```
GO111MODULE=on go get github.com/gruzovator/zipper@latest
```
or
```
go get github.com/gruzovator/zipper
```
 

## Usage

Command to pack assets-src dir into assets/assets.go excluding `bin` dir and all `txt` file :
```
zipper -src assets-src -dest assets/assets.go -pkg assets -exclude bin/** -exclude **.txt
```

Output go file provides:

* var ZippedFiles []byte
* func NewZippedFilesFS() http.FileSystem

Exclude/Include patterns examples:

- all txt files (including subdirs): `**.txt`
- all txt files in base dir: `*.txt`
- data dir: `data/**`

Patterns format description: https://github.com/gobwas/glob  
