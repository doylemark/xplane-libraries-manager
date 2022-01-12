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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Home(w fyne.Window) fyne.CanvasObject {
	heading := widget.NewLabelWithStyle("Home", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	bar := widget.NewSeparator()

	return container.NewVBox(heading, bar)
}

func Settings(w fyne.Window) fyne.CanvasObject {
	simPathBinding := binding.NewString()
	simPathBinding.Set(fyne.CurrentApp().Preferences().String("SimPath"))
	trackVersionsBinding := binding.NewBool()
	trackVersionsBinding.Set(fyne.CurrentApp().Preferences().Bool("TrackVersions"))

	heading := widget.NewLabelWithStyle("Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	bar := widget.NewSeparator()

	pathLabel := widget.NewLabel("X-Plane Location")
	pathEntry := widget.NewEntry()
	pathEntry.Bind(simPathBinding)
	pathEntry.PlaceHolder = "C://X-Plane 11"
	pathBtn := widget.NewButton("Choose Directory", func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, e error) {
			fyne.CurrentApp().Preferences().SetString("SimPath", lu.Path())
			simPathBinding.Set(lu.Path())
		}, w)
	})

	trackBar := widget.NewSeparator()
	trackVersions := widget.NewLabel("Library Versioning")
	trackVersionsBox := widget.NewCheck("Track Library Versions", nil)
	trackVersionsBox.OnChanged = func(b bool) {
		fyne.CurrentApp().Preferences().SetBool("TrackVersions", b)
		fmt.Println(b)
		fmt.Println(trackVersionsBinding.Get())
		trackVersionsBinding.Set(b)
	}
	trackVersionsBox.Bind(trackVersionsBinding)
	trackVersionsBox.Show()

	return container.NewVBox(heading, bar, pathLabel, pathEntry, pathBtn, trackBar, trackVersions, trackVersionsBox)
}

func CreateTabs(w fyne.Window) *container.AppTabs {
	home := container.NewTabItem("Home", Home(w))
	home.Icon = theme.FolderIcon()

	settings := container.NewTabItem("Settings", Settings(w))
	settings.Icon = theme.SettingsIcon()

	tabs := container.NewAppTabs(home, settings)

	return tabs
}
