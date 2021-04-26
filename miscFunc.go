// miscFunc.go

/*
	Source file auto-generated on Mon, 26 Apr 2021 08:45:09 using Gotk3 Objects Handler v1.7.8
	©2018-21 hfmrow https://hfmrow.github.io

	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 hfmrow - 7z folder v1.6 github.com/hfmrow/7z-folder

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	glfsfo "github.com/hfmrow/genLib/files/filesOperations"
)

// checkFiles:
func checkFiles() (filesList []string, size int64, next bool, err error) {

	var (
		brokenSl         []string
		removedFilesSize int64
	)
	next = true
	filesList = fos.FilesGetAsRaw(
		// this callback is used to filtering / prevent broken links errors
		func(file *glfsfo.FileDetails) bool {

			var ok = true
			size += file.Size

			if obj.CheckbuttonFollowSymlink.GetActive() {
				if file.IsSymlink {
					if _, err = os.Lstat(file.LinkTarget); err != nil {
						err = nil
						brokenSl = append(brokenSl, file.FilenameFull)
						removedFilesSize += file.Size
						ok = false
					}
				}
			}
			// if ok {
			// 	upDir := CheckOutsideRoot(baseDir, file.FilenameFull)
			// 	// Remove unwanted files (Specific to '7za' command)
			// 	switch {
			// 	case upDir == -1, upDir == 0:
			// 		if upDir == 0 && !file.IsDir {
			// 			ok = false
			// 		} else {
			// 			// file is at the same level of the base directory
			// 			file.FileAsRaw = file.FilenameFull
			// 		}
			// 	case upDir > 0:
			// 		// file is in a sub-dir of the base directory
			// 		ok = false
			// 	case upDir < -1:
			// 		ok = false
			// 		fmt.Println("WARNING [ButtonCompressClicked/Rel]: Outside the scope of the current root directory")
			// 	}
			// }
			return ok
		})
	// Check for broken Slinks
	if len(brokenSl) > 0 {
		var tmpBSl []string
		for _, sl := range brokenSl {
			tmpBSl = append(tmpBSl, TruncatePath(sl, 2))
		}
		if DialogMessage(obj.MainWindow, sts["confirm"],
			nl+nl+sts["brkSl"]+nl+nl+strings.Join(tmpBSl, GetOsLineEnd()),
			nil, DLG_OPT.DLG_INFO|DLG_OPT.DLG_BTN_RELIEF_NONE, sts["cancel"], sts["continue"]) != 1 {
			next = false
			return
		} else {
			for i := len(filesList) - 1; i >= 0; i-- {
				for _, bsl := range tmpBSl {
					if filesList[i] == bsl {
						filesList = append(filesList[:i], filesList[i+1:]...)
					}
				}
			}
			size -= removedFilesSize
		}
	}

	return
}

// buildCmdLineArgs:
func buildCmdLineArgs(filename string, filelist []string) (args []string) {
	// Archive arg
	switch {
	case obj.RadioButtonAppend.GetActive(), obj.RadioButtonNew.GetActive():
		args = append(args, "a")
	case obj.RadioButtonUpdate.GetActive():
		args = append(args, "u")
	}
	// Follow symlink arg
	if obj.CheckbuttonFollowSymlink.GetActive() {
		args = append(args, "-l")
	}
	// Compressor arg
	args = append(args, "-m0=lzma2")
	// Comp Lvl arg
	args = append(args, "-mx="+fmt.Sprintf("%d", obj.SpinbuttonCompressionLvl.GetValueAsInt()))
	// Cpu used arg
	args = append(args, "-mmt="+fmt.Sprintf("%d", runtime.NumCPU()))
	// Dict size mb
	args = append(args, "-md="+fmt.Sprintf("%dm", obj.SpinbuttonDictSize.GetValueAsInt()))
	// Dest filename
	args = append(args, filename)
	// input filenames
	return append(args, filelist...)
}

// runCommand:
func runCommand(filename string, args []string, size int64) {
	// Progressbar display
	ticker = time.NewTicker(time.Millisecond * 100)
	go func() {
		for _ = range ticker.C {
			// idle run to permit gtk3 working right during goroutine
			glib.IdleAdd(func() {
				obj.Progressbar.Pulse()
			})
		}
	}()
	var (
		term  []byte
		timer = BenchNew(false)
		err   error
		fi    os.FileInfo
	)
	// TODO remove on prod
	fmt.Println(args)

	timer.Lapse()
	execCmd = exec.Command(opt.SevenZip, args...)
	opt.PrevFilename = filename
	execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	term, err = execCmd.CombinedOutput()
	if err == nil {
		outTerm, ok := getTermResult(term)
		glib.IdleAdd(func() {
			if !ok {
				DialogMessage(obj.MainWindow, "\n\n"+sts["unexpectedEnd"]+nl+nl+strings.Join(outTerm, GetOsLineEnd()),
					nil, DLG_OPT.DLG_INFO|DLG_OPT.DLG_BTN_RELIEF_NONE, sts["continue"])
			} else {
				if fi, err = os.Stat(filename); checkErr(err) {
					err = nil
					return
				}
				size = fi.Size()
				obj.LabelPackedSize.SetLabel(HumanReadableSizeAlt(size))

				if err = os.Rename(filename, ExtEnsure(filename, ".7z")); checkErr(err) {
					err = nil
					return
				}
				timer.Stop()
				obj.LabelDuration.SetLabel(timer.Lapses[0].StringShort)
			}
		})
	} else {
		glib.IdleAdd(func() {
			// Get terminal error
			splitted := strings.Split(string(term), GetOsLineEnd())
			for idx, line := range splitted {
				if strings.Contains(line, "ERROR:") || strings.Contains(line, "WARNING:") {
					splitted = splitted[idx:]
				}
			}
			DialogMessage(obj.MainWindow, sts["info"],
				nl+nl+sts["permissionError"]+nl+nl+strings.Join(splitted, GetOsLineEnd()),
				nil, DLG_OPT.DLG_ERROR|DLG_OPT.DLG_BTN_RELIEF_NONE, sts["ok"])
		})
	}
	// Stop progressbar
	ticker.Stop()
	obj.Progressbar.SetFraction(0)
}

// dispStatusbar: Display statusbar information
func dispStatusbar() {
	if len(opt.cmdLineList) > 0 {
		statusbar.Set(fmt.Sprintf("%s", TruncatePath(baseDir, 2)), 0)
		statusbar.Set(fmt.Sprintf("%d", len(fos.Files)), 1)
	}
}

// checkErr: show error dialog if not nil and return true if an error was found.
func checkErr(e error) bool {
	if e != nil {
		glib.IdleAdd(func() {
			DialogMessage(obj.MainWindow, sts["error"], "\n\n"+sts["unexpectedError"]+nl+nl+e.Error(),
				nil, DLG_OPT.DLG_ERROR|DLG_OPT.DLG_BTN_RELIEF_NONE, sts["continue"])
		})
		return true
	}
	return false
}

// storeFiles: retrieve files
func storeFiles() error {

	var (
		err error
		fi  os.FileInfo
	)

	if checkErr(fos.ResetAll()) {
		return nil
	}

	fos.SkipDup = true
	// fos.ReadDirDepth = -1
	// fos.IncludeDir = true

	fos.BackupFilesDetails()
	backupValue := fos.AppendToFiles
	for _, path := range opt.cmdLineList {
		if fi, err = os.Lstat(path); err == nil {
			if fi.IsDir() {
				fos.AppendToFiles = backupValue
			} else {
				fos.AppendToFiles = true
			}
			err = fos.GetFilesDetails(path)
		}
		if checkErr(err) {
			err = nil
			continue
		}
	}
	// Determine which base folder we will work with.
	if len(baseDir) == 0 {
		baseDir = opt.cmdLineList[0]
		if fi, err := os.Stat(baseDir); err == nil {
			if !fi.IsDir() {
				baseDir = filepath.Dir(baseDir)
			}
		}
	}

	// // Control directory upstream if there is.
	// if upDir := CheckOutsideRoot(baseDir, fos.SrceDir); upDir > 0 {

	// 	txt := fmt.Sprintf("\n\n<span foreground=\"#760000FF\"><b>%s %d %s</b></span>\n\n<span foreground=\"#583939FF\"><b>%s</b></span>\n<span foreground=\"#3465A4FF\">%s</span>\n<span foreground=\"#583939FF\"><b>%s</b></span>\n<span foreground=\"#3465A4FF\">%s</span>\n\n<b>%s</b>",
	// 		sts["outBaseDir1"], upDir, sts["outBaseDir2"], sts["outBaseDir3"],
	// 		TruncatePath(baseDir, int(math.Abs(float64(upDir))+2)), sts["outBaseDir4"],
	// 		TruncatePath(fos.SrceDir, int(math.Abs(float64(upDir))+2)), sts["notAllowed"])

	// 	DialogMessage(obj.MainWindow, sts["info"], txt, nil,
	// 		DLG_OPT.DLG_BTN_INFO|DLG_OPT.DLG_BTN_RELIEF_NONE|DLG_OPT.DLG_MARKUP, sts["continue"])
	// 	// remove files
	// 	fos.UndoLastFilesChanges()
	// }

	// fos.SortFiles()
	// fmt.Println("baseDir:", baseDir)

	glib.IdleAdd(func() {
		obj.LabelFileTreated.SetLabel(fmt.Sprintf("%d "+sts["sbWFiles"], len(fos.Files)))
		obj.LabelUnpackedSize.SetLabel(HumanReadableSizeAlt(int64(fos.GetTotalSize())))
		dispStatusbar()
	})
	return nil
}

// getTermResult: Parse terminal output to get compressed files after 7za was done.
func getTermResult(inBytes []byte) ([]string, bool) {
	var (
		ok       bool
		regStart = regexp.MustCompile(`^(Everything is Ok)`)
		eol      = GetTextEOL(inBytes)
		lines    = strings.Split(string(inBytes), eol)
	)
	for idx := 0; idx < len(lines); idx++ {
		if regStart.MatchString(lines[idx]) {
			ok = true
		}
	}
	return lines, ok
}

//computeFilename:
func computeFilename(filename string) (outFilename string) {
	var (
		value int
		digit string
	)

	filename = filepath.Base(filename)
	filename = BaseNoExt(filename)

	digitE := regexp.MustCompile(`(?m)(\[\d*\])`)
	if digitE.MatchString(filename) {
		digit = digitE.FindString(filename)
		outFilename = filename[:len(filename)-len(digit)]
		digit = strings.TrimSuffix(strings.TrimPrefix(digit, "["), "]")
		value, _ = strconv.Atoi(digit)
		outFilename += fmt.Sprintf("[%02d]", value+1)
	} else {
		outFilename = filename + fmt.Sprintf("[%02d]", value)
	}

	return outFilename + ".tmp"
}

// sortNumbered: sort by numbering methode
func sortNumbered(fName string, inString []string) (newList []string) {
	var present bool
	ext := filepath.Ext(fName)
	fName = BaseNoExt(fName)

	reg := regexp.MustCompile(fName + `(\[\d*\])` + ext)
	var tmpWordList []string
	for idx := len(inString) - 1; idx >= 0; idx-- {
		if reg.MatchString(filepath.Base(inString[idx])) {
			tmpWordList = append(tmpWordList, filepath.Base(inString[idx]))
			newList = append(newList, inString[idx])
			present = true
		}
	}
	if present {
		numberedWords := new(WordWithDigit)
		numberedWords.Init(tmpWordList)
		sort.SliceStable(newList, func(i, j int) bool {
			return numberedWords.FillWordToMatchMaxLength(newList[i], true) < numberedWords.FillWordToMatchMaxLength(newList[j], true)
		})
	} else {
		newList = append(newList, fName+ext)
	}
	return newList
}

// SpinbuttonSetValues: Configure spin button
func SpinbuttonSetValues(sb *gtk.SpinButton, min, max, value int, step ...int) (err error) {
	incStep := 1
	if len(step) != 0 {
		incStep = step[0]
	}
	if ad, err := gtk.AdjustmentNew(float64(value), float64(min), float64(max), float64(incStep), 0, 0); err == nil {
		sb.Configure(ad, 1, 0)
	}
	return err
}

// cmdLineArg: Check for existing cmd line arguments
func cmdLineArg() (list []string) {
	lengthArgs := len(os.Args)
	if lengthArgs > 1 {
		for idx := 1; idx < lengthArgs; idx++ {
			list = append(list, os.Args[idx])
			obj.MainWindow.ShowAll()
		}
	}
	return list
}

// CheckCmd: Check for command if exist
func CheckCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return false
	}
	return true
}
