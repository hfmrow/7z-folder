// gohAssets.go

/*
	Source file auto-generated on Mon, 26 Apr 2021 08:45:09 using Gotk3 Objects Handler v1.7.8
	©2018-21 hfmrow https://hfmrow.github.io

	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 hfmrow - 7z folder v1.6 github.com/hfmrow/7z-folder

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

/************************************************************************/
/* NOTICE: This file is completely rewritten when updating the assets. */
/**********************************************************************/

package main

import (
	"embed"
	"log"
)

// This functionality does not require explicit encoding of the files, at each
// compilation, the files are inserted into the resulting binary. Thus, updating
// assets is only required when new files are added to be embedded in order to
// create and declare the variables to which the files are linked.
// assetsDeclarationsUseEmbedded: Use native Go 'embed' package to include files
// content at runtime.
func assetsDeclarationsUseEmbedded(embedded ...bool) {
	mainGlade = readEmbedFile("assets/glade/main.glade")
	icon7zFolder500x48 = readEmbedFile("assets/images/icon-7zFolder-500x48.png")
	icon7zip48 = readEmbedFile("assets/images/icon-7zip-48.png")
	logOut20 = readEmbedFile("assets/images/Log-Out-20.png")
	signCancel20 = readEmbedFile("assets/images/Sign-cancel-20.png")
	signSelect20 = readEmbedFile("assets/images/Sign-Select-20.png")
	zipYellow20 = readEmbedFile("assets/images/zip-yellow-20.png")
}

// readEmbedFile: read 'embed' file system and return []byte data.
func readEmbedFile(filename string) (out []byte) {
	var err error
	out, err = embeddedFiles.ReadFile(filename)
	if err != nil {
		log.Printf("Unable to read embedded file: %s, %v\n", filename, err)
	}
	return
}

//go:embed assets/glade/main.glade
//go:embed assets/images/icon-7zFolder-500x48.png
//go:embed assets/images/icon-7zip-48.png
//go:embed assets/images/Log-Out-20.png
//go:embed assets/images/Sign-cancel-20.png
//go:embed assets/images/Sign-Select-20.png
//go:embed assets/images/zip-yellow-20.png
var embeddedFiles embed.FS
