package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"time"
)

type zipOptions struct {
	excludePatterns  []string
	includePatterns  []string
	useModTime       *time.Time
	useFileMode      *os.FileMode
	progressCallback func(file string)
}

type ZipOption func(*zipOptions)

func WithExcludePatterns(pp ...string) ZipOption {
	return func(o *zipOptions) {
		o.excludePatterns = append(o.excludePatterns, pp...)
	}
}

func WithIncludePatterns(pp ...string) ZipOption {
	return func(o *zipOptions) {
		o.includePatterns = append(o.includePatterns, pp...)
	}
}

func WithModTime(t time.Time) ZipOption {
	return func(o *zipOptions) {
		o.useModTime = &t
	}
}

func WithFileMode(m os.FileMode) ZipOption {
	return func(o *zipOptions) {
		o.useFileMode = &m
	}
}

func WithProgressCallback(fn func(filePath string)) ZipOption {
	return func(o *zipOptions) {
		o.progressCallback = fn
	}
}

func Zip(w io.Writer, srcPath string, opts ...ZipOption) error {
	var options zipOptions
	for _, o := range opts {
		o(&options)
	}

	filter, err := NewGlobFilter(options.includePatterns, options.excludePatterns)
	if err != nil {
		return err
	}

	zipOut := zip.NewWriter(w)
	defer zipOut.Close()

	z := &zipper{
		writer:      zipOut,
		srcPath:     srcPath,
		filter:      filter,
		useModTime:  options.useModTime,
		useFileMode: options.useFileMode,
	}

	progressCallback := func(string) {}
	if options.progressCallback != nil {
		progressCallback = options.progressCallback
	}

	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(srcPath, path)
		if path == "." || info.IsDir() || !z.filter.Match(relPath) {
			return nil
		}
		progressCallback(relPath)
		return z.zip(path, info)
	})
	return err
}

type zipper struct {
	writer      *zip.Writer
	srcPath     string
	filter      GlobFilter
	useModTime  *time.Time
	useFileMode *os.FileMode
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
	if z.useModTime != nil {
		zipFileHeader.Modified = *z.useModTime
	}
	if z.useFileMode != nil {
		zipFileHeader.SetMode(*z.useFileMode)
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
