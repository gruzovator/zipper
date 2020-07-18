package main

import (
	"fmt"
	"os"

	"github.com/gobwas/glob"
)

type GlobFilter struct {
	include []glob.Glob
	exclude []glob.Glob
}

func NewGlobFilter(includePatterns, excludePatterns []string) (GlobFilter, error) {
	var globFilter GlobFilter
	for _, p := range includePatterns {
		g, err := glob.Compile(p, os.PathSeparator)
		if err != nil {
			return globFilter, fmt.Errorf("bad include pattern: %s", p)
		}
		globFilter.include = append(globFilter.include, g)
	}
	for _, p := range excludePatterns {
		g, err := glob.Compile(p, os.PathSeparator)
		if err != nil {
			return globFilter, fmt.Errorf("bad exclude pattern: %s", p)
		}
		globFilter.exclude = append(globFilter.exclude, g)
	}
	return globFilter, nil
}

// Match return true if path doesn't matches any exclude pattern and matches any include pattern (if there include patterns)
func (f GlobFilter) Match(path string) bool {
	if matchAnyPattern(f.exclude, path) {
		return false
	}

	if len(f.include) != 0 && !matchAnyPattern(f.include, path) {
		return false
	}

	return true
}

func matchAnyPattern(globs []glob.Glob, filePath string) bool {
	for _, g := range globs {
		if g.Match(filePath) {
			return true
		}
	}

	return false
}
