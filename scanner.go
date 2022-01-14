package main

import (
	"fmt"
	"io/ioutil"
)

const (
	path = "/Users/markdoyle/X-Plane 11/Custom Scenery"
)

type Scanner struct {
	installedLibraries map[string]*Library
}

func newScanner() *Scanner {
	return &Scanner{
		installedLibraries: map[string]*Library{},
	}
}

func (s *Scanner) scan() {
	dirs, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Println(err)
	}

	for _, f := range dirs {
		if f.IsDir() {
			lib := &Library{
				name:        f.Name(),
				isInstalled: true,
			}

			s.installedLibraries[lib.name] = lib
		}
	}
}
