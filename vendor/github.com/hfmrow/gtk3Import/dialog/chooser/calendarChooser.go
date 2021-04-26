// calendarChooser.go

/// +build ignore

/*
	Source file auto-generated on Sat, 19 Oct 2019 22:06:16 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019 H.F.M - calendar chooser v2.0
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	This structure and methods give a convenient way to display an use a calendar window to
	setting and retrieving calendar values.

Usage:
			var err error
			var ok bool

			calData := new(gidgcr.CalendarData)
			calData.FromLayout("2019.09.21.04.30.32")

			// []string{"Reset", "Ok"} means there is two button, "ok" will be emitted whether the button > 0 is clicked
			cal := gidgcr.CalendarNew(mainObjects.MainWindow, calData, "Test calendar", logout48, []string{"Reset", "Ok"}, crossIcon48, tickIcon48)
			cal.ButtonsSize = 24
			cal.DisplayTime = true

			ok, err = cal.Run()
			if err != nil {
				log.Fatalf("%s\n", err.Error())
			}

			if ok {
				fmt.Printf("%#v\n", cal.Result.ToLayout())
			}

	- The layout parameter describes the format of a time value.
	  It should be the magical reference date: "Mon Jan 2 15:04:05 MST 2006"
*/

package gtk3Import

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gipf "github.com/hfmrow/gtk3Import/pixbuff"
)

/*
*	Calendar implementation
 */

// This structure and methods give a convenient way to display and
// use a calendar window to setting and retrieving calendar values.
type Calendar struct {
	TitleWindow      string
	ImageTop         interface{}
	LabelHour        string
	ButtonsImages    []interface{}
	ButtonsLabels    []string
	ButtonsSize      int
	ButtonOkPosition int // Specify the position of the "validate" button (start at 0 from left)
	Width            int
	Height           int
	Resizable        bool
	KeepAbove        bool
	Result           *CalendarData
	DisplayTime      bool

	parentWindow        *gtk.Window
	dbs                 *gidg.DialogBoxStructure
	calendar            *gtk.Calendar
	spinH, spinM, spinS *gtk.SpinButton
}

// CalendarNew: create new structure and initialize it.
func CalendarNew(parentWindow *gtk.Window, calData *CalendarData, titleWindow string,
	topImage interface{}, buttons []string, imagesButtons ...interface{}) *Calendar {
	cal := new(Calendar)
	cal.Init(parentWindow, calData, titleWindow, topImage, buttons, imagesButtons...)
	return cal
}

// Init: Initialize structure
func (cal *Calendar) Init(parentWindow *gtk.Window, calData *CalendarData, titleWindow string,
	topImage interface{}, buttons []string, imagesButtons ...interface{}) {

	cal.parentWindow = parentWindow
	cal.TitleWindow = titleWindow
	cal.ImageTop = topImage
	cal.ButtonsLabels = buttons
	cal.ButtonsImages = imagesButtons

	cal.Result = calData
	cal.LabelHour = "Time"
	cal.Width = 100
	cal.Height = 100
	cal.ButtonsSize = 18
	cal.Width = 425
	cal.Resizable = false
}

// Show: Display about box
func (cal *Calendar) Run() (ok bool, err error) {
	if err = cal.build(); err == nil {
		cal.dbs.KeepAbove = cal.KeepAbove
		cal.dbs.WidgetExpend = true
		cal.dbs.SetSize(cal.Width, cal.Height)
		cal.dbs.Resizable = cal.Resizable
		cal.dbs.IconsSize = cal.ButtonsSize
		// Get cal values if ok button is pressed
		cal.dbs.ResponseCallBack = func(dlg *gtk.Dialog, resp int) {
			if resp == cal.ButtonOkPosition {
				cal.Result.ToCalData(cal)
				ok = true
			}
		}
		cal.dbs.Run()
	}
	return
}

