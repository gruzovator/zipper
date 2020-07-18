package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"text/template"
	"time"
)

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprint(out, "\nZipper is a tool to pack directories into go file.\n\n")
		fmt.Fprintf(out, "Version: %s\n\n", version())
		flag.PrintDefaults()
	}
	srcDir := flag.String("src", "", "source dir or file path (required)")
	destFile := flag.String("dest", "",
		"output go file path (required), use - to print zipped data string to stdout")
	pkgName := flag.String("pkg", "", "package name for output go file (required)")
	zipdataVarName := flag.String("var", "ZippedFiles", "name of variable to store data")
	ignoreModTimes := flag.Bool("ignore-modtimes", false, "use default constant for zipped files mod time")
	ignoreFileModes := flag.Bool("ignore-filemodes", false, "use default file mod fro zipped files")
	verboseMode := flag.Bool("verbose", false, "print zipped files paths")
	printVersion := flag.Bool("version", false, "print version")
	var (
		includePatterns strArrayFlags
		excludePatterns strArrayFlags
	)
	flag.Var(&includePatterns, "include",
		"glob pattern to include files (e.g. use **.txt to include only txt files), glob format: github.com/gobwas/glob")
	flag.Var(&excludePatterns, "exclude",
		"glob pattern to exclude files (e.g. use bin/** to exclude bin dir), glob format: github.com/gobwas/glob")
	flag.Parse()
	if *printVersion {
		fmt.Fprintln(os.Stderr, version())
		os.Exit(2)
	}
	if *srcDir == "" || *destFile == "" || *pkgName == "" {
		flag.Usage()
		os.Exit(2)
	}

	zipOptions := []ZipOption{
		WithExcludePatterns(excludePatterns...),
		WithIncludePatterns(includePatterns...),
	}
	if *ignoreModTimes {
		modTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
		zipOptions = append(zipOptions, WithModTime(modTime))
	}
	if *ignoreFileModes {
		zipOptions = append(zipOptions, WithFileMode(0644))
	}
	if *verboseMode {
		printFile := func(f string) {
			fmt.Println(f)
		}
		zipOptions = append(zipOptions, WithProgressCallback(printFile))
	}

	var zippedFiles bytes.Buffer
	if err := Zip(&zippedFiles, *srcDir, zipOptions...); err != nil {
		panic(err)
	}

	encodedData := base64.StdEncoding.EncodeToString(zippedFiles.Bytes())

	if *destFile == "-" {
		fmt.Print(encodedData)
		return
	}

	if err := os.MkdirAll(filepath.Dir(*destFile), os.ModePerm); err != nil {
		panic(err)
	}

	outFile, err := os.Create(*destFile)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	tmpl := template.Must(template.New("").Parse(outFileTemplate))
	err = tmpl.Execute(outFile, map[string]interface{}{
		"pkg":  *pkgName,
		"var":  *zipdataVarName,
		"data": encodedData,
	})
	if err != nil {
		panic(err)
	}
}

type strArrayFlags []string

func (i *strArrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *strArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return "unknown"
	}
	return info.Main.Version
}
