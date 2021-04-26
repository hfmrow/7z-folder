// dialogMessage.go

/*
*	©2018-21 H.F.M. MIT license
*	Handle gtk3 Dialogs.
*
*	Message, Question/Response dialogs, file/dir dialogs, Notifications.
 */

package gtk3Import

import (
	"log"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	gipf "github.com/hfmrow/gtk3Import/pixbuff"
)

// _ExampleDialogMessage: Sample example.
func _ExampleDialogMessage() {

	parentWindow, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	DLG_OPT := DialogMessageOptNew()

	// Modal example
	resp := DialogMessage(
		parentWindow,
		"Confirmation",
		"Continue with this parameter ?",
		nil,
		DLG_OPT.DLG_INFO|
			DLG_OPT.DLG_BTN_RELIEF_NONE,
		"Cancel",
		"Continue")

	switch resp {
	case 0:
		log.Printf("Canceled: %d\n", resp)
	case 1:
		log.Printf("Continue clicked: %d\n", resp)
	}

	// Non-modal example
	DialogMessage(
		parentWindow,
		"Confirmation",
		"Continue with this parameter ?",
		nil,
		DLG_OPT.DLG_INFO|
			DLG_OPT.DLG_BTN_RELIEF_NONE|
			DLG_OPT.DLG_NON_MODAL,
		"Cancel",
		"Continue",
		// 'callBackResp' function:
		func(resp int) {
			switch resp {
			case 0:
				log.Printf("Canceled: %d\n", resp)
			case 1:
				log.Printf("Continue clicked: %d\n", resp)
			}
		})
}

/***************************
* DlgMessage implementation.
 ***************************/
var dialogType = map[string]gtk.MessageType{
	"info": gtk.MESSAGE_INFO, "inf": gtk.MESSAGE_INFO,
	"warning": gtk.MESSAGE_WARNING, "wrn": gtk.MESSAGE_WARNING,
	"question": gtk.MESSAGE_QUESTION, "qst": gtk.MESSAGE_QUESTION,
	"error": gtk.MESSAGE_ERROR, "err": gtk.MESSAGE_ERROR,
	"other": gtk.MESSAGE_OTHER, "oth": gtk.MESSAGE_OTHER,
	"infoWithMarkup": gtk.MESSAGE_INFO, "infMU": gtk.MESSAGE_INFO,
	"warningWithMarkup": gtk.MESSAGE_WARNING, "wrMUn": gtk.MESSAGE_WARNING,
	"questionWithMarkup": gtk.MESSAGE_QUESTION, "qstMU": gtk.MESSAGE_QUESTION,
	"errorWithMarkup": gtk.MESSAGE_ERROR, "errMU": gtk.MESSAGE_ERROR,
	"otherWithMarkup": gtk.MESSAGE_OTHER, "othMU": gtk.MESSAGE_OTHER,
}

type DlgOptions int

const (
	DLG_NON_MODAL DlgOptions = 1 << iota
	DLG_MARKUP
	DLG_BTN_RELIEF_NONE
	DLG_INFO
	DLG_WARNING
	DLG_QUESTION
	DLG_ERROR
	DLG_OTHER
)

type dlgOptions struct {
	DLG_NON_MODAL,
	DLG_MARKUP,
	DLG_BTN_RELIEF_NONE,
	DLG_INFO,
	DLG_WARNING,
	DLG_QUESTION,
	DLG_ERROR,
	DLG_OTHER DlgOptions
}

func DialogMessageOptNew() *dlgOptions {
	dlgo := new(dlgOptions)
	dlgo.DLG_NON_MODAL = DLG_NON_MODAL
	dlgo.DLG_MARKUP = DLG_MARKUP
	dlgo.DLG_BTN_RELIEF_NONE = DLG_BTN_RELIEF_NONE
	dlgo.DLG_INFO = DLG_INFO
	dlgo.DLG_WARNING = DLG_WARNING
	dlgo.DLG_QUESTION = DLG_QUESTION
	dlgo.DLG_ERROR = DLG_ERROR
	dlgo.DLG_OTHER = DLG_OTHER
	return dlgo
}

/*
Old parameters format:
func DialogMessage(
window *gtk.Window,
dlgType interface{},
title interface{},
text interface{},
iconFileName interface{},
parameters ...interface{}) (value int)*/

/*
func DialogMessage(window *gtk.Window, title, text string, parameters ...interface{}) (value int)
*/

