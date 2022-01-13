/*
 * CDDL HEADER START
 *
 * This file and its contents are supplied under the terms of the
 * Common Development and Distribution License ("CDDL"), version 1.0.
 * You may only use this file in accordance with the terms of version
 * 1.0 of the CDDL.
 *
 * A full copy of the text of the CDDL should have accompanied this
 * source.  A copy of the CDDL is also available via the Internet at
 * http://www.illumos.org/license/CDDL.
 *
 * Copyright 2022 Mark Doyle. All rights reserved.
 */

package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	masterListUrl = "https://forums.x-plane.org/index.php?/forums/topic/90776-master-list-of-libraries/"
)

type LibraryFetcher struct {
	existingLibCount int
	progress         chan Library
	Libs             []Library
}

type Library struct {
	Name        string
	DownloadUrl string
	Url         string
	Version     string
}

func newLibraryFetcher() *LibraryFetcher {
	return &LibraryFetcher{
		progress: make(chan Library),
	}
}

func (fetcher *LibraryFetcher) getMasterLibraries(creds []*http.Cookie) error {
	resp, err := http.Get(masterListUrl)

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return err
	}

	frames := doc.Find("iframe")
	fetcher.existingLibCount = frames.Length()

	var libs []Library
	frames.Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("data-embed-src")

		if exists {
			go func() {
				lib := getLibraryInfo(src, creds)

				if lib != nil {
					fetcher.progress <- *lib
					libs = append(libs, *lib)
				}
			}()
		}
	})

	return nil
}

func getLibraryInfo(url string, creds []*http.Cookie) *Library {
	resp, err := makeAuthorizedGet(creds, url)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	lib := &Library{
		Url:  url,
		Name: doc.Find("title").Text(),
	}

	// find lib version
	doc.Find("span").Each(func(i int, s *goquery.Selection) {
		role, exists := s.Attr("data-role")

		if exists && role == "versionTitle" {
			lib.Version = s.Text()
		} else {
			lib.Version = "1.0.0"
		}
	})

	// find lib downloadurl
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("href")

		if s.Text() == "Download this file" && exists {
			lib.DownloadUrl = src
		}
	})

	if lib.DownloadUrl == "" || lib.Name == "" {
		return nil
	}

	return lib
}
