// progress.go

/*
	Source file auto-generated on Thu, 17 Oct 2019 02:33:04 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	This structure implement a progressbar.
*/

package gtk3Import

import (
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	gipf "github.com/hfmrow/gtk3Import/pixbuff"
)

type ProgressBarStruct struct {
	RefreshMs       uint
	GifImageName    string
	varPath         interface{}
	gifImage        *gtk.Image
	box             *gtk.Box
	boxPosition     int
	fMain, fEnd     func() error
	progressBar     *gtk.ProgressBar
	TimeOutContinue bool
	ticker          *time.Ticker
}

// ProgressBarNew: creat a new 'ProgressBarStruct'
func ProgressBarNew(progressBar *gtk.ProgressBar) (pbs *ProgressBarStruct) {
	pbs = new(ProgressBarStruct)
	pbs.RefreshMs = 100
	pbs.progressBar = progressBar
	return
}

// ProgressGifNew: create 'ProgressBarStruct' frome given information.
// Note: 'varPath' could be []byte, filename or *gtk.Image
func ProgressGifNew(varPath interface{}, box *gtk.Box, position int) (pbs *ProgressBarStruct) {
	pbs = new(ProgressBarStruct)
	pbs.box = box
	pbs.boxPosition = position
	pbs.varPath = varPath
	return
}

// Init: structure initialization with start and end callback functions.
func (pbs *ProgressBarStruct) Init(fMain, fEnd func() error) {
	pbs.fMain, pbs.fEnd = fMain, fEnd
}

// StartGif: launch progressbar as gif animated image.
func (pbs *ProgressBarStruct) StartGif() (err error) {

	switch i := pbs.varPath.(type) {
	case *gtk.Image:
		pbs.gifImage = i
	default:
		pbs.gifImage, err = gipf.GetAnimationImage(i)
		if err != nil {
			return err
		}
	}
	pbs.box.Add(pbs.gifImage)
	pbs.box.ReorderChild(pbs.gifImage, pbs.boxPosition)
	pbs.gifImage.SetHAlign(gtk.ALIGN_FILL)
	pbs.gifImage.SetHExpand(true)

	pbs.gifImage.Show()
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err = pbs.fMain()
	}()

	waitGroup.Wait()
	if err != nil {
		return
	}
	glib.IdleAdd(func() {
		// remove gif image from box.
		pbs.box.Remove(pbs.gifImage)
		err = pbs.fEnd()
	})
	return
}

// StartTicker: start gtk progressbar using golang ticker function.
func (pbs *ProgressBarStruct) StartTicker() (err error) {

	pbs.ticker = time.NewTicker(time.Millisecond * time.Duration(pbs.RefreshMs))
	go func() {
		for _ = range pbs.ticker.C {
			glib.IdleAdd(func() {
				pbs.progressBar.Pulse()
			})
		}
	}()

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err = pbs.fMain()
	}()

	waitGroup.Wait()
	if err != nil {
		return
	}
	glib.IdleAdd(func() {
		pbs.progressBar.SetFraction(0)
		if pbs.fEnd != nil {
			err = pbs.fEnd()
		}
		pbs.ticker.Stop()
	})

	return
}

// StartTimeOut: start gtk progressbar using 'glib.TimeoutAdd' command
func (pbs *ProgressBarStruct) StartTimeOut() (err error) {

	pbs.TimeOutContinue = true
	glib.TimeoutAdd(pbs.RefreshMs, func() bool {
		glib.IdleAdd(func() {
			pbs.progressBar.Pulse()
		})
		if pbs.TimeOutContinue {
			return true
		} else {
			glib.IdleAdd(func() {
				pbs.progressBar.SetFraction(0)
				if pbs.fEnd != nil {
					err = pbs.fEnd()
				}
			})
			return false
		}
	})

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		pbs.fMain()
	}()
	waitGroup.Wait()
	pbs.TimeOutContinue = false
	return
}