// DlgMessage: Display message dialog with multiples buttons.
// return get <0 for cross closed or >-1 correspondig to buttons order representation.
// dlgType accepte: "info", "warning", "question", "error", "other", by adding
// 'DlgOptions' type flags. 'parameters' must be in order:
// (title, text, iconFileName, options, buttons..., callBackResp)
// 'iconFileName' can be a []byte or a string, "" or nil -> default image.
// 'buttons' must be comma separated buttons names: "Cancel", "Ok", "Exit"...
// 'callBackResp' is used for non-modal dialog message in form of: func(resp int), this
// function is required (even if it is not used), at the tail of the parameters if the
// 'DLG_NON_MODAL' flag is set, this function permit to retrieve the response value.
func DialogMessage(window *gtk.Window, parameters ...interface{}) (value int) {

	/********************************************************************************************************
	 *             -- When all sources codes update was done, function signature will become: --            *
	 * func DialogMessage(window *gtk.Window, title, text string, parameters ...interface{}) (value int) {  *
	 ********************************************************************************************************/

	var (
		response  int
		msgDialog *gtk.MessageDialog
		box       *gtk.Box
		err       error

		reliefBtnNone = false
		useMarkup,
		passed,
		messTypeIsSet bool
		modal       = true
		messageType gtk.MessageType
		dlgFlag     gtk.DialogFlags = gtk.DIALOG_MODAL

		// executed when a non modal dialog message has been created,
		// used to retrieve the response.
		callBackResp func(resp int) = nil

		title, text, dlgType string
		iconFileName         interface{}
		buttons              []string
	)

	// 'parameters' parsing
	for idx := 0; idx < len(parameters); idx++ {
		switch params := parameters[idx].(type) {
		/* will be remove after all sources code was updated */
		/*    This part is made for compatibility reasons    */
		case string:
			switch idx {
			case 0:
				for _, lt := range []string{
					"info",
					"warning",
					"question",
					"error",
					"other",
					"infoWithMarkup",
					"warningWithMarkup",
					"questionWithMarkup",
					"errorWithMarkup",
					"otherWithMarkup"} {
					if lt == params {

						log.Printf("[DialogMessage], This setting up options is outdated, please use %s instead\n", "'DlgOptions' flags")

						dlgType = lt
						title = parameters[idx+1].(string)
						text = parameters[idx+2].(string)

						switch ifn := parameters[idx+3].(type) {
						case nil, interface{}:
							iconFileName = ifn
						}

						prm, ok := parameters[idx+4].(DlgOptions)
						if ok {
							switch {
							case prm&DLG_INFO != 0:
								messTypeIsSet = true
								messageType = gtk.MESSAGE_INFO
							case prm&DLG_WARNING != 0:
								messTypeIsSet = true
								messageType = gtk.MESSAGE_WARNING
							case prm&DLG_QUESTION != 0:
								messTypeIsSet = true
								messageType = gtk.MESSAGE_QUESTION
							case prm&DLG_ERROR != 0:
								messTypeIsSet = true
								messageType = gtk.MESSAGE_ERROR
							case prm&DLG_OTHER != 0:
								messTypeIsSet = true
								messageType = gtk.MESSAGE_OTHER
							}
							if prm&DLG_BTN_RELIEF_NONE != 0 {
								reliefBtnNone = true
							}
							if prm&DLG_MARKUP != 0 {
								useMarkup = true
							}
							if prm&DLG_NON_MODAL != 0 {
								modal = false
								dlgFlag = gtk.DIALOG_DESTROY_WITH_PARENT
							}
						} else {
							idx--
						}
						idx = idx + 5
						// Retrieving buttons
						for idxSub := idx; idxSub < len(parameters); idxSub++ {
							if btnLbl, ok := parameters[idxSub].(string); ok {
								buttons = append(buttons, btnLbl)
								continue
							}
							idx = idxSub - 1
							break
						}
						passed = true
						break
					}
					if passed {
						break
					}
				}
				if passed {
					continue
				}

				title = params
				continue
			case 1:
				text = params
				continue
			case 4:
				for idxSub := idx; idxSub < len(parameters); idxSub++ {
					if btnLbl, ok := parameters[idxSub].(string); ok {
						buttons = append(buttons, btnLbl)
						continue
					}
					idx = idxSub - 1
					break
				}
				continue
			}
			/*^^^ will be remove after all sources code was updated ^^^*/

		//TODO Enable this when the new function signature become active.
		// case []string:
		// 	for _, btnLbl := range params {
		// 		buttons = append(buttons, btnLbl)
		// 	}
		case DlgOptions:
			switch {
			case params&DLG_INFO != 0:
				messTypeIsSet = true
				messageType = gtk.MESSAGE_INFO
			case params&DLG_WARNING != 0:
				messTypeIsSet = true
				messageType = gtk.MESSAGE_WARNING
			case params&DLG_QUESTION != 0:
				messTypeIsSet = true
				messageType = gtk.MESSAGE_QUESTION
			case params&DLG_ERROR != 0:
				messTypeIsSet = true
				messageType = gtk.MESSAGE_ERROR
			case params&DLG_OTHER != 0:
				messTypeIsSet = true
				messageType = gtk.MESSAGE_OTHER
			}
			if params&DLG_BTN_RELIEF_NONE != 0 {
				reliefBtnNone = true
			}
			if params&DLG_MARKUP != 0 {
				useMarkup = true
			}
			if params&DLG_NON_MODAL != 0 {
				modal = false
			}
		case func(resp int):
			callBackResp = params
		case interface{}, nil:
			iconFileName = params
		}
	}

	if !messTypeIsSet {
		log.Printf("[DialogMessage], This setting up options is outdated, please use %s instead\n", "'DlgOptions' flags")
		messageType = dialogType[dlgType]

		if strings.Contains(dlgType, "WithMarkup") {
			msgDialog.SetProperty("use-markup", true)
			log.Printf("[DialogMessage], This setting up options is outdated, please use %s instead\n", "'DlgOptions' flags")
		}
	}

	msgDialog = gtk.MessageDialogNew(window,
		dlgFlag,
		messageType,
		gtk.BUTTONS_NONE,
		"")

	// Modal ?
	if !modal {
		if callBackResp == nil {
			log.Fatalf("[DialogMessage], A callback function 'func(resp int)', is required at the tail of parameters if the 'DLG_NON_MODAL' flag is set.\n")
		}
		msgDialog.SetModal(modal)
		msgDialog.Connect("response", func(dlg *gtk.MessageDialog, resp int) {
			callBackResp(resp)
			dlg.Destroy()
		})
	}

	// Image
	switch iconFileName.(type) {
	case string:
		if len(iconFileName.(string)) != 0 {
			if box, err = msgDialog.GetContentArea(); err == nil {
				gipf.SetPict(box, iconFileName, 18)
			} else {
				log.Printf("[DialogMessage] Could not get content area: %v\n", err)
			}
		}
	}

	msgDialog.SetSkipTaskbarHint(true)
	msgDialog.SetKeepAbove(true)
	// Add button(s)
	for idx, btn := range buttons {
		button, err := msgDialog.AddButton(btn, gtk.ResponseType(idx))
		if err != nil {
			log.Printf("[DialogMessage] Button #%d could not be created: %v", idx, err)
		}
		buttonBox, err := button.GetParent()
		if err != nil {
			log.Printf("[DialogMessage] Button #%d, Unable to get parent: %v", idx, err)
		}
		buttonBox.(*gtk.ButtonBox).SetProperty("hexpand", false)
		buttonBox.(*gtk.ButtonBox).SetBorderWidth(0)
		buttonBox.(*gtk.ButtonBox).SetMarginTop(0)
		buttonBox.(*gtk.ButtonBox).SetMarginBottom(0)
		buttonBox.(*gtk.ButtonBox).SetMarginEnd(0)
		buttonBox.(*gtk.ButtonBox).SetMarginStart(0)

		if reliefBtnNone {
			button.SetProperty("relief", gtk.RELIEF_NONE)
		}
		button.SetHAlign(gtk.ALIGN_END)
		button.SetHExpand(false)
		button.SetBorderWidth(2)
		button.SetMarginTop(0)
		button.SetMarginBottom(0)
		button.SetMarginEnd(0)
		button.SetMarginStart(0)
	}

	if useMarkup {
		msgDialog.SetProperty("use-markup", true)
	}
	msgDialog.SetProperty("text", text)
	msgDialog.SetTitle(title)

	if modal {
		response = int(msgDialog.Run())
		msgDialog.Destroy()
	} else {
		msgDialog.ShowAll()
	}
	return response
}
