// dnd.go

// Source file auto-generated on Tue, 23 Jul 2019 04:14:20 using Gotk3ObjHandler v1.3.6 ©2019 H.F.M

/*
	Drag & drop handling .

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"bytes"
	"log"
	"net/url"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type DragNDropStruct struct {
	// gtkObject that receive DND
	Object interface{}
	// Contain the files list received
	FilesList *[]string
	// Callback called after data was received.
	callBackRecieveDone func()
	// Callback called during data reception. If false is returned, the loop ends
	callBackOnRecieve func(item interface{}, context *gdk.DragContext) bool
	// To build dnd context
	targets []gtk.TargetEntry
}

// DragNDropNew: configure controls to receive dndcontent. "filesList" can be "nil"
func DragNDropNew(objects interface{}, filesList *[]string,
	callBackRecieveDone func(),
	callBackOnRecieve ...func(item interface{}, context *gdk.DragContext) bool) *DragNDropStruct {

	ds := new(DragNDropStruct)
	ds.Object = objects

	ds.callBackOnRecieve = nil
	if len(callBackOnRecieve) > 0 {
		ds.callBackOnRecieve = callBackOnRecieve[0]
	}
	ds.callBackRecieveDone = callBackRecieveDone

	switch filesList {
	case nil:
		ds.FilesList = new([]string)
	default:
		ds.FilesList = filesList
	}
	ds.init()
	return ds
}

// Dispatching reciever object type (TreeView, Button ...)
func (ds *DragNDropStruct) init() {

	// Build DnD context
	targetTypes := []string{
		"x-special/mate-icon-list",
		"text/uri-list",
		"UTF8_STRING",
		"COMPOUND_TEXT",
		"TEXT",
		"STRING",
		"text/plain;charset=utf-8",
		"text/plain"}

	for _, tType := range targetTypes {
		te, err := gtk.TargetEntryNew(tType, gtk.TARGET_OTHER_APP, 0)
		if err != nil {
			log.Fatal(err)
		}
		ds.targets = append(ds.targets, *te)
	}

	switch ds.Object.(type) {

	case *gtk.Window:
		ds.Object.(*gtk.Window).DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			ds.targets,
			gdk.ACTION_COPY)
		ds.Object.(*gtk.Window).Connect("drag-data-received", ds.dndFilesReceived)

	case *gtk.Entry:
		ds.Object.(*gtk.Entry).DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			ds.targets,
			gdk.ACTION_COPY)
		ds.Object.(*gtk.Entry).Connect("drag-data-received", ds.dndFilesReceived)

	case *gtk.EventBox:
		ds.Object.(*gtk.EventBox).DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			ds.targets,
			gdk.ACTION_COPY)
		ds.Object.(*gtk.EventBox).Connect("drag-data-received", ds.dndFilesReceived)

	case *gtk.TreeView:
		ds.Object.(*gtk.TreeView).DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			ds.targets,
			gdk.ACTION_COPY)
		ds.Object.(*gtk.TreeView).Connect("drag-data-received", ds.dndFilesReceived)

	case *gtk.TextView:
		ds.Object.(*gtk.TextView).DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			ds.targets,
			gdk.ACTION_COPY)
		ds.Object.(*gtk.TextView).Connect("drag-data-received", ds.dndFilesReceived)

	case *gtk.Button:
		ds.Object.(*gtk.Button).DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			ds.targets,
			gdk.ACTION_COPY)
		ds.Object.(*gtk.Button).Connect("drag-data-received", ds.dndFilesReceived)

	case *gtk.Image:
		ds.Object.(*gtk.Image).DragDestSet(
			gtk.DEST_DEFAULT_ALL,
			ds.targets,
			gdk.ACTION_COPY)
		ds.Object.(*gtk.Image).Connect("drag-data-received", ds.dndFilesReceived)

		// Already handled natively !!!
		// case *gtk.FileChooserButton:
		// 	ds.Object.(*gtk.FileChooserButton).DragDestSet(
		// 		gtk.DEST_DEFAULT_ALL,
		// 		ds.targets,
		// 		gdk.ACTION_COPY)
		// 	ds.Object.(*gtk.FileChooserButton).Connect("drag-data-received", ds.dndFilesReceived)
	}
}

// ButtonInFilesReceived: Store in files list
func (ds *DragNDropStruct) dndFilesReceived(object interface{}, context *gdk.DragContext, x, y int, selData *gtk.SelectionData, info, time uint) {

	*ds.FilesList = (*ds.FilesList)[:0]
	data := selData.GetData()
	list := strings.Split(string(data), getTextEOL(data))

	for _, item := range list {
		if len(item) != 0 {

			if ds.callBackOnRecieve != nil {
				// For other type than string, callback permit to handle anything
				if !ds.callBackOnRecieve(item, context) {
					break
				}
			} else {
				// Default handling as a string
				if u, err := url.PathUnescape(item); err == nil {

					*ds.FilesList = append(*ds.FilesList, strings.TrimPrefix(u, "file://"))
				}
			}
		}
	}
	if ds.callBackRecieveDone != nil {
		ds.callBackRecieveDone()
	}
}

// GetTextEOL: Get EOL from text bytes (CR, LF, CRLF)
func getTextEOL(inTextBytes []byte) (outString string) {

	var (
		bCR   = []byte{0x0D}
		bLF   = []byte{0x0A}
		bCRLF = []byte{0x0D, 0x0A}
	)

	switch {
	case bytes.Contains(inTextBytes, bCRLF):
		outString = string(bCRLF)
	case bytes.Contains(inTextBytes, bCR):
		outString = string(bCR)
	default:
		outString = string(bLF)
	}
	return
}