// build:
func (cal *Calendar) build() (err error) {
	var boxHMS *gtk.Box
	var lblHour *gtk.Label
	var sep1, sep2, sep3 *gtk.Separator
	var pixbuf *gdk.Pixbuf
	var imageTop *gtk.Image

	// Create new dialogbox structure
	cal.dbs = gidg.DialogBoxNew(cal.parentWindow, nil, cal.TitleWindow, cal.ButtonsLabels...)
	cal.dbs.ButtonsImages = cal.ButtonsImages

	// Build widgets used to this Calendar window.
	pixbuf, _ = gipf.GetPixBuf(cal.ImageTop) // If an error occurs pixbuf will be nul and imageTop too. So, no image displayed.
	if imageTop, err = gtk.ImageNewFromPixbuf(pixbuf); err == nil {
		if sep1, err = gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); err == nil {
			if sep2, err = gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); err == nil {
				if sep3, err = gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL); err == nil {
					if cal.calendar, err = gtk.CalendarNew(); err == nil {
						// Validate for double-click on Day
						cal.calendar.Connect("day-selected-double-click", func() {
							cal.dbs.Dialog.Response(gtk.ResponseType(cal.ButtonOkPosition))
						})
						if boxHMS, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2); err == nil {
							if cal.spinH, err = gtk.SpinButtonNewWithRange(0, 23, 1); err == nil {
								if cal.spinM, err = gtk.SpinButtonNewWithRange(0, 59, 1); err == nil {
									if cal.spinS, err = gtk.SpinButtonNewWithRange(0, 59, 1); err == nil {
										if lblHour, err = gtk.LabelNew(cal.LabelHour); err == nil {
											// // Define default values for calendar
											if cal.Result == nil {
												cal.Result = CalendarDataNew()
											}
											// Set default values to calendar
											cal.Result.ToCalendar(cal)

											// Add some properties
											cal.dbs.WidgetsProps.AddProperty("margin-top", 1)
											cal.dbs.WidgetsProps.AddProperty("margin-bottom", 1)
											cal.dbs.WidgetsProps.AddProperty("margin-start", 2)
											cal.dbs.WidgetsProps.AddProperty("margin-end", 2)
											// cal.dbs.WidgetsProps.AddProperty("relief", gtk.RELIEF_NONE)

											// Add time controls if requested
											if cal.DisplayTime {
												boxHMS.PackStart(lblHour, true, true, 2)
												boxHMS.Add(cal.spinH)
												boxHMS.Add(cal.spinM)
												boxHMS.Add(cal.spinS)
											}
											// Add widgets to the [DialogBox] structure
											cal.dbs.Widgets = []gtk.IWidget{
												imageTop,
												sep1,
												cal.calendar,
												sep2,
												boxHMS,
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
	return
}

type CalendarData struct {
	Year,
	Month,
	Day,
	Hour,
	Min,
	Sec int
	BlankTime time.Time // Used when no time is defined

	layout string
}

// CalendarDataNew: create a new calendarData structure
func CalendarDataNew() (calData *CalendarData) {
	calData = new(CalendarData)
	calData.Init(time.Now().UTC())
	return
}

// Init: set values relatives "t" or time.Time{} if not specified.
func (calData *CalendarData) Init(t ...time.Time) {

	calData.BlankTime = time.Time{} //  Init to nul date. means no date time was specified.

	teaTime := calData.BlankTime
	if len(t) > 0 {
		teaTime = t[0]
	}

	calData.Year = teaTime.Year()
	calData.Month = int(teaTime.Month())
	calData.Day = teaTime.Day()
	calData.Hour = teaTime.Hour()
	calData.Min = teaTime.Minute()
	calData.Sec = teaTime.Second()
	calData.layout = "2006.01.02.15.04.05"

}

// calDataToCalendar:
func (calData *CalendarData) ToCalendar(cal *Calendar) {
	cal.calendar.SetProperty("year", uint(calData.Year))
	cal.calendar.SetProperty("month", uint(calData.Month-1))
	cal.calendar.SetProperty("day", uint(calData.Day))
	cal.spinH.SetValue(float64(calData.Hour))
	cal.spinM.SetValue(float64(calData.Min))
	cal.spinS.SetValue(float64(calData.Sec))
}

// calendarToCalData:
func (calData *CalendarData) ToCalData(cal *Calendar) {
	Y, M, D := cal.calendar.GetDate()
	calData.Year, calData.Month, calData.Day = int(Y), int(M)+1, int(D)
	calData.Hour, calData.Min, calData.Sec = cal.spinH.GetValueAsInt(), cal.spinM.GetValueAsInt(), cal.spinS.GetValueAsInt()
}

// ToTime:
func (calData *CalendarData) ToTime() (outTime time.Time) {
	var err error
	if outTime, err = time.Parse(calData.layout, calData.ToLayout()); err != nil {
		outTime = calData.BlankTime
	}
	return
}

// ToTimeAsUTC: convert given time to native UTC (based on local time zone),
// entered time H,M,S will be considered as UTC (useful for files time
// comparaison that always calculated using UTC).
func (calData *CalendarData) ToTimeAsUTC() (outTime time.Time) {
	var err error
	if outTime, err = time.Parse(calData.layout, calData.ToLayout()); err != nil {
		outTime = calData.BlankTime
	} else {
		_, offset := outTime.Local().Zone()
		outTime = time.Unix(outTime.Unix()-int64(offset), 0)
	}
	return
}

// FromTime
func (calData *CalendarData) FromTime(inTime time.Time) {
	calData.FromLayout(inTime.Format(calData.layout))
}

// ToLayout: fill missing char to complies with "2006.01.02.15.04.05",
// the magical reference format.
func (calData *CalendarData) ToLayout() (out string) {
	return fmt.Sprintf("%4d.%02d.%02d.%02d.%02d.%02d", calData.Year, calData.Month, calData.Day, calData.Hour, calData.Min, calData.Sec)
}

// FromLayout: build CalendarData using this layout: "2006.01.02.15.04.05"
func (calData *CalendarData) FromLayout(in string) (err error) {
	sl := strings.Split(in, ".")
	if calData.Year, err = strconv.Atoi(sl[0]); err == nil {
		if calData.Month, err = strconv.Atoi(sl[1]); err == nil {
			if calData.Day, err = strconv.Atoi(sl[2]); err == nil {
				if calData.Hour, err = strconv.Atoi(sl[3]); err == nil {
					if calData.Min, err = strconv.Atoi(sl[4]); err == nil {
						calData.Sec, err = strconv.Atoi(sl[5])
					}
				}
			}
		}
	}
	return
}
