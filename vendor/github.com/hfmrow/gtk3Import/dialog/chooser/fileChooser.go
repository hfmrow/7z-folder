// fileChooser.go

/*
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 https://github.com/hfmrow - GtkFileChooser helper
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/****************************
* FileChooser implementation.
 ****************************/
var (
	UserAbortError = errors.New("Action aborted by user !")

	fileChooserAction = map[string]gtk.FileChooserAction{
		"select-folder": gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		"create-folder": gtk.FILE_CHOOSER_ACTION_CREATE_FOLDER,
		"open":          gtk.FILE_CHOOSER_ACTION_OPEN,
		"open-entry":    gtk.FILE_CHOOSER_ACTION_SAVE, // 'Open' with 'entry' like in 'save' dialog
		"save":          gtk.FILE_CHOOSER_ACTION_SAVE,
	}
)

type fcOpt int16

const (
	FC_DLG_DEFAULT fcOpt = 1 << iota
	FC_DLG_SEL_FLD
	FC_DLG_CREATE_FLD
	FC_DLG_OPEN
	FC_DLG_OPEN_ENTRY // Open with 'entry', like in 'save' dialog
	FC_DLG_SAVE
	FC_NOT_ASK_OVERWRITE
	FC_NOT_KEEP_ABOVE
	FC_PREVIEW
	FC_NOT_MODAL
	FC_MULTI_SELECT
	FC_NOT_LOCAL_ONLY
)

type fcOptions struct {
	FC_DLG_DEFAULT,
	FC_DLG_SEL_FLD,
	FC_DLG_CREATE_FLD,
	FC_DLG_OPEN,
	FC_DLG_OPEN_ENTRY,
	FC_DLG_SAVE,
	FC_NOT_ASK_OVERWRITE,
	FC_NOT_KEEP_ABOVE,
	FC_PREVIEW,
	FC_NOT_MODAL,
	FC_MULTI_SELECT,
	FC_NOT_LOCAL_ONLY fcOpt
}

// FileChooserOptionsNew: create a new structure containing the available
// options/dialogType used as the parameters for the 'FileChooserN' function.
func FileChooserOptionsNew() *fcOptions {
	fco := new(fcOptions)
	fco.FC_DLG_DEFAULT = FC_DLG_DEFAULT
	fco.FC_DLG_SEL_FLD = FC_DLG_SEL_FLD
	fco.FC_DLG_CREATE_FLD = FC_DLG_CREATE_FLD
	fco.FC_DLG_OPEN = FC_DLG_OPEN
	fco.FC_DLG_OPEN_ENTRY = FC_DLG_OPEN_ENTRY
	fco.FC_DLG_SAVE = FC_DLG_SAVE
	fco.FC_NOT_ASK_OVERWRITE = FC_NOT_ASK_OVERWRITE
	fco.FC_NOT_KEEP_ABOVE = FC_NOT_KEEP_ABOVE
	fco.FC_PREVIEW = FC_PREVIEW
	fco.FC_NOT_MODAL = FC_NOT_MODAL
	fco.FC_MULTI_SELECT = FC_MULTI_SELECT
	fco.FC_NOT_LOCAL_ONLY = FC_NOT_LOCAL_ONLY
	return fco
}

