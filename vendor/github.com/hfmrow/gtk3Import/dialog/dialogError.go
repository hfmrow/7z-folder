// dialogError.go

package gtk3Import

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"

	gler "github.com/hfmrow/genLib/tools/errors"
)

// DialogError: Display error messaged dialog returning true in error case.
// options: devMode, forceExit, markupEnabled bool
// NOTICE: exit option must not be used if a "defer" function is initiated !
func DialogError(window *gtk.Window, title, text string, err error, options ...bool) bool {

	var (
		devMode              = true
		forceExit            = false
		DMOptions DlgOptions = DLG_ERROR
	)

	switch {
	case len(options) == 1:
		devMode = options[0]

	case len(options) == 2:
		devMode = options[0]
		forceExit = options[1]

	case len(options) == 3:
		devMode = options[0]
		forceExit = options[1]
		if options[2] {
			DMOptions = DMOptions | DLG_MARKUP
		}
	case len(options) > 3:
		fmt.Printf("DialogError: Bad, number of arguments, %v\n", options)
		os.Exit(1)
	}

	if gler.Check(err) {
		dispText := fmt.Sprintf("\n\n"+text+":\n%s", err.Error())
		if len(text) == 0 {
			dispText = fmt.Sprintf("\n\n" + err.Error() + "\n")
		}
		if devMode {
			if DialogMessage(
				window,
				title,
				dispText,
				nil, DMOptions,
				"Stop",
				"Continue") == 0 {
				os.Exit(1)
			}
		} else {
			DialogMessage(
				window,
				title,
				dispText,
				nil, DMOptions,
				"Ok")
			if forceExit {
				os.Exit(1)
			}
		}
		return true
	}
	return false
}
