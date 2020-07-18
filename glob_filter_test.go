package main

import (
	"reflect"
	"testing"
)

func TestGlobFilter_Match(t *testing.T) {
	cases := []struct {
		name      string
		include   []string
		exclude   []string
		paths     []string
		wantMatch []string
	}{
		{
			name:      "all paths should match empty filter",
			include:   nil,
			exclude:   nil,
			paths:     []string{"1.txt", "2/2.jpg", ".data"},
			wantMatch: []string{"1.txt", "2/2.jpg", ".data"},
		},
		{
			name:      "match all jpg files",
			include:   []string{"**.jpg"},
			exclude:   nil,
			paths:     []string{"1.txt", "2/2.jpg", ".data/1.jpg", "2.jpg"},
			wantMatch: []string{"2/2.jpg", ".data/1.jpg", "2.jpg"},
		},
		{
			name:      "exclude all files from data dir",
			include:   nil,
			exclude:   []string{"data/**"},
			paths:     []string{"1.txt", "data/1.jpg", "data/1.bin", "data/bin/1.txt"},
			wantMatch: []string{"1.txt"},
		},
		{
			name:      "exclude all files from data and include only txt files",
			include:   []string{"**.txt"},
			exclude:   []string{"data/**"},
			paths:     []string{"1.txt", "data/1.jpg", "data/1.bin", "data/bin/1.txt"},
			wantMatch: []string{"1.txt"},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			f, err := NewGlobFilter(c.include, c.exclude)
			if err != nil {
				t.Fatal(err)
			}
			var matchedPaths []string
			for _, path := range c.paths {
				if !f.Match(path) {
					continue
				}
				matchedPaths = append(matchedPaths, path)
			}
			if !reflect.DeepEqual(c.wantMatch, matchedPaths) {
				t.Errorf("want: %v, got: %v", c.wantMatch, matchedPaths)
			}
		})
	}
}
