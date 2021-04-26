// main.go

/*
	Source file auto-generated on Mon, 26 Apr 2021 08:45:09 using Gotk3 Objects Handler v1.7.8
	©2018-21 hfmrow https://hfmrow.github.io

	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 hfmrow - 7z folder v1.6 github.com/hfmrow/7z-folder

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	/* Be or not to be ... in dev mode ... */
	devMode = true

	/* Build directory for tempDir */
	doTempDir = false

	/* Set to true when you choose using embedded assets functionality */
	assetsDeclarationsUseEmbedded(!devMode)

	/* Init work directory */
	absoluteRealPath, optFilename = getAbsRealPath()

	/* Init Options */
	opt = new(MainOpt)
	opt.Init()

	/* Read Options */
	err = opt.Read()
	if err != nil {
		fmt.Printf("%s\n%v\n", "Options file not found or error on parsing.", err)
	}

	/* Init gtk display */
	mainStartGtk(fmt.Sprintf("%s %s  %s %s %s",
		Name,
		Vers,
		"©"+YearCreat,
		Creat,
		LicenseAbrv),
		opt.MainWinWidth,
		opt.MainWinHeight, true)
}

func mainApplication() {
	/* Init AboutBox */
	About = AboutInfosNew(obj.MainWindow)
	About.ImageTopHeight = 40
	About.FillInfos(
		"About "+Name,
		Name,
		Vers,
		Creat,
		YearCreat,
		LicenseAbrv,
		LicenseShort,
		Repository,
		Descr,
		icon7zFolder500x48,
		signSelect20)

	/* Translate init. */
	translate = MainTranslateNew(absoluteRealPath + opt.LanguageFilename)
	if devMode {
		translate.Sentences = sts
		err := translate.write(absoluteRealPath + opt.LanguageFilename)
		if err != nil {
			fmt.Printf("%s\n%v\n", "Cannot write language file or error on parsing.", err)
		}
	}
	sts = translate.Sentences

	if CheckCmd(opt.SevenZip) {
		/* Init spin button*/
		SpinbuttonSetValues(obj.SpinbuttonCompressionLvl, 0, 9, 9)
		SpinbuttonSetValues(obj.SpinbuttonDictSize, 4, 4096, 64, 4)

		// Init files operations structure
		fos, err = FilesOpStructNew()
		if checkErr(err) {
			return
		}

		// init dnd
		dndStruct = DragNDropNew(obj.MainWindow, nil, func() {

			opt.cmdLineList = *dndStruct.FilesList

			if checkErr(storeFiles()) {
				return
			}
			obj.LabelDuration.SetLabel("")
			obj.LabelPackedSize.SetLabel("")
			obj.LabelFileTreated.SetLabel("")
			obj.LabelUnpackedSize.SetLabel("")
			obj.LabelOverwritten.SetLabel("")

			dispStatusbar()
		})

		/* Update gtk conctrols with stored values into opt */
		opt.UpdateObjects()

		/* Statusbar init. */
		statusbar = StatusBarStructureNew(obj.Statusbar, []string{sts["sbCurrPath"], sts["count"]})

		/* Cmdline args */
		opt.cmdLineList = cmdLineArg()
		if len(opt.cmdLineList) > 0 {
			storeFiles()
		}
	} else {
		DialogMessage(obj.MainWindow, sts["warning"],
			"\n\n<span foreground=\"#204A87FF\">"+sts["no7za"]+"  </span><a href=\"https://www.7-zip.org/\">https://www.7-zip.org/</a>\n\nOn Ubuntu, use <span background=\"#666666FF\">  <span foreground=\"#4EBC06FF\"><b>$</b></span> <span foreground=\"#FE6C00FF\"><b>sudo</b></span> <span foreground=\"#F57900FF\"><b>apt-get</b></span> <span foreground=\"#DDEEFFFF\">install <b>p7zip-full  </b></span></span> to install the required command.",
			nil, DLG_OPT.DLG_MARKUP|DLG_OPT.DLG_INFO|DLG_OPT.DLG_BTN_RELIEF_NONE, sts["ok"])
		os.Exit(1)
	}
}

/*************************************\
/* Executed just before closing all. */
/************************************/
func onShutdown() bool {
	var err error

	// Update 'opt' with GtkObjects and save it
	if err = opt.Write(); err == nil {
		// What you want to execute before closing the app.
		// Return:
		// 	true for exit applicaton
		//	false does not exit application
	}
	if err != nil {
		log.Fatalf("Unexpected error on exit: %v", err)
	}
	return true
}
