# Zipper

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
