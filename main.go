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

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("Hello World")

	username, password := ReadUser("user.txt")

	if username == nil || password == nil {
		// get user to login again
	}

	creds, err := Login(*username, *password)

	if err != nil {
		// get user to login again
		panic("Login Failed")
	}

	fetcher := newLibraryFetcher()
	go fetcher.GetMasterLibraries(creds)

	for lib := range fetcher.progress {
		fmt.Printf("%#v\n\n", lib)
	}

	w.SetContent(widget.NewLabel("Hello World!"))
	// w.ShowAndRun()
}

func hasUser(m map[string]string) bool {
	_, name := m["username"]
	_, pass := m["password"]

	return (name && pass)
}
