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

import "fmt"

// https://forums.x-plane.org/index.php?/files/file/17092-flags-of-the-usa-states-real-flag-ii/&do=download&r=197551&confirm=1&t=1&csrfKey=827dc6d2fddb3350ff71da85c0cad228
// https://forums.x-plane.org/index.php?/files/file/27545-the-fruit-stand-aircraft-library/&do=download&r=488580&confirm=1&t=1&csrfKey=827dc6d2fddb3350ff71da85c0cad228
var (
	libraries = map[string]string{
		"3D_people_library":                      "https://forums.x-plane.org/index.php?/files/file/26611-3d-people-library",
		"AR_Library":                             "https://forums.x-plane.org/index.php?/files/file/44586-ar_library-librer%C3%ADa-de-objetos-argentinos",
		"BS2001 Object Library":                  "https://forums.x-plane.org/index.php?/files/file/28045-bs2001-object-library",
		"CDB-Library":                            "https://www.mediafire.com/file/avnjtzt8hdy3n4g/CDB-Library.zip/file",
		"cemetry":                                "https://forums.x-plane.org/index.php?/files/file/24521-cemetery-library",
		"Europe_RoadTraffic":                     "https://forums.x-plane.org/index.php?/files/file/21719-custom-traffic-for-europe-library/r=218080&confirm=1&t=1",
		"european_vehicles_library_uwespeed":     "https://forums.x-plane.org/index.php?/files/file/24708-library-european-vehicles-static",
		"Faib":                                   "custom...",
		"ff_library_extended_LOD":                "https://forums.x-plane.org/index.php?/files/file/12836-ff-library-extended-lod-version",
		"FJS_Scenery_Library_v1.7":               "https://forums.x-plane.org/index.php?/files/file/28594-fjs_scenery_library",
		"flags_of_USA_states":                    "https://forums.x-plane.org/index.php?/files/file/17092-flags-of-the-usa-states-real-flag-ii/&confirm=1&t=1",
		"flags_of_the_world":                     "https://forums.x-plane.org/index.php?/files/file/17090-flags-of-the-world-real-flag-ii/r=197592&confirm=1&t=18",
		"FlyAgi_Vegetation":                      "https://storage.flyagi.de/flyagi/vegetation/FlyAgi_Vegetation.zip",
		"Flyby_Planes":                           "https://forums.x-plane.org/index.php?/files/file/28295-flyby-planes-library",
		"THE-FRUIT-STAND Aircraft Library v3.0":  "https://forums.x-plane.org/index.php?/files/file/27545-the-fruit-stand-aircraft-library/&r=488580&confirm=1&t=1",
		"german_traffic_library":                 "https://forums.x-plane.org/index.php?/files/file/19282-german-traffic-library",
		"gt_library":                             "https://forums.x-plane.org/index.php?/files/file/29461-ground-textures-library",
		"The_Handy_Objects_Library":              "https://forums.x-plane.org/index.php?/files/file/24261-the-handy-objects-library",
		"ZZ_DF-Hard_Surface":                     "https://forums.x-plane.org/index.php?/files/file/13129-hard-surface-library",
		"MisterX_Library":                        "https://forums.x-plane.org/index.php?/files/file/28167-misterx-library-and-static-aircraft-extension",
		"Serviced Aircraft North America Part 1": "https://forums.x-plane.org/index.php?/files/file/15363-north-american-serviced-aircraft-library-part-1",
		"Serviced Aircraft North America Part 2": "https://forums.x-plane.org/index.php?/files/file/15365-north-american-serviced-aircraft-library-part-2",
		"Serviced Aircraft North America Part 3": "https://forums.x-plane.org/index.php?/files/file/15386-north-american-serviced-aircraft-library-part-3",
		"ADSparks_library":                       "https://forums.x-plane.org/index.php?/files/file/47600-parking-stands-signs",
		"People_LIB":                             "http://www.x-plane.at/drupal/sites/www.x-plane.at/files/PEOPLE_LIB/People_LIB%201.11.zip",
		"PPlibrary":                              "https://forums.x-plane.org/index.php?/files/file/37088-pavement-paintings-library-pplibrary",
		"RE_Library V1.8":                        "https://forums.x-plane.org/index.php?/files/file/24722-re_library-airport-buildings-and-related-objects",
		"ruscenery":                              "https://ruscenery.x-air.ru/files/RuScenery.zip",
		"R2_Library":                             "http://r2.xpl.cz/download.php?id=3",
		"Sea_Life":                               "https://forums.x-plane.org/index.php?/files/file/28296-sea-life-library",
		"Static_GA_Aircraft_NZ":                  "https://forums.x-plane.org/index.php?/files/file/27150-static-ga-aircraft-new-zealand",
		"Waves_Library":                          "https://forums.x-plane.org/index.php?/files/file/25439-waves-library",
	}
)

type Library struct {
	name             string
	url              string
	isInstalled      bool
	downloadProgress float64
}

func (l Library) install() {
	fmt.Println("Installing @", l.url)
	l.downloadProgress = 100
	l.isInstalled = true
}

func getAllLibraries() []*Library {
	var libs []*Library

	scanner := newScanner()
	scanner.scan()

	for name, url := range libraries {
		_, ok := scanner.installedLibraries[name]

		lib := &Library{
			name:        name,
			url:         url,
			isInstalled: ok,
		}

		libs = append(libs, lib)
	}

	return libs
}
