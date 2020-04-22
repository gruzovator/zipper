package main

import (
	"path/filepath"
)

// Filter is a collection of include/exclude patterns for file names.
type Filter struct {
	IncludePatterns []string
	ExcludePatterns []string
}

// Match returns true if basename for the path matches filter patterns.
func (f Filter) Match(path string) bool {
	basePath := filepath.Base(path)

	if matchAnyPattern(f.ExcludePatterns, basePath) {
		return false
	}

	if len(f.IncludePatterns) != 0 && !matchAnyPattern(f.IncludePatterns, basePath) {
		return false
	}

	return true
}

func matchAnyPattern(patterns []string, filePath string) bool {
	for _, p := range patterns {
		if ok, _ := filepath.Match(p, filePath); ok {
			return true
		}
	}

	return false
}
