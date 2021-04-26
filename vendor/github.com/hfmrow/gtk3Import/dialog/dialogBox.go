// dialogBox.go

/*
  Create a Dialog, who accept GtkWidgets defined into structure.
  The Structure contain all needed options to fill most of usages,
  including scrolling capabilities.

Usage:
	if tw, err := gtk.TreeViewNew(); err == nil {
		dbs := gi.DialogBoxNew(MainWindow, gtk.DIALOG_DESTROY_WITH_PARENT, tw, "test Title", "No", "Yes", "why")
		dbs.ButtonsImages = dbs.ValuesToInterfaceSlice("assets/images/Sign-cancel-20.png", "", signSelect20) // "signSelect20" is []byte of image
		dbs.ScrolledArea = true
		result=dbs.Run()
		// Do what you want with "result"
	}

	- A full examples:
		- In file "aboutBox.go", that use the majority of the functionalities allowed by the code below.
		- Another full example in file "calendar.go"
*/

package gtk3Import

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	gipf "github.com/hfmrow/gtk3Import/pixbuff"
	gitsws "github.com/hfmrow/gtk3Import/tools/widgets"
)

// DialogBoxStructure: Wrap a Dialog with desired count of
// buttons and widgets. The structure have defaults parameters
type DialogBoxStructure struct {
	BoxHAlign,
	BoxVAlign gtk.Align

	BoxOrientation gtk.Orientation

	BoxHExpand,
	BoxVExpand,
	WidgetExpend,
	SkipTaskbarHint,
	KeepAbove,
	Resizable,
	ScrolledArea,
	WidgetFill,
	MarkupLabel,
	LabelLineWrap bool
	Padding uint

	width,
	height,
	posX,
	posY,
	Response,
	IconsSize int

	Title,
	Text,
	CssName string

	Buttons           []string
	ButtonsImages     []interface{} // image representation from file or []byte, depending on type
	Dialog            *gtk.Dialog
	ButtonReliefStyle gtk.ReliefStyle
	Widgets           []gtk.IWidget
	WidgetsProps      gitsws.WidgetProperties // each property applyed to each object
	// Widgets below will be added after scrolled window without any property,
	// they must be applyed mlanually on creation by the caller (useful to add
	// some checkboxes or others after a TreeView).
	WidgetsOutOfScrolledArea []gtk.IWidget

	// Function used when a button is clicked or window is closed
	// This public function is wrapped with "internalResponseCallBack"
	ResponseCallBack func(dlg *gtk.Dialog, response int)

	dialogFlag gtk.DialogFlags

	modal,
	buttonsWithImages bool

	label          *gtk.Label
	box            *gtk.Box
	scrolledWindow *gtk.ScrolledWindow
	window         *gtk.Window

	internalResponseCallBack func(dlg *gtk.Dialog, response int)
}

// DialogBoxNew: Create a new structure to wrap a GtkDialog including
// defaults parameters. "widget" can be "nul" to only display "text".
func DialogBoxNew(window *gtk.Window, widget gtk.IWidget, title string, buttons ...string) *DialogBoxStructure {
	dbs := new(DialogBoxStructure)
	dbs.DialogBoxInit(window, widget, title, buttons...)
	return dbs
}

func (dbs *DialogBoxStructure) DialogBoxInit(window *gtk.Window, widget gtk.IWidget, title string, buttons ...string) {
	dbs.window = window
	dbs.ButtonReliefStyle = gtk.RELIEF_NONE
	dbs.BoxHAlign = gtk.ALIGN_FILL
	dbs.BoxVAlign = gtk.ALIGN_FILL
	dbs.BoxOrientation = gtk.ORIENTATION_VERTICAL
	dbs.dialogFlag = gtk.DIALOG_DESTROY_WITH_PARENT
	dbs.width, dbs.height,
		dbs.posX, dbs.posY = 640, 480, -1, -1
	dbs.IconsSize = 18
	dbs.BoxHExpand = true
	dbs.BoxVExpand = true
	dbs.WidgetExpend = true
	dbs.WidgetFill = true
	dbs.SkipTaskbarHint = true
	dbs.KeepAbove = true
	dbs.Resizable = true
	dbs.Title = title
	dbs.LabelLineWrap = true
	dbs.Padding = 0
	dbs.CssName = "CustomDialog"
	dbs.ResponseCallBack = func(dlg *gtk.Dialog, response int) { // Init default callback
	}
	dbs.internalResponseCallBack = func(dlg *gtk.Dialog, response int) { // Init internal and wrap default callback
		dbs.Response = response
		dbs.ResponseCallBack(dlg, response)
	}

	if widget != nil {
		dbs.Widgets = append([]gtk.IWidget{widget}, dbs.Widgets...) // Prepend
	}
	if len(buttons) == 0 {
		dbs.Buttons = []string{"Ok"}
	} else {
		dbs.Buttons = buttons
	}
}

// Run: function return "value" < 0 for cross closed or >= 0
// corresponding to buttons order representation starting with 0 at left.
// The dialog is destroyed at the end of process.
// Only the Modal Dialog window is allowed with this function.
// Use RunForResults() to allow non Modal usage.
func (dbs *DialogBoxStructure) Run() (value int) {
	dbs.modal = true
	dbs.buildDialog()
	value = int(dbs.Dialog.Run())
	dbs.Dialog.Destroy()
	return
}

