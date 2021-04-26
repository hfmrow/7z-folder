// translate.go

// File generated on Fri, 16 Apr 2021 22:27:39 using Gotk3 Objects Translate v1.5.2 2019-21 hfmrow

/*
* 	This program comes with absolutely no warranty.
*	See the The MIT License (MIT) for details:
*	https://opensource.org/licenses/mit-license.php
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

// initGtkObjectsText: read translations from structure and set them to objects.
func (trans *MainTranslate) initGtkObjectsText() {
	trans.setTextToGtkObjects(&obj.ButtonCancel.Widget, "ButtonCancel")
	trans.setTextToGtkObjects(&obj.ButtonCompress.Widget, "ButtonCompress")
	trans.setTextToGtkObjects(&obj.ButtonExit.Widget, "ButtonExit")
	trans.setTextToGtkObjects(&obj.CheckbuttonAutoInc.Widget, "CheckbuttonAutoInc")
	trans.setTextToGtkObjects(&obj.CheckbuttonFollowSymlink.Widget, "CheckbuttonFollowSymlink")
	trans.setTextToGtkObjects(&obj.FrameLabelCompressionLvl.Widget, "FrameLabelCompressionLvl")
	trans.setTextToGtkObjects(&obj.FrameLabelOptions.Widget, "FrameLabelOptions")
	trans.setTextToGtkObjects(&obj.FrameLabelWriteAs.Widget, "FrameLabelWriteAs")
	trans.setTextToGtkObjects(&obj.ImageMainTop.Widget, "ImageMainTop")
	trans.setTextToGtkObjects(&obj.LabelDictSize.Widget, "LabelDictSize")
	trans.setTextToGtkObjects(&obj.LabelDuration.Widget, "LabelDuration")
	trans.setTextToGtkObjects(&obj.LabelFileTreated.Widget, "LabelFileTreated")
	trans.setTextToGtkObjects(&obj.LabelOverwritten.Widget, "LabelOverwritten")
	trans.setTextToGtkObjects(&obj.LabelPackedSize.Widget, "LabelPackedSize")
	trans.setTextToGtkObjects(&obj.LabelUnpackedSize.Widget, "LabelUnpackedSize")
	trans.setTextToGtkObjects(&obj.Progressbar.Widget, "Progressbar")
	trans.setTextToGtkObjects(&obj.RadioButtonAppend.Widget, "RadioButtonAppend")
	trans.setTextToGtkObjects(&obj.RadioButtonNew.Widget, "RadioButtonNew")
	trans.setTextToGtkObjects(&obj.RadioButtonUpdate.Widget, "RadioButtonUpdate")
	trans.setTextToGtkObjects(&obj.SpinbuttonCompressionLvl.Widget, "SpinbuttonCompressionLvl")
	trans.setTextToGtkObjects(&obj.SpinbuttonDictSize.Widget, "SpinbuttonDictSize")
	trans.setTextToGtkObjects(&obj.Statusbar.Widget, "Statusbar")
}

// Translations structure declaration. To be used in main application.
var translate = new(MainTranslate)

// sts: some sentences/words used in the application. Mostly used in Development mode.
// You must add there all sentences used in your application. Or not ...
// They'll be added to language file each time application started
// when "devMode" is set at true.
var sts = map[string]string{
	`no`:    `No`,
	`error`: `An error occured`,
	`info`:  `Information`,
	`no7za`: `7Zip command not present (7za), please install it
before trying using this program again.

Informations can be found at:`,
	`outBaseDir3`: `The previous base directory was:`,
	`continue`:    `Continue`,
	`yes`:         `Yes`,
	`errExec`:     `File execution error`,
	`brkSl`: `Some broken symbolic links were found:
continue and ignore this alert ?`,
	`savef`: `Save file`,
	`archUpdt`: `archive
updated`,
	`for`:             `For`,
	`sbCurrFile`:      `Current language file: `,
	`permissionError`: `Some files appear unreadable or inaccessible.`,
	`noArgs`:          `You haven't chosen any file to compress.`,
	`warning`:         `Warning`,
	`notAllowed`:      `This is not allowed, please retry with different manner.`,
	`outBaseDir2`:     `upstream directories from the current base.`,
	`allow`:           `Allow`,
	`noFile`:          `There is no file / dir to compress.`,
	`sbCurrPath`:      `Path:`,
	`noWCFld`:         `No way to choose this folder`,
	`outBaseDir4`:     `The new base directory will be changed to:`,
	`cancel`:          `Cancel`,
	`openf`:           `Open file`,
	`unexpectedError`: `Unexpected error: `,
	`deny`:            `Deny`,
	`archOvrWrt`: `archive
overwritten`,
	`noWCFile`:      `No way to choose this file`,
	`unexpectedEnd`: `An unexpected issue has been encountered:`,
	`sbWFiles`:      `File(s)`,
	`archAppnd`: `Added files
to the archive`,
	`retry`:       `Retry`,
	`added`:       `Added`,
	`outBaseDir1`: `The recently added files seem to be`,
	`ok`:          `Ok`,
	`confirm`:     `Ask for confirmation`,
	`count`:       `Count:`,
	`removed`:     `Removed`,
	`useFName`:    `Use this filename ?`,
}

// Translations structure with methods
type MainTranslate struct {
	// Public
	ProgInfos    progInfo
	Language     language
	Options      parsingFlags
	ObjectsCount int
	Objects      []object
	Sentences    map[string]string
	// Private
	objectsLoaded bool
}

// MainTranslateNew: Initialise new translation structure and assign language file content to GtkObjects.
// devModeActive, indicate that the new sentences must be added to previous language file.
func MainTranslateNew(filename string, devModeActive ...bool) (mt *MainTranslate) {
	var err error
	mt = new(MainTranslate)
	if err = mt.read(filename); err == nil {
		mt.initGtkObjectsText()
		if len(devModeActive) != 0 {
			if devModeActive[0] {
				mt.Sentences = sts
				err := mt.write(filename)
				if err != nil {
					fmt.Printf("%s\n%s\n", "Cannot write actual sentences to language file.", err.Error())
				}
			}
		}
	} else {
		fmt.Printf("%s\n%s\n", "Error loading language file !\nNot an error when you just creating from glade Xml or GOH project file.", err.Error())
	}
	return
}

// readFile: language file.
func (trans *MainTranslate) read(filename string) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		if err = json.Unmarshal(textFileBytes, &trans); err == nil {
			trans.objectsLoaded = true
		}
	}
	return
}

// Write json datas to file
func (trans *MainTranslate) write(filename string) (err error) {
	var out bytes.Buffer
	var jsonData []byte
	if jsonData, err = json.Marshal(&trans); err == nil && trans.objectsLoaded {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), 0644)
		}
	}
	return
}

type parsingFlags struct {
	SkipLowerCase,
	SkipEmptyLabel,
	SkipEmptyName,
	DoBackup bool
}

type progInfo struct {
	Name,
	Version,
	Creat,
	MainObjStructName,
	GladeXmlFilename,
	TranslateFilename,
	ProjectRootDir,
	GohProjFile string
}

type language struct {
	LangNameLong,
	LangNameShrt,
	Author,
	Date,
	Updated string
	Contributors []string
}

type object struct {
	Class,
	Id,
	Label,
	Tooltip,
	Text,
	Uri,
	Comment string
	LabelMarkup,
	LabelWrap,
	TooltipMarkup bool
	Idx int
}

// Define available property within objects
type propObject struct {
	Class string
	Label,
	LabelMarkup,
	LabelWrap,
	Tooltip,
	TooltipMarkup,
	Text,
	Uri bool
}

// Property that exists for Gtk3 Object ...	(Used for Class capability)
var propPerObjects = []propObject{
	{Class: "GtkButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkMenuButton", Label: true, Tooltip: true, TooltipMarkup: true},

	// {Class: "GtkToolButton", Label: true, Tooltip: true, TooltipMarkup: true},    // Deprecated since 3.10
	// {Class: "GtkImageMenuItem", Label: true, Tooltip: true, TooltipMarkup: true}, // Deprecated since 3.10

	{Class: "GtkMenuItem", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkCheckMenuItem", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkRadioMenuItem", Label: true, Tooltip: true, TooltipMarkup: true},

	{Class: "GtkToggleButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkLabel", Label: true, LabelMarkup: true, Tooltip: true, TooltipMarkup: true, LabelWrap: true},
	{Class: "GtkSpinButton", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkEntry", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkCheckButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkProgressBar", Tooltip: true, TooltipMarkup: true, Text: true},
	{Class: "GtkSearchBar", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkImage", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkRadioButton", Label: true, LabelMarkup: false, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkComboBoxText", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkComboBox", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkLinkButton", Label: true, Tooltip: true, TooltipMarkup: true, Uri: true},
	{Class: "GtkSwitch", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkTreeView", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkFileChooserButton", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkTextView", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkSourceView", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkStatusbar", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkScrolledWindow", Tooltip: true, TooltipMarkup: true},
}

// setTextToGtkObjects: read translations from structure and set them to object.
// like this: setTextToGtkObjects(&mainObjects.TransLabelHint.Widget, "TransLabelHint")
func (trans *MainTranslate) setTextToGtkObjects(obj *gtk.Widget, objectId string) {
	for _, currObject := range trans.Objects {
		if currObject.Id == objectId {
			for _, props := range propPerObjects {
				if currObject.Class == props.Class {
					if props.Label {
						obj.SetProperty("label", currObject.Label)
						if props.LabelMarkup {
							obj.SetProperty("use-markup", currObject.LabelMarkup)
							obj.SetProperty("label", strings.ReplaceAll(currObject.Label, "&", "&amp;"))
						}
					}
					if props.LabelWrap {
						obj.SetProperty("wrap", currObject.LabelWrap)
					}
					if props.Tooltip && !currObject.TooltipMarkup {
						obj.SetProperty("tooltip_text", currObject.Tooltip)
					}
					if props.Tooltip && currObject.TooltipMarkup {
						obj.SetProperty("tooltip_markup", strings.ReplaceAll(currObject.Tooltip, "&", "&amp;"))
					}
					if props.Text {
						obj.SetProperty("text", currObject.Text)
					}
					if props.Uri {
						obj.SetProperty("uri", currObject.Uri)
					}
					break
				}
			}
			break
		}
	}
}
