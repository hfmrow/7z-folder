// cssWdgScnLoad.go

/*
	Load or read from data css style for an object(widget) or for entire screen.
*/

package gtk3Import

import (
	"io/ioutil"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// CssWidgetLoad: Load or read from data and apply css to
// widget if it's provided. Apply to screen otherwise.
func CssWdgScnLoad(filename string, wdgt ...*gtk.Widget) (err error) {

	data := []byte(filename)

	if data, err = ioutil.ReadFile(filename); err == nil {
		err = CssWdgScnBytes(data, wdgt...)
	}
	return
}

func CssWdgScnBytes(data []byte, wdgt ...*gtk.Widget) (err error) {

	var cssProv *gtk.CssProvider

	if cssProv, err = gtk.CssProviderNew(); err == nil {

		if err = cssProv.LoadFromData(string(data)); err == nil {
			if len(wdgt) == 0 {
				var screen *gdk.Screen
				if screen, err = gdk.ScreenGetDefault(); err == nil {
					gtk.AddProviderForScreen(screen, cssProv, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
				}
			} else {
				var styleContext *gtk.StyleContext
				if styleContext, err = wdgt[0].GetStyleContext(); err == nil {
					styleContext.AddProvider(cssProv, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
				}
			}
		}
	}
	return
}