// RunNonModal: allow non modal Dialog "respCallBack()" is executed on button click or exit window
func (dbs *DialogBoxStructure) RunNonModal(respCallBack ...func(dlg *gtk.Dialog, response int)) {
	dbs.modal = false
	if len(respCallBack) > 0 {
		dbs.ResponseCallBack = respCallBack[0]
	}
	dbs.internalResponseCallBack = func(dlg *gtk.Dialog, response int) { // Default callback function
		dbs.Response = response
		dbs.ResponseCallBack(dlg, response)
		dbs.Dialog.Destroy()
	}
	dbs.buildDialog()
}

// SetSize:
func (dbs *DialogBoxStructure) SetSize(width, height int) {
	dbs.width = width
	dbs.height = height
}

// SetPosition:
func (dbs *DialogBoxStructure) SetPosition(posX, posY int) {
	dbs.posX = posX
	dbs.posY = posY
}

// GetSize:
func (dbs *DialogBoxStructure) GetSize() (width, height int) {
	return dbs.Dialog.GetSize()
}

// GetPosition:
func (dbs *DialogBoxStructure) GetPosition() (posX, posY int) {
	return dbs.Dialog.GetPosition()
}

// BringToFront: Set window position to be over all others windows
// without staying on top whether another window come to be selected.
func (dbs *DialogBoxStructure) BringToFront() {
	dbs.Dialog.Deiconify()
	dbs.Dialog.ShowAll()
	dbs.Dialog.GrabFocus()
}

// buildDialog: Create the dialog window with defined parameters.
func (dbs *DialogBoxStructure) buildDialog() (err error) {
	var btnObj *gtk.Button

	// Build Dialog
	if dbs.modal {
		dbs.dialogFlag = gtk.DIALOG_MODAL
	}
	if dbs.Dialog, err = gtk.DialogNewWithButtons(dbs.Title, dbs.window, dbs.dialogFlag); err == nil {
		// Dialog options
		dbs.Dialog.SetDefaultSize(dbs.width, dbs.height)

		if dbs.posX != -1 && dbs.posY != -1 {
			dbs.Dialog.Move(dbs.posX, dbs.posY)
		}

		dbs.Dialog.SetName(dbs.CssName)
		dbs.Dialog.SetSkipTaskbarHint(dbs.SkipTaskbarHint)
		dbs.Dialog.SetKeepAbove(dbs.KeepAbove)
		dbs.Dialog.SetResizable(dbs.Resizable)
		dbs.Dialog.SetModal(dbs.modal)
		dbs.Dialog.Connect("response", dbs.internalResponseCallBack) // Default callback function

		if dbs.box, err = dbs.Dialog.GetContentArea(); err == nil {
			if dbs.label, err = gtk.LabelNew(""); err == nil {
				// Markup & Label options
				dbs.label.SetSizeRequest(dbs.box.GetSizeRequest()) // Size, same as parent
				dbs.label.SetLineWrap(dbs.LabelLineWrap)
				dbs.label.SetLineWrapMode(pango.WRAP_WORD)
				if dbs.MarkupLabel {
					dbs.label.SetLabel(dbs.Text)
				} else {
					dbs.label.SetMarkup(dbs.Text)
				}
			}
		}
	}

	if err != nil {
		return
	}

	// Control
	dbs.buttonsWithImages = len(dbs.ButtonsImages) != 0
	if len(dbs.Buttons) != len(dbs.ButtonsImages) && dbs.buttonsWithImages {
		log.Fatalf("You must provide an image or an empty string for each button\nButton(s) count: %d, Image(s) count: %d\n",
			len(dbs.Buttons), len(dbs.ButtonsImages))
	}

	// Box options
	dbs.box.SetHAlign(dbs.BoxHAlign)
	dbs.box.SetVAlign(dbs.BoxVAlign)
	dbs.box.SetOrientation(dbs.BoxOrientation)
	dbs.box.SetHExpand(dbs.BoxHExpand)
	dbs.box.SetVExpand(dbs.BoxVExpand)

	// Buttons
	for idxBtn, btnLbl := range dbs.Buttons {
		if btnObj, err = dbs.Dialog.AddButton(btnLbl, gtk.ResponseType(idxBtn)); err == nil {
			if dbs.buttonsWithImages {
				gipf.SetPict(btnObj, dbs.ButtonsImages[idxBtn], dbs.IconsSize)
				btnObj.SetRelief(dbs.ButtonReliefStyle)
			}
		}
	}

	// Packing
	if len(dbs.Text) > 0 {
		dbs.box.PackStart(dbs.label, dbs.WidgetExpend, dbs.WidgetFill, dbs.Padding)
	}
	if len(dbs.Widgets) > 0 {
		if dbs.ScrolledArea {
			if dbs.scrolledWindow, err = gtk.ScrolledWindowNew(nil, nil); err == nil {
				dbs.box.PackStart(dbs.scrolledWindow, dbs.WidgetExpend, dbs.WidgetFill, dbs.Padding) // Add ScrolledWindow to box
			} else {
				return
			}
		}

		for _, wdg := range dbs.Widgets {
			if wdg != nil {
				dbs.WidgetsProps.PropsToWidget(wdg) // Set properties to widget

				if dbs.ScrolledArea {
					dbs.scrolledWindow.Add(wdg) // Add objects to ScrolledWindow
				} else {
					dbs.box.PackStart(wdg, dbs.WidgetExpend, dbs.WidgetFill, dbs.Padding) // Add objects to box
				}
			}
		}
		if len(dbs.WidgetsOutOfScrolledArea) > 0 {
			for _, wdg := range dbs.WidgetsOutOfScrolledArea {
				dbs.box.PackEnd(wdg, false, false, dbs.Padding) // Add objects to box
			}
		}
	}
	// The show must go on
	dbs.Dialog.ShowAll()
	return
}
