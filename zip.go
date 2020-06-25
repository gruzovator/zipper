package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// timestamp to use when modification time should be ignored
var defaultModTime = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

type zipper struct {
	srcPath        string
	filter         Filter
	ignoreModTimes bool
	useFileMode    os.FileMode
	writer         *zip.Writer
}

type Option func(*zipper)

func WithIncludePatterns(patterns []string) Option {
	return func(z *zipper) {
		z.filter.IncludePatterns = patterns
	}
}

func WithIncludePatternsStr(patterns string) Option {
	return WithIncludePatterns(splitPatternsString(patterns))
}

func WithExcludePatterns(patterns []string) Option {
	return func(z *zipper) {
		z.filter.ExcludePatterns = patterns
	}
}

func WithExcludePatternsStr(patterns string) Option {
	return WithExcludePatterns(splitPatternsString(patterns))
}

func WithIgonreModTimes(flagValue bool) Option {
	return func(z *zipper) {
		z.ignoreModTimes = flagValue
	}
}

func WithFileMode(m os.FileMode) Option {
	return func(z *zipper) {
		z.useFileMode = m
	}
}

func Zip(w io.Writer, srcPath string, opts ...Option) error {
	zipOut := zip.NewWriter(w)
	defer zipOut.Close()

	z := &zipper{
		srcPath: srcPath,
		writer:  zipOut,
	}
	for _, o := range opts {
		o(z)
	}

	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == "." || info.IsDir() || !z.filter.Match(path) {
			return nil
		}
		return z.zip(path, info)
	})
	return err
}

func (z *zipper) zip(filePath string, fileInfo os.FileInfo) error {
	zipFileHeader, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	var relPath = filePath
	if z.srcPath != filePath {
		r, err := filepath.Rel(z.srcPath, filePath)
		if err != nil {
			return err
		}
		relPath = r
	}

	zipFileHeader.Name = filepath.ToSlash(relPath)
	zipFileHeader.Method = zip.Deflate
	if z.ignoreModTimes {
		zipFileHeader.Modified = defaultModTime
	}
	if z.useFileMode != 0 {
		zipFileHeader.SetMode(z.useFileMode)
	}

	w, err := z.writer.CreateHeader(zipFileHeader)
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)

	return err
}

func splitPatternsString(s string) []string {
	var patterns []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		patterns = append(patterns, p)
	}
	return patterns
}
