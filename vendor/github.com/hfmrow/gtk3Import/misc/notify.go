// notify.go

/*
*	©2018 H.F.M. MIT license
*	Handle gtk3 Dialogs.
*
*	Message, Question/Response dialogs, file/dir dialogs, Notifications.
 */

package gtk3Import

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Notify: Display a notify message at the top right of screen.
func Notify(title, text string) {
	const appID = "h.f.m"
	app, _ := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	//Shows an application as soon as the app starts
	app.Connect("activate", func() {
		notif := glib.NotificationNew(title)
		notif.SetBody(text)
		app.SendNotification(appID, notif)
	})
	app.Run(nil)
	app.Quit()
}
