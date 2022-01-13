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
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func home(w fyne.Window, c chan Library) fyne.CanvasObject {
	var libs []Library

	heading := widget.NewLabelWithStyle("Home"+fmt.Sprint(len(libs)), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	bar := widget.NewSeparator()

	list := widget.NewList(
		func() int {
			return len(libs)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("Template Object"),
				layout.NewSpacer(),
				canvas.NewText("(installed)", color.RGBA{79, 220, 124, 1}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(libs[id].Name)
		},
	)
	go func() {
		for lib := range c {
			fmt.Println(lib.Name)
			libs = append(libs, lib)
			list.Refresh()
			heading.SetText("Home (" + fmt.Sprint(len(libs)) + " libraries found)")
		}
	}()

	return container.NewBorder(container.NewVBox(heading, bar), nil, nil, nil, list)
}

func settings(w fyne.Window) fyne.CanvasObject {
	simPathBinding := binding.NewString()
	simPathBinding.Set(fyne.CurrentApp().Preferences().String("SimPath"))

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

	return container.NewVBox(heading, bar, pathLabel, pathEntry, pathBtn)
}

func createTabs(w fyne.Window, c chan Library) *container.AppTabs {
	home := container.NewTabItem("Home", home(w, c))
	home.Icon = theme.FolderIcon()

	settings := container.NewTabItem("Settings", settings(w))
	settings.Icon = theme.SettingsIcon()

	tabs := container.NewAppTabs(home, settings)

	return tabs
}
