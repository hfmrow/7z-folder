// gohObjects.go

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
	"github.com/gotk3/gotk3/gtk"
)

// Control over all used objects from glade.
var obj *MainControlsObj

/****************************************/
/* Main structure Declaration. You may */
/* add your own declarations here,    */
/* all is preserved on updates       */
/************************************/
type MainControlsObj struct {
	ButtonCancel             *gtk.Button
	ButtonCompress           *gtk.Button
	ButtonExit               *gtk.Button
	CheckbuttonAutoInc       *gtk.CheckButton
	CheckbuttonFollowSymlink *gtk.CheckButton
	EvenboxTop               *gtk.EventBox
	FrameCompressionLvl      *gtk.Frame
	FrameLabelCompressionLvl *gtk.Label
	FrameLabelOptions        *gtk.Label
	FrameLabelWriteAs        *gtk.Label
	FrameOptions             *gtk.Frame
	FrameWriteAs             *gtk.Frame
	ImageMainTop             *gtk.Image
	LabelDictSize            *gtk.Label
	LabelDuration            *gtk.Label
	LabelFileTreated         *gtk.Label
	LabelOverwritten         *gtk.Label
	LabelPackedSize          *gtk.Label
	LabelUnpackedSize        *gtk.Label
	mainUiBuilder            *gtk.Builder
	MainWindow               *gtk.Window
	Progressbar              *gtk.ProgressBar
	RadioButtonAppend        *gtk.RadioButton
	RadioButtonNew           *gtk.RadioButton
	RadioButtonUpdate        *gtk.RadioButton
	SpinbuttonCompressionLvl *gtk.SpinButton
	SpinbuttonDictSize       *gtk.SpinButton
	Statusbar                *gtk.Statusbar
}

/*******************************************/
/* GtkObjects Initialisation. You may add */
/* your own declarations as you wish,    */
/* best way is to add them by grouping  */
/* objects names. Yours modifications  */
/* will be preserved on updates.      */
/*************************************/
// gladeObjParser: Initialise Gtk3 Objects into obj structure.
func gladeObjParser() {
	obj.ButtonCancel = loadObject("ButtonCancel").(*gtk.Button)
	obj.ButtonCompress = loadObject("ButtonCompress").(*gtk.Button)
	obj.ButtonExit = loadObject("ButtonExit").(*gtk.Button)
	obj.CheckbuttonAutoInc = loadObject("CheckbuttonAutoInc").(*gtk.CheckButton)
	obj.CheckbuttonFollowSymlink = loadObject("CheckbuttonFollowSymlink").(*gtk.CheckButton)
	obj.EvenboxTop = loadObject("EvenboxTop").(*gtk.EventBox)
	obj.FrameCompressionLvl = loadObject("FrameCompressionLvl").(*gtk.Frame)
	obj.FrameLabelCompressionLvl = loadObject("FrameLabelCompressionLvl").(*gtk.Label)
	obj.FrameLabelOptions = loadObject("FrameLabelOptions").(*gtk.Label)
	obj.FrameLabelWriteAs = loadObject("FrameLabelWriteAs").(*gtk.Label)
	obj.FrameOptions = loadObject("FrameOptions").(*gtk.Frame)
	obj.FrameWriteAs = loadObject("FrameWriteAs").(*gtk.Frame)
	obj.ImageMainTop = loadObject("ImageMainTop").(*gtk.Image)
	obj.LabelDictSize = loadObject("LabelDictSize").(*gtk.Label)
	obj.LabelDuration = loadObject("LabelDuration").(*gtk.Label)
	obj.LabelFileTreated = loadObject("LabelFileTreated").(*gtk.Label)
	obj.LabelOverwritten = loadObject("LabelOverwritten").(*gtk.Label)
	obj.LabelPackedSize = loadObject("LabelPackedSize").(*gtk.Label)
	obj.LabelUnpackedSize = loadObject("LabelUnpackedSize").(*gtk.Label)
	obj.MainWindow = loadObject("MainWindow").(*gtk.Window)
	obj.Progressbar = loadObject("Progressbar").(*gtk.ProgressBar)
	obj.RadioButtonAppend = loadObject("RadioButtonAppend").(*gtk.RadioButton)
	obj.RadioButtonNew = loadObject("RadioButtonNew").(*gtk.RadioButton)
	obj.RadioButtonUpdate = loadObject("RadioButtonUpdate").(*gtk.RadioButton)
	obj.SpinbuttonCompressionLvl = loadObject("SpinbuttonCompressionLvl").(*gtk.SpinButton)
	obj.SpinbuttonDictSize = loadObject("SpinbuttonDictSize").(*gtk.SpinButton)
	obj.Statusbar = loadObject("Statusbar").(*gtk.Statusbar)
}