/*
  FileChooser: Display a file chooser dialog box.
   - By default, 'keepAbove', 'setModal', 'askOverwrite' are enabled.
       title = "": auto choice based on dialog type.

   - Parameters: title, filename string, options fcOpt
       where 'options' means: 'dlgType' as FC_DLG_XXX and FC_XXX for dialog options.
       i.e: "title dialog", "/home/user/.config/x.txt", FC_DLG_SAVE|FC_ASK_OVERWRITE|FC_MODAL'
       means open a saving dialog, ask for overwrite and set it as modal dialog.

       NOTE: 'FileChooserOptionsNew' gives a structure which contains all the available options
       and can be declared as a global variable and then used whenever needed.
*/
func FileChooserN(window *gtk.Window, title, filename string, options fcOpt) (outFilename string, result bool, err error) {
	var (
		preview,
		folder bool
		fileChooser *gtk.FileChooserDialog

		kpAbove,
		modal,
		overwrt,
		setLocalOnly = true, true, true, true
		setSelectMultiple = false

		firstBtn,
		scndBtn = "Cancel", "Ok"

		dlgType,
		tmpTitle string

		fcOptionsGet = func(o fcOpt) {
			// Dialog type
			if o&FC_DLG_SEL_FLD != 0 {
				dlgType = "select-folder"
			}
			if o&FC_DLG_CREATE_FLD != 0 {
				dlgType = "create-folder"
			}
			if o&FC_DLG_OPEN != 0 {
				dlgType = "open"
			}
			if o&FC_DLG_OPEN_ENTRY != 0 {
				dlgType = "open-entry"
			}
			if o&FC_DLG_SAVE != 0 {
				dlgType = "save"
			}
			// Options
			if o&FC_MULTI_SELECT != 0 {
				setSelectMultiple = true
			}
			if o&FC_NOT_LOCAL_ONLY != 0 {
				setLocalOnly = false
			}
			if o&FC_NOT_ASK_OVERWRITE != 0 {
				overwrt = false
			}
			if o&FC_NOT_KEEP_ABOVE != 0 {
				kpAbove = false
			}
			if o&FC_PREVIEW != 0 {
				preview = true
			}
			if o&FC_NOT_MODAL != 0 {
				modal = false
			}
		}
	)

	fcOptionsGet(options)
	switch dlgType {
	case "create-folder":
		tmpTitle = "Create folder"
		folder = true
	case "select-folder":
		tmpTitle = "Select directory"
		folder = true
	case "open", "open-entry":
		tmpTitle = "Select file to open"
	case "save":
		tmpTitle = "Select file to save"
	}
	if len(title) == 0 {
		title = tmpTitle
	}

	if fileChooser, err = gtk.FileChooserDialogNewWith2Buttons(title, window, fileChooserAction[dlgType],
		firstBtn, gtk.RESPONSE_CANCEL, scndBtn, gtk.RESPONSE_ACCEPT); err != nil {
		return
	}

	fileChooser.SetLocalOnly(setLocalOnly)
	fileChooser.SetSelectMultiple(setSelectMultiple)
	fileChooser.SetDoOverwriteConfirmation(overwrt)
	fileChooser.SetModal(modal)
	fileChooser.SetKeepAbove(kpAbove)
	fileChooser.SetSkipPagerHint(true)
	fileChooser.SetSkipTaskbarHint(true)

	if preview {
		if previewImage, err := gtk.ImageNew(); err == nil {
			previewImage.Show()
			var pixbuf *gdk.Pixbuf
			fileChooser.SetPreviewWidget(previewImage)
			fileChooser.Connect("update-preview", func(fc *gtk.FileChooserDialog) {
				if _, err = os.Stat(fc.GetFilename()); !os.IsNotExist(err) {
					if pixbuf, err = gdk.PixbufNewFromFile(fc.GetFilename()); err == nil {
						fileChooser.SetPreviewWidgetActive(true)
						if pixbuf.GetWidth() > 640 || pixbuf.GetHeight() > 480 {
							if pixbuf, err = gdk.PixbufNewFromFileAtScale(fc.GetFilename(), 200, 200, true); err != nil {
								log.Fatalf("Image '%s' cannot be loaded, got error: %s", fc.GetFilename(), err.Error())
							}
						}
						previewImage.SetFromPixbuf(pixbuf)
					} else {
						fileChooser.SetPreviewWidgetActive(false)
					}
				}
			})
		}
	}
	fileChooser.SetFilename(filename)

	if dlgType == "save" {
		fileChooser.SetCurrentName(filepath.Base(filename))
	}
	if folder {
		fileChooser.SetCurrentFolder(filename)
	} else {
		fileChooser.SetCurrentFolder(filepath.Dir(filename))
	}

	if dlgType == "open-entry" {
		fileChooser.SetDoOverwriteConfirmation(false)
		fileChooser.SetCurrentName(filepath.Base(filename))
	}
	resp := fileChooser.Run()

	switch resp {
	case gtk.RESPONSE_CANCEL, gtk.RESPONSE_DELETE_EVENT:
		result = false
		err = UserAbortError
	case gtk.RESPONSE_ACCEPT:
		result = true
		outFilename = fileChooser.GetFilename()

		if dlgType == "open-entry" {
			if _, err = os.Stat(outFilename); err != nil {
				result = false
			}
		}
	}

	fileChooser.Destroy()
	return
}

