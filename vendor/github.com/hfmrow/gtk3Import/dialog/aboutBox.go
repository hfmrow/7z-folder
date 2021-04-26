// aboutBox.go

/*
	Source file auto-generated on Sat, 19 Oct 2019 22:06:16 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 H.F.M - about box v2.0

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	glsg "github.com/hfmrow/genLib/strings"

	gipops "github.com/hfmrow/gtk3Import/pango/pangoSimple"
	gipf "github.com/hfmrow/gtk3Import/pixbuff"
)

/*
*	Aboutbox implementation
 */

type AboutInfos struct {
	TitleWindow,
	AppName,
	AppVers,
	AppCreats,
	YearCreat,
	LicenseAbrv,
	LicenseShort,
	Repository,
	CssName, // Used to access dialog fromm CSS
	Description string

	ImageTop,
	ImageOkButton interface{}

	ImageTopHeight,
	MaxCharWidthText,
	ImageOkButtonSize,
	Width,
	HttpKeepFromEnd, // how many elements from end of http adress will be displayed in description and licence text.
	Height int

	Resizable,
	KeepAbove bool

	parentWindow *gtk.Window
	DlgBoxStruct *DialogBoxStructure
}

func AboutInfosNew(parentWindow *gtk.Window) *AboutInfos {

	ab := new(AboutInfos)
	ab.parentWindow = parentWindow
	ab.MaxCharWidthText = 80
	ab.ImageTopHeight = 30
	ab.CssName = "AboutBox"
	ab.HttpKeepFromEnd = 2
	ab.Width = 425
	ab.Height = 100
	ab.ImageOkButtonSize = 24
	ab.Resizable = false

	return ab
}

// InitFillInfos: Initialize structure
func (ab *AboutInfos) FillInfos(titleWindow, appName, appVers, appCreat,
	yearCreat, licenseAbrv, licenseShort, repository, description string, topImage, okBtnIcon interface{}) {

	ab.TitleWindow = titleWindow
	ab.AppName = appName
	ab.AppVers = appVers
	ab.AppCreats = appCreat
	ab.YearCreat = "©" + yearCreat
	ab.LicenseAbrv = licenseAbrv
	ab.LicenseShort = licenseShort
	ab.Repository = repository
	ab.Description = glsg.FormatText(description, ab.MaxCharWidthText, true)
	ab.ImageTop = topImage
	ab.ImageOkButton = okBtnIcon
}

// InitFillInfos: Initialize structure
func (ab *AboutInfos) InitFillInfos(parentWindow *gtk.Window, titleWindow, appName, appVers, appCreat,
	yearCreat, licenseAbrv, licenseShort, repository, description string, topImage, okBtnIcon interface{}) {

	log.Println("WARNING: method out of date. Use 'AboutInfosNew', set wanted option inside structure,\nfollowed by 'FillInfos' with application informations.")

	if ab.MaxCharWidthText == 0 {
		ab.MaxCharWidthText = 80
	}

	ab.parentWindow = parentWindow
	ab.TitleWindow = titleWindow
	ab.AppName = appName
	ab.AppVers = appVers
	ab.AppCreats = appCreat
	ab.YearCreat = "©" + yearCreat
	ab.LicenseAbrv = licenseAbrv
	ab.LicenseShort = licenseShort
	ab.Repository = repository
	ab.Description = glsg.FormatText(description, ab.MaxCharWidthText, true)
	ab.ImageTop = topImage
	ab.ImageOkButton = okBtnIcon
}

// Show: Display about box
func (ab *AboutInfos) Show() (err error) {
	if err = ab.build(); err == nil {

		ab.DlgBoxStruct.CssName = ab.CssName
		ab.DlgBoxStruct.KeepAbove = ab.KeepAbove
		ab.DlgBoxStruct.SetSize(ab.Width, ab.Height)
		ab.DlgBoxStruct.Resizable = ab.Resizable
		ab.DlgBoxStruct.IconsSize = ab.ImageOkButtonSize
		ab.DlgBoxStruct.Run()
	}
	return
}

