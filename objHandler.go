// objHandler.go

// Source file auto-generated on Mon, 04 Mar 2019 17:22:55 using Gotk3ObjHandler v1.0 ©2019 H.F.M

/*
	7zFolder v1.0 ©2019 H.F.M

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/gotk3/gotk3/gtk"

	"github.com/gotk3/gotk3/glib"
)

// ButtonCompressClicked:
func ButtonCompressClicked() {
	var (
		filename string
		args     []string
		// tmpListFile     []string
		ok, askFilename bool
		err             error
		size            int64
		fi              os.FileInfo
	)
	if len(fos.Files) > 0 {
		// filename = fos.Files[0].FilenameFull
		filename = baseDir
		// The case where we have a directory name with dot
		EE_opt := EE_DEFAULT
		if fi, err = os.Stat(filename); err == nil {
			if fi.IsDir() {
				EE_opt = EE_PRESERVE_FILENAME
			}
			filename = ExtEnsure(filename, ".7z", EE_opt)

			if obj.CheckbuttonAutoInc.GetActive() {
				path := filepath.Dir(filename)
				// path := baseDir
				fosFilename, err := FilesOpStructNew()
				if checkErr(err) {
					return
				}
				fosFilename.ReadDirDepth = 0
				err = fosFilename.GetFilesDetailsReadDir(path)
				if checkErr(err) {
					return
				}
				tmpListFile := fosFilename.FilesGetAsRaw()
				tmpList := sortNumbered(filepath.Base(filename), tmpListFile)
				computedFN := computeFilename(tmpList[len(tmpList)-1])
				filename = filepath.Join(path, computedFN)

				if DialogMessage(obj.MainWindow, sts["confirm"],
					"\n\n"+sts["useFName"]+nl+TruncatePath(ExtEnsure(filename, ".7z"), 2),
					nil, DLG_OPT.DLG_INFO|DLG_OPT.DLG_BTN_RELIEF_NONE, sts["no"], sts["yes"]) != 1 {
					askFilename = true
				}
			}

			// if tmpListFile, size, ok, err = checkFiles(); checkErr(err) {
			if _, size, ok, err = checkFiles(); checkErr(err) {
				return
			}
			if !ok {
				return
			}
			ok = true
			if askFilename {
				FC_opt := fcOptions.FC_DLG_DEFAULT
				if obj.RadioButtonNew.GetActive() {
					FC_opt = fcOptions.FC_NOT_ASK_OVERWRITE
				}
				if filename, ok, err = FileChooser(
					obj.MainWindow, sts["chse7zFname"],
					ExtEnsure(filename, ".7z"),
					fcOptions.FC_DLG_SAVE|FC_opt); !ok {
					inErr := fmt.Errorf("Filename selection error.\n")
					if err != nil {
						inErr = fmt.Errorf("%v%v\n", inErr, err)
					}
					if checkErr(inErr) {
						return
					}
				}
			}
			if ok && err == nil {
				glib.IdleAdd(func() {
					if obj.RadioButtonNew.GetActive() {
						if _, err = os.Stat(filename); err == nil {
							if err = os.Remove(filename); checkErr(err) {
								return
							}
							obj.LabelOverwritten.SetLabel(sts["archOvrWrt"])
						}
						err = nil
					} else {
						if obj.RadioButtonAppend.GetActive() {
							obj.LabelOverwritten.SetLabel(sts["archUpdt"])
						} else {
							obj.LabelOverwritten.SetLabel(sts["archUpdt"])
						}
					}
				})
			} else if checkErr(err) {
				return
			}
			// Build 7z command
			args = buildCmdLineArgs(filename, opt.cmdLineList)
			// Use of goroutines to start an independent thread.
			go runCommand(filename, args, size)
		}
	} else {
		DialogMessage(obj.MainWindow, sts["warning"], "\n\n"+sts["noFile"], nil,
			DLG_OPT.DLG_NON_MODAL|DLG_OPT.DLG_WARNING|DLG_OPT.DLG_BTN_RELIEF_NONE, sts["continue"],
			func(resp int) {
				fmt.Println(resp)
			})
	}
	checkErr(err)
}

// ButtonCancelClicked:
func ButtonCancelClicked() {
	if execCmd != nil {
		// Stop progressbar
		ticker.Stop()
		obj.Progressbar.SetFraction(0)
		// Killing current 7za processe if exist.
		syscall.Kill(-execCmd.Process.Pid, syscall.SIGKILL)
		if _, err := os.Stat(ExtEnsure(opt.PrevFilename, ".tmp")); !os.IsNotExist(err) {
			os.RemoveAll(ExtEnsure(opt.PrevFilename, ".tmp"))
		}
	}
}

// ButtonExitClicked:
func ButtonExitClicked() {
	obj.MainWindow.Hide()
	windowDestroy(nil)
}

// EvenboxTopButtonReleaseEvent DialogBox
func EvenboxTopButtonReleaseEvent() {
	About.Show()
}

// RadioButtonGrpChanged: generic callback to handle archive mode creation.
func RadioButtonGrpChanged(gtkObj interface{}) {
	switch o := gtkObj.(type) {
	case *gtk.RadioButton:
		opt.Overwrite = obj.RadioButtonNew.GetActive()
		opt.Append = obj.RadioButtonAppend.GetActive()
		opt.Update = obj.RadioButtonUpdate.GetActive()
	case *gtk.CheckButton:
		obj.FrameWriteAs.SetSensitive(!o.GetActive())
		if o.GetActive() {
			opt.OverwriteBack = obj.RadioButtonNew.GetActive()
			opt.AppendBack = obj.RadioButtonAppend.GetActive()
			opt.UpdateBack = obj.RadioButtonUpdate.GetActive()

			obj.RadioButtonNew.SetActive(true)
		} else {
			obj.RadioButtonNew.SetActive(opt.OverwriteBack)
			obj.RadioButtonAppend.SetActive(opt.AppendBack)
			obj.RadioButtonUpdate.SetActive(opt.UpdateBack)
		}
	}
}
