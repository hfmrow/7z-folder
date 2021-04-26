// gohOptions.go

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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"time"

	glfs "github.com/hfmrow/genLib/files"
	glfsfo "github.com/hfmrow/genLib/files/filesOperations"
	glss "github.com/hfmrow/genLib/slices"
	glsg "github.com/hfmrow/genLib/strings"
	gltsbh "github.com/hfmrow/genLib/tools/bench"
	gltsushe "github.com/hfmrow/genLib/tools/units/human_readable"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gidgcr "github.com/hfmrow/gtk3Import/dialog/chooser"
	gimc "github.com/hfmrow/gtk3Import/misc"
)

// App infos
var (
	Name         = "7z folder"
	Vers         = "v1.6"
	Descr        = "Compress folder and files to 7z format. Some options are available. The 7za used command store filesystem permissions (such as UNIX owner/group permissions or NTFS ACLs). Not designed for large backup/archival purposes. On Ubuntu, use 'sudo apt-get install p7zip-full' to install required command."
	Creat        = "hfmrow"
	YearCreat    = "2019-21"
	LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
	LicenseAbrv  = "License (MIT)"
	Repository   = "github.com/hfmrow/7z-folder"

	// Vars declarations
	opt *MainOpt
	err error
	absoluteRealPath,
	optFilename,
	tempDir string
	doTempDir,
	devMode,
	VSCode bool
	nl      = "\n"
	execCmd *exec.Cmd
	ticker  *time.Ticker

	// Lib mapping

	// AboutBox
	AboutInfosNew = gidg.AboutInfosNew
	About         *gidg.AboutInfos

	// D&D
	DragNDropNew = gimc.DragNDropNew
	dndStruct    *gimc.DragNDropStruct

	// Dialog
	DialogMessage         = gidg.DialogMessage
	DLG_OPT               = gidg.DialogMessageOptNew()
	FileChooser           = gidgcr.FileChooserN
	fcOptions             = gidgcr.FileChooserOptionsNew()
	StatusBarStructureNew = gimc.StatusBarStructureNew
	statusbar             *gimc.StatusBar

	// files
	ExtEnsure            = glfs.ExtEnsure
	EE_PRESERVE_FILENAME = glfs.EE_PRESERVE_FILENAME
	EE_DEFAULT           = glfs.EE_DEFAULT
	TruncatePath         = glfs.TruncatePath
	BaseNoExt            = glfs.BaseNoExt
	CheckOutsideRoot     = glfs.CheckOutsideRoot
	FilesOpStructNew     = glfsfo.FilesOpStructNew
	fos                  *glfsfo.FilesOpStruct

	// Text
	GetTextEOL           = glsg.GetTextEOL
	GetOsLineEnd         = glsg.GetOsLineEnd
	HumanReadableSizeAlt = gltsushe.HumanReadableSizeAlt
	// Check                = gltses.Check
	IsExistSlIface = glss.IsExistSlIface

	// Bench
	BenchNew = gltsbh.BenchNew

	// Misc
	baseDir string
)

type MainOpt struct {
	// File signature
	FileSign []string

	MainWinWidth,
	MainWinHeight int
	LanguageFilename,
	SevenZip string

	FollowSL,
	AutoInc,
	Overwrite,
	Append,
	Update,
	OverwriteBack,
	AppendBack,
	UpdateBack bool

	DictSize,
	CompLvl int

	PrevFilename string

	cmdLineList []string
}

// Main options initialisation
func (opt *MainOpt) Init() {

	opt.LanguageFilename = "assets/lang/eng.lang"
	opt.MainWinWidth = 100
	opt.MainWinHeight = 100

	opt.SevenZip = "7za"
	opt.DictSize = 32 // Good value 64m
	opt.CompLvl = 5   // usually set to 9

	// opt.Overwrite = true
	// opt.OverwriteBack = opt.Overwrite
}

// Variables -> Objects.
func (opt *MainOpt) UpdateObjects() {
	obj.MainWindow.Resize(opt.MainWinWidth, opt.MainWinHeight)

	obj.CheckbuttonFollowSymlink.SetActive(opt.FollowSL)
	obj.SpinbuttonCompressionLvl.SetValue(float64(opt.CompLvl))
	obj.SpinbuttonDictSize.SetValue(float64(opt.DictSize))

	// RadioButton handled with callback fn, depending on
	// 'CheckbuttonAutoInc' state
	if !opt.AutoInc {
		obj.RadioButtonNew.SetActive(opt.Overwrite)
		obj.RadioButtonAppend.SetActive(opt.Append)
		obj.RadioButtonUpdate.SetActive(opt.Update)
	} else {
		obj.RadioButtonNew.SetActive(opt.OverwriteBack)
		obj.RadioButtonAppend.SetActive(opt.AppendBack)
		obj.RadioButtonUpdate.SetActive(opt.UpdateBack)
		obj.CheckbuttonAutoInc.SetActive(opt.AutoInc)
	}
}

// Objects -> Variables.
func (opt *MainOpt) UpdateOptions() {
	opt.MainWinWidth, opt.MainWinHeight = obj.MainWindow.GetSize()

	opt.FollowSL = obj.CheckbuttonFollowSymlink.GetActive()
	opt.CompLvl = obj.SpinbuttonCompressionLvl.GetValueAsInt()
	opt.DictSize = obj.SpinbuttonDictSize.GetValueAsInt()

	// RadioButton handled with callback fn, depending on
	// 'CheckbuttonAutoInc' state
	if opt.AutoInc = obj.CheckbuttonAutoInc.GetActive(); !opt.AutoInc {
		opt.Overwrite = obj.RadioButtonNew.GetActive()
		opt.Append = obj.RadioButtonAppend.GetActive()
		opt.Update = obj.RadioButtonUpdate.GetActive()
	}
}

// Read Options from file
func (opt *MainOpt) Read() (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(optFilename); err == nil {
		err = json.Unmarshal(textFileBytes, &opt)
	}
	return err
}

// Write Options to file
func (opt *MainOpt) Write() (err error) {
	var jsonData []byte
	var out bytes.Buffer
	opt.UpdateOptions()
	opt.FileSign = []string{Name, Vers, "©" + YearCreat, Creat, Repository, LicenseAbrv}
	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), 0644)
		}
	}
	return err
}
