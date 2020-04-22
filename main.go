package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	var (
		src             string
		dest            string
		pkg             string
		variable        string
		includePatterns string
		excludePatterns string
	)
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprint(out, "\nZipper is a tool to pack directories into go file.\n\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&src, "src", "", "source dir or file path")
	flag.StringVar(&dest, "dest", "", "dest go file path or - to print packed data string to stdout")
	flag.StringVar(&pkg, "pkg", "", "go file package name")
	flag.StringVar(&variable, "var", "ZippedFiles", "name of const variable to store data (optional)")
	flag.StringVar(&includePatterns, "include", "",
		"list of filename patterns to include, e.g.: *.css,*.html (optional)")
	flag.StringVar(&excludePatterns, "exclude", "",
		"list of filename patterns to exclude, e.g.: *.txt,*.bin (optional)")
	flag.Parse()

	if src == "" || dest == "" || pkg == "" {
		flag.Usage()
		os.Exit(2)
	}

	var zippedFiles bytes.Buffer
	err := Zip(&zippedFiles, src,
		WithIncludePatternsStr(includePatterns),
		WithExcludePatternsStr(excludePatterns),
	)
	if err != nil {
		panic(err)
	}

	encodedData := base64.StdEncoding.EncodeToString(zippedFiles.Bytes())

	if dest == "-" {
		fmt.Print(encodedData)
		return
	}

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		panic(err)
	}

	outFile, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	tmpl := template.Must(template.New("").Parse(outFileTemplate))
	err = tmpl.Execute(outFile, map[string]interface{}{
		"pkg":  pkg,
		"var":  variable,
		"data": encodedData,
	})
	if err != nil {
		panic(err)
	}
}

const outFileTemplate = `package {{.pkg}}

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"net/http"

	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

func New{{.var}}FS() vfs.FileSystem {
	zipReader, err := zip.NewReader(bytes.NewReader(decoded{{.var}}), int64(len(decoded{{.var}})))
	if err != nil {
		panic(err)
	}

	return zipfs.New(&zip.ReadCloser{Reader: *zipReader}, "/")
}

func New{{.var}}HttpFS() http.FileSystem {
	return httpfs.New(New{{.var}}FS())
}

var decoded{{.var}} []byte

func init() {
	var err error
	decoded{{.var}}, err = base64.StdEncoding.DecodeString({{.var}})
	if err != nil {
		panic(err)
	}
}

const {{.var}} = ` + "`{{.data}}`" + `
`
