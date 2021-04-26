// colorChooser.go

package gtk3Import

import (
	"fmt"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	gipfmcrr "github.com/hfmrow/gtk3Import/pixbuff/misc/RGBA"
)

/************************************
*  ColorChooserDialog implementation.
 ************************************/

// ColorChooserDialogAndGetHexValue: Open color chooser dialog and retrieve value. String return is
// converted like this: "rgb(239,41,41)" to "#EF2929".
func ColorChooserDialogAndGetFloat64Val(parentWindow *gtk.Window, title string, orgCol []float64) (outFloat64 []float64, outStr string) {

	// Build Color chooser dialog en give it some basic parameters.
	if ColorChooserDialog, err := gtk.ColorChooserDialogNew(title, parentWindow); err == nil {
		if len(orgCol) != 0 {

			RGBA := gdk.NewRGBA(orgCol...)
			fmt.Println(RGBA.String())
			ColorChooserDialog.ColorChooser.SetRGBA(RGBA)
		}
		ColorChooserDialog.SetSkipTaskbarHint(true)
		ColorChooserDialog.SetKeepAbove(true)
		ColorChooserDialog.SetSizeRequest(10, 10)
		ColorChooserDialog.SetResizable(false)
		ColorChooserDialog.SetModal(true)
		ColorChooserDialog.ColorChooser.SetUseAlpha(true)

		switch ColorChooserDialog.Run() {
		case gtk.RESPONSE_OK:

			ColorChooserDialog.Close()
			return ColorChooserDialog.ColorChooser.GetRGBA().Floats(),
				gipfmcrr.ColGdkRGBA2Hex(ColorChooserDialog.ColorChooser.GetRGBA())

		case gtk.RESPONSE_CANCEL:

			ColorChooserDialog.Close()
		}
	}
	return outFloat64, outStr
}
