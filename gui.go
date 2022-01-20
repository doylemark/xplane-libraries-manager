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

var (
	red   = color.RGBA{237, 66, 69, 1}
	green = color.RGBA{79, 220, 122, 1}
)

func home(w fyne.Window) fyne.CanvasObject {
	selectedForInstall := binding.NewStringList()

	bar := widget.NewSeparator()
	installList := widget.NewLabel("")

	scanner := newScanner()
	scanner.scan()

	var libs []Library

	for name, url := range libraries {
		_, installed := scanner.installedLibraries[name]

		lib := Library{
			name:                 name,
			url:                  url,
			isSelectedForInstall: false,
			isInstalled:          installed,
		}

		libs = append(libs, lib)
	}

	list := widget.NewList(
		func() int {
			return len(libs)
		},
		func() fyne.CanvasObject {
			chk := widget.NewCheck("template str", func(b bool) {})

			chk.OnChanged = func(b bool) {
				var lib Library
				var libIndex int

				for i, l := range libs {
					if l.name == chk.Text {
						lib = l
						libIndex = i
					}
				}

				if b && !lib.isInstalled {
					libs[libIndex].isSelectedForInstall = true
					selectedForInstall.Append(chk.Text)
				} else {
					var newSelected []string
					s, err := selectedForInstall.Get()

					if err != nil {
						fmt.Println(err)
						return
					}

					for _, item := range s {
						if item != chk.Text {
							newSelected = append(newSelected, item)
						}
					}

					selectedForInstall.Set(newSelected)
				}

				installList.SetText(fmt.Sprint(selectedForInstall.Length()))
			}

			hb := container.NewHBox(
				chk,
				layout.NewSpacer(),
				canvas.NewText("", green),
			)

			return container.NewPadded(hb)
		},
		func(id widget.ListItemID, co fyne.CanvasObject) {
			lib := libs[id]
			cont := co.(*fyne.Container).Objects[0].(*fyne.Container)
			cont.Objects[0].(*widget.Check).Text = lib.name

			if lib.isInstalled {
				cont.Objects[0].(*widget.Check).Disable()
				cont.Objects[0].(*widget.Check).SetChecked(true)
				cont.Objects[2].(*canvas.Text).Text = "(installed)"
				cont.Objects[2].(*canvas.Text).Color = green
			} else {
				cont.Objects[0].(*widget.Check).Enable()

				if !lib.isSelectedForInstall {
					cont.Objects[0].(*widget.Check).SetChecked(false)
				}
			}

			if lib.isSelectedForInstall {
				cont.Objects[0].(*widget.Check).SetChecked(true)
			}

			// fmt.Println("updated", items[id].name)

			cont.Objects[0].(*widget.Check).Refresh()
			cont.Objects[2].(*canvas.Text).Refresh()
		},
	)

	// list.OnSelected = func(id widget.ListItemID) {
	// 	p := fyne.CurrentApp().Preferences().String("SimPath") + "/" + items[id].name
	// 	pLabel := widget.NewLabelWithStyle(p, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	// 	dialogContent := container.NewVBox(pLabel)

	// 	dialog.ShowCustom(items[id].name, "close", dialogContent, w)
	// 	list.UnselectAll()
	// }

	return container.NewBorder(container.NewVBox(installList, bar), nil, nil, nil, list)
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

	return container.NewPadded(container.NewVBox(heading, bar, pathLabel, pathEntry, pathBtn))
}

func createTabs(w fyne.Window) *container.AppTabs {
	home := container.NewTabItem("Home", home(w))
	home.Icon = theme.FolderIcon()

	settings := container.NewTabItem("Settings", settings(w))
	settings.Icon = theme.SettingsIcon()

	tabs := container.NewAppTabs(home, settings)

	return tabs
}