// FileChooser: Display a file chooser dialog.
// dlgType: "open", "save", "create-folder", "select-folder" as dlgType.
// title = "": auto choice based on dialog type.
// options: 1-keepAbove, 2-enablePreviewImages, 3-setModal, 4-askOverwrite
// Default:	1- true, 2- false, 3- true, 4- true
func FileChooser(window *gtk.Window, dlgType, title, filename string, options ...bool) (outFilename string, result bool, err error) {

	log.Printf("WARNING: [FileChooser] Outdated version for this function, update to 'FileChooserN' is required as soon as possible...")

	var (
		preview,
		folder bool
		fileChooser *gtk.FileChooserDialog
		kpAbove,
		modal,
		overwrt = true, true, true
		firstBtn,
		scndBtn = "Cancel", "Ok"
		tmpTitle string
	)

	for idx, opt := range options {
		switch idx {
		case 0:
			kpAbove = opt
		case 1:
			preview = opt
		case 2:
			modal = opt
		case 3:
			overwrt = opt
		}
	}

	switch dlgType {
	case "create-folder":
		tmpTitle = "Create folder"
		folder = true
	case "select-folder":
		tmpTitle = "Select directory"
		folder = true
	case "open", "open-entry":
		tmpTitle = "Select file to open"
	case "save":
		tmpTitle = "Select file to save"
	}
	if len(title) == 0 {
		title = tmpTitle
	}

	if fileChooser, err = gtk.FileChooserDialogNewWith2Buttons(title, window, fileChooserAction[dlgType],
		firstBtn, gtk.RESPONSE_CANCEL, scndBtn, gtk.RESPONSE_ACCEPT); err != nil {
		return
	}

	if preview {
		if previewImage, err := gtk.ImageNew(); err == nil {
			previewImage.Show()
			var pixbuf *gdk.Pixbuf
			fileChooser.SetPreviewWidget(previewImage)
			fileChooser.Connect("update-preview", func(fc *gtk.FileChooserDialog) {
				if _, err = os.Stat(fc.GetFilename()); !os.IsNotExist(err) {
					if pixbuf, err = gdk.PixbufNewFromFile(fc.GetFilename()); err == nil {
						fileChooser.SetPreviewWidgetActive(true)
						if pixbuf.GetWidth() > 640 || pixbuf.GetHeight() > 480 {
							if pixbuf, err = gdk.PixbufNewFromFileAtScale(fc.GetFilename(), 200, 200, true); err != nil {
								log.Fatalf("Image '%s' cannot be loaded, got error: %s", fc.GetFilename(), err.Error())
							}
						}
						previewImage.SetFromPixbuf(pixbuf)
					} else {
						fileChooser.SetPreviewWidgetActive(false)
					}
				}
			})
		}
	}

	fileChooser.SetFilename(filename)

	if dlgType == "save" {
		fileChooser.SetCurrentName(filepath.Base(filename))
	}

	if folder {
		fileChooser.SetCurrentFolder(filename)
	} else {
		fileChooser.SetCurrentFolder(filepath.Dir(filename))
	}
	fileChooser.SetDoOverwriteConfirmation(overwrt)
	fileChooser.SetModal(modal)
	fileChooser.SetSkipPagerHint(true)
	fileChooser.SetSkipTaskbarHint(true)
	fileChooser.SetKeepAbove(kpAbove)

	if dlgType == "open-entry" {
		fileChooser.SetDoOverwriteConfirmation(false)
		fileChooser.SetCurrentName(filepath.Base(filename))
	}
	resp := fileChooser.Run()
	switch resp {
	case gtk.RESPONSE_CANCEL, gtk.RESPONSE_DELETE_EVENT:

		result = false
		err = UserAbortError
	case gtk.RESPONSE_ACCEPT:

		result = true
		outFilename = fileChooser.GetFilename()

		if dlgType == "open-entry" {
			if _, err = os.Stat(outFilename); err != nil {
				result = false
			}
		}
	}

	fileChooser.Destroy()
	return
}
