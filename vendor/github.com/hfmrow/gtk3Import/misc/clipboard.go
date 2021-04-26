// clipboard.go

/*
*	©2019 H.F.M. MIT license
*	This is a simple clipboard handler to use with gotk3 "https://github.com/gotk3/gotk3"
 */

package gtk3Import

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type Clipboard struct {
	Entity *gtk.Clipboard
}

// ClipboardNew: Create new clipboard structure
func ClipboardNew() (c *Clipboard, err error) {

	c = new(Clipboard)
	err = c.Init()
	return
}

// Init: Initialise clipboard
func (c *Clipboard) Init() (err error) {

	c.Entity, err = gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	return
}

// GetText: Get text from clipboard
func (c *Clipboard) GetText() (clipboardContent string, err error) {
	return c.Entity.WaitForText()
}

// SetText: Set text to clipboard
func (c *Clipboard) SetText(clipboardContent string) {
	c.Entity.SetText(clipboardContent)
}

// Store: the current clipboard data somewhere so that it will stay around
// after the application has quit.
func (c *Clipboard) Store() {
	c.Entity.Store()
}
