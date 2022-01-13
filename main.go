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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
)

const (
	sig   = "com.markdoyle.libmanager"
	title = "X-Plane Library Manager"
)

func main() {
	a := app.NewWithID(sig)
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow(title)

	if a.Preferences().String("SimPath") == "" {
		promptSimPath(a, w)
	}

	cookies := refreshAuth()

	if cookies == nil {
		// prompt login
		fmt.Println("you don't got no smoke")
	}

	fetcher := newLibraryFetcher()
	go fetcher.getMasterLibraries(cookies)

	tabs := createTabs(w, fetcher.progress)

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(500, 600))
	w.ShowAndRun()
}

func promptSimPath(a fyne.App, w fyne.Window) {
	info := dialog.NewInformation("X-Plane path", "Please select your X-Plane 11 path when prompted", w)
	info.SetOnClosed(func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, e error) {
			if lu == nil {
				closeInfo := dialog.NewInformation("Missing X-Plane Path", "No Path was provided, closing", w)
				closeInfo.SetOnClosed(func() {
					w.Close()
				})
				closeInfo.Show()
				return
			}

			a.Preferences().SetString("SimPath", lu.Path())
			dialog.ShowInformation("X-Plane Path Set", "X-Plane Path set to \""+lu.Path()+"\".\n You can modify your path in settings", w)
		}, w)
	})
	info.Show()
}
