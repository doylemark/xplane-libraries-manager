package main

import (
	"fmt"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
)

const (
	path = "/Users/markdoyle/X-Plane 11/Custom Scenery"
)

type Scanner struct {
	installedLibraries map[string]Library
	changeSignal       chan struct{}
}

func newScanner() *Scanner {
	s := &Scanner{
		installedLibraries: map[string]Library{},
		changeSignal:       make(chan struct{}),
	}
	s.watch()

	return s
}

func (s *Scanner) scan() {
	dirs, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Println(err)
	}

	for _, f := range dirs {
		if f.IsDir() {
			lib := Library{
				name:        f.Name(),
				isInstalled: true,
			}

			s.installedLibraries[lib.name] = lib
		}
	}
}

func (s *Scanner) watch() {
	fmt.Println("Starting watcher")
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		fmt.Println(err)
		return
	}

	// defer watcher.Close()

	go func() {
		for {
			select {
			case e, ok := <-watcher.Events:
				if ok {
					fmt.Println(e, ok, "event")
					fmt.Println("update")
				}
			}
		}
	}()

	err = watcher.Add(path)

	if err != nil {
		fmt.Println("error", err)
	}
}
