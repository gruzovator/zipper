package main_test

import (
	"reflect"
	"testing"

	. "github.com/gruzovator/zipper"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		name   string
		filter Filter
		paths  []string
		want   []string
	}{
		{
			name:   "empty filter: any path matches",
			filter: Filter{},
			paths:  []string{"1.txt", "./assets/2.jpg", "3"},
			want:   []string{"1.txt", "./assets/2.jpg", "3"},
		},
		{
			name: "filter only jpg files",
			filter: Filter{
				IncludePatterns: []string{"*.jpg"},
			},
			paths: []string{"1.txt", "./assets/2.jpg", "3"},
			want:  []string{"./assets/2.jpg"},
		},
		{
			name: "filter all except jpg files",
			filter: Filter{
				ExcludePatterns: []string{"*.jpg"},
			},
			paths: []string{"1.txt", "./assets/2.jpg", "3"},
			want:  []string{"1.txt", "3"},
		},
		{
			name: "filter all jpg files except test*.jpg",
			filter: Filter{
				IncludePatterns: []string{"*.jpg"},
				ExcludePatterns: []string{"test*.jpg"},
			},
			paths: []string{"1.jpg", "./assets/2.jpg", "tests/test_1.jpg", "3.bin"},
			want:  []string{"1.jpg", "./assets/2.jpg"},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			var filteredPaths []string
			for _, path := range c.paths {
				if !c.filter.Match(path) {
					continue
				}
				filteredPaths = append(filteredPaths, path)
			}
			if !reflect.DeepEqual(c.want, filteredPaths) {
				t.Errorf("want: %v, got: %v", c.want, filteredPaths)
			}
		})
	}
}