// build:
func (ab *AboutInfos) build() (err error) {
	var labelAppName, labelVersion, labelYearCreator, labelDescriptionLbl,
		labelDescription, labelRepolinkLbl, labelRepolink, labelLicense *gtk.Label
	var sep1, sep2, sep3 *gtk.Separator
	var box *gtk.Box
	var pixbuf *gdk.Pixbuf
	var imageTop *gtk.Image

	ps := gipops.PangoSimpleNew()
	pc := gipops.PangoColorNew()

	// Create new dialogbox structure
	ab.DlgBoxStruct = DialogBoxNew(ab.parentWindow, box, ab.TitleWindow, "Ok")
	ab.DlgBoxStruct.ButtonsImages = []interface{}{ab.ImageOkButton}

	// Add markup for some elements
	name, repo, lic := doMarkup(ab.AppName, ab.Repository, ab.LicenseShort, ab.HttpKeepFromEnd)

	// Build widgets used to this About box window
	pixbuf, _ = gipf.GetPixBuf(ab.ImageTop, ab.ImageTopHeight)
	// If an error occurs pixbuf will be nul and imageTop too. So, no image displayed.
	if imageTop, err = gtk.ImageNewFromPixbuf(pixbuf); err == nil {
		if labelAppName, err = gtk.LabelNew(""); err == nil {
			labelAppName.SetMarkup("\n" + name)

			if labelVersion, err = gtk.LabelNew(ab.AppVers + "\n"); err == nil {
				if labelYearCreator, err = gtk.LabelNew(ab.YearCreat + " " + ab.AppCreats + "\n"); err == nil {

					if labelDescriptionLbl, err = gtk.LabelNew(""); err == nil {
						labelDescriptionLbl.SetMarkup("\n<b>" + ps.ApplyMarkup(`<span foreground="`+pc.Brown+`DD">`, `Description:`) + "</b>")

						if labelDescription, err = gtk.LabelNew(""); err == nil {
							labelDescription.SetMarkup(ps.MarkupHttpClickable(ab.Description+"\n", ab.HttpKeepFromEnd))

							if labelRepolinkLbl, err = gtk.LabelNew(""); err == nil {
								labelRepolinkLbl.SetMarkup("\n<b>" + ps.ApplyMarkup(`<span foreground="`+pc.Brown+`DD">`, `Source repository:`) + "</b>")

								if labelRepolink, err = gtk.LabelNew(""); err == nil {
									labelRepolink.SetMarkup(repo + "\n")

									if labelLicense, err = gtk.LabelNew(""); err == nil {
										labelLicense.SetMarkup("\n" + lic + "\n")

										if sep1, err = gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); err == nil {
											if sep2, err = gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); err == nil {
												if sep3, err = gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); err == nil {

													// Add some properties
													ab.DlgBoxStruct.WidgetsProps.AddProperty("margin-top", 2)
													ab.DlgBoxStruct.WidgetsProps.AddProperty("justify", gtk.JUSTIFY_CENTER)
													ab.DlgBoxStruct.WidgetsProps.AddProperty("wrap", true)

													// Add widgets to the [DialogBox] structure
													ab.DlgBoxStruct.Widgets = []gtk.IWidget{
														imageTop,
														labelAppName,
														labelVersion,
														labelYearCreator,
														sep1,
														labelDescriptionLbl,
														labelDescription,
														labelRepolinkLbl,
														labelRepolink,
														sep2,
														labelLicense,
														sep3}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func doMarkup(appName, repo, licence string, keepFromEnd ...int) (outAppName, outRepo, outLicense string) {
	kfe := 2
	if len(keepFromEnd) > 0 {
		kfe = keepFromEnd[0]
	}
	ps := gipops.PangoSimpleNew()
	pc := gipops.PangoColorNew()

	outAppName = ps.ApplyMarkup(`<b>`, appName)
	outAppName = ps.ApplyMarkupAgain(`<span font_size="xx-large">`, outAppName)
	outAppName = ps.ApplyMarkupAgain(`<span foreground="`+pc.Brown+`">`, outAppName)

	outRepo = ps.ApplyMarkup(`<a href="`+"https://"+repo+`">`, repo)

	outLicense = ps.MarkupHttpClickable(licence, kfe)
	return
}
