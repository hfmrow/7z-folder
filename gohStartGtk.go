// gohStartGtk.go

/*
	Source file auto-generated on Fri, 16 Apr 2021 22:07:08 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 hfmrow - 7z folder v1.6 github.com/hfmrow/7z-folder
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

/*******************************/
/* Gtk3 Window Initialisation */
/*****************************/
func mainStartGtk(winTitle string, width, height int, center bool) {
	obj = new(MainControlsObj)
	gtk.Init(nil)
	if newBuilder(mainGlade) == nil {
		// Init tempDir and Remove it on quit if requested.
		if doTempDir {
			tempDir = tempMake(Name)
			defer os.RemoveAll(tempDir)
		}
		// Parse Gtk objects
		gladeObjParser()
		// Objects Signals initialisations
		signalsPropHandler()
		/* Fill control with images */
		assignImages()
		// Set Window Properties
		if center {
			obj.MainWindow.SetPosition(gtk.WIN_POS_CENTER)
		}
		obj.MainWindow.SetTitle(winTitle)
		obj.MainWindow.SetDefaultSize(width, height)
		obj.MainWindow.Connect("delete-event", windowDestroy)
		// Start main application ...
		mainApplication()
		//	Start Gui loop
		obj.MainWindow.ShowAll()
		gtk.Main()
	} else {
		log.Fatal("Builder initialisation error.")
	}
}

// windowDestroy: on closing/destroying the gui window.
func windowDestroy(window interface{}) {
	if onShutdown() {
		gtk.MainQuit()
	}
}
