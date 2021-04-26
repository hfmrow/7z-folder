// filesOperations.go

/*
	Copyright ©2020-21 H.F.M - File operations v1.2 https://github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package fileOp

// /*
// #include <stdio.h>
// #include <dirent.h>

// int ldir(char *root)
// {
//     DIR *folder;
//     struct dirent *entry;
//     int files = 0;
// // test
//     folder = opendir(root);
//     if(folder == NULL)
//     {
//         perror("Unable to read directory");
//         return(1);
//     }

//     while( (entry=readdir(folder)) )
//     {
//         files++;
//         printf("File %3d: %s\n",
//                 files,
//                 entry->d_name
//               );
//     }

//     closedir(folder);

//     return(0);
// }
// */
// import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	// glss "github.com/hfmrow/genLib/slices"
)

// var (
// 	isExistSlIface = glss.IsExistSlIface
// )

// FilesOpStruct: Contains all the information of the scanned files
// (directories / files). This structure also contains many methods
// for handling files, filtering, directories, perms, owner options
// relating to the linux filesystem.
type FilesOpStruct struct {
	Masks []string

	Files     []FileDetails
	filesUndo []FileDetails // Used as backup, for cases were we want to go back.

	ExcludeMasks,
	AppendToFiles,
	Backup,
	IncludeDir,
	RealUserOwner,
	// don't add duplicate files default = true
	SkipDup bool

	SrceDir string

	Owner *OwnerStruct
	Perms *PermsStruct

	// ReadDir function specific options
	ReadDirDepth int // default 1
}

// FileDetails: Contains all the information about a scanned file.
// More convenient way than the native 'os.FileInfo' structure,
// all information is directly accessible.
type FileDetails struct {
	BaseDir,
	FileRel,
	LinkTarget,
	FilenameFull,
	Error string

	Mode os.FileMode

	IsDir,
	IsSymlink bool

	Size int64
	// Used in callback function for FilesRaw method, that permit to eventually change
	// filename on condition (target slink instead on slink...)
	FileAsRaw string

	info os.FileInfo
}

// restoreRight: NOT used
func (fos *FilesOpStruct) restoreRight(osPathname string) {
	fi, _ := os.Stat(osPathname)
	os.Chmod(osPathname, os.ModePerm&0)
	switch {
	case fi.IsDir():
		os.Chmod(osPathname, fos.Perms.Dir)
	default:
		os.Chmod(osPathname, fos.Perms.File)
	}
}

// FilesOpStructNew: Create a new 'FilesOpStruct' that hold some
// useful file/dir functions. The goal of it is to facilitate file
// interactions, can create copy, files/dir perserving
// owner/permissions. Files can be stored with all information
// needed to recreate them.
func FilesOpStructNew() (fos *FilesOpStruct, err error) {

	fos = new(FilesOpStruct)
	return fos, fos.init()
}

// init: Initialize permissions and current/real user information.
func (fos *FilesOpStruct) init() (err error) {

	fos.Perms = PermsStructNew()
	if fos.Owner, err = OwnerStructNew(); err != nil {
		fmt.Println("Error collecting real user information.")
	}
	fos.ReadDirDepth = 1
	fos.SkipDup = true
	return
}

// init: Initialize permissions and current/real user information.
func (fos *FilesOpStruct) ResetAll() error {

	// if len(fos.Masks) > 0 {
	fos.Masks = fos.Masks[:0]
	// }
	// if len(fos.Files) > 0 {
	fos.Files = fos.Files[:0]
	// }

	fos.ExcludeMasks,
		fos.AppendToFiles,
		fos.Backup,
		fos.IncludeDir,
		fos.RealUserOwner = false, false, false, false, false

	fos.SrceDir = ""
	return fos.init()
}

// SetOwner: Set filename user owner, uid and gid are defined in
// main structure, fos.RealUserOwner must be enabled. The
// 'callbeforefunc' permit to run an operation just before chown
// command (like saving png, jpg ...).
func (fos *FilesOpStruct) SetOwner(path string, callbeforefunc ...interface{}) (err error) {

	if fos.RealUserOwner {
		if len(callbeforefunc) > 0 {
			callbeforefunc[0].(func())()
		}
		err = os.Chown(path, fos.Owner.Uid, fos.Owner.Gid)
	}
	return
}

func (fos *FilesOpStruct) checkForBackup(filename string) (err error) {

	if fos.Backup {
		if _, err = os.Stat(filename); err == nil {
			err = os.Rename(filename, filename+"~")
		}
	}
	return
}

// WriteFile: can set owner (if requested).
func (fos *FilesOpStruct) WriteFile(filename string, data []byte, perm os.FileMode) (err error) {

	if err = fos.checkForBackup(filename); err == nil {

		if err = ioutil.WriteFile(filename, data, perm); err == nil {
			err = fos.SetOwner(filename)
		}
	}
	return
}

// Copy: file/dir tree preserving permissions & owner (if requested).
func (fos *FilesOpStruct) CopyFile(src, dest string) (err error) {

	var (
		info os.FileInfo
		d, s *os.File
	)
	if info, err = os.Lstat(src); err == nil {
		if err = fos.MkdirAll(filepath.Dir(dest), fos.Perms.Dir); err == nil {
			if err = fos.checkForBackup(dest); err == nil {
				if d, err = os.Create(dest); err == nil {
					if err = os.Chmod(d.Name(), info.Mode()); err == nil {
						if s, err = os.Open(src); err == nil {
							_, err = io.Copy(d, s)
							fos.SetOwner(dest)
							d.Close()
							s.Close()
						}
					}
				}
			}
		}
	}
	return err
}

// GetTotalSize: retrieve total files size currently stored.
func (fos *FilesOpStruct) GetTotalSize() uint64 {
	var size uint64
	for _, fd := range fos.Files {
		size += uint64(fd.Size)
	}
	return size
}

// FilesGetAsRaw: Get filenames in simpliest way as []string.
// Were provided, 'callbackOnAdd' function allow to filtering some
// file details before adding or not this one to results, in this
// case, 'FileDetails.FileAsRaw' must contain desired outpout
// that will be added to the returned strings slice.
func (fos *FilesOpStruct) FilesGetAsRaw(callbackOnAdd ...func(file *FileDetails) bool) []string {
	var (
		fn       func(file *FileDetails) bool = nil
		outFiles []string
		currFile string
	)
	if len(callbackOnAdd) > 0 {
		fn = callbackOnAdd[0]
	}
	for _, f := range fos.Files {
		switch fn {
		case nil:
			currFile = f.FilenameFull
		default:
			if fn(&f) {
				currFile = f.FileAsRaw
			} else {
				continue
			}
		}
		outFiles = append(outFiles, currFile)
	}
	return outFiles
}

// FileDetailsNew: Make a new 'FileDetails' Structure.
func (fos *FilesOpStruct) FileDetailsNew(path string, stat os.FileInfo) (fd *FileDetails, err error) {

	rel, err := filepath.Rel(fos.SrceDir, path)
	if err != nil {
		return nil, err
	}

	if rel == "." {
		return nil, nil
	}

	var linkTarget string
	isLink := stat.Mode()&os.ModeSymlink == os.ModeSymlink
	if isLink {
		if linkTarget, err = os.Readlink(path); err != nil {
			return nil, err
		}
	}

	fd = new(FileDetails)
	fd = &FileDetails{
		FilenameFull: path,
		BaseDir:      fos.SrceDir,
		FileRel:      rel,
		Mode:         stat.Mode(),
		Size:         stat.Size(),
		IsDir:        stat.IsDir(),
		IsSymlink:    isLink,
		LinkTarget:   linkTarget,
		info:         stat}

	return fd, nil
}

// func asPtrAndLength(s string) (*C.char, int) {
// 	addr := &s
// 	hdr := (*reflect.StringHeader)(unsafe.Pointer(addr))

// 	p := (*C.char)(unsafe.Pointer(hdr.Data))
// 	n := hdr.Len

// 	// reflect.StringHeader stores the Data field as a uintptr, not
// 	// a pointer, so ensure that the string remains reachable until
// 	// the uintptr is converted.
// 	runtime.KeepAlive(addr)

// 	return p, n
// }

// readDir:
func (fos *FilesOpStruct) readDir(path string, depthRecurse int) (filesDetails []FileDetails, err error) {

	var (
		currFilename string
		fd           *FileDetails
		tmpFD        []FileDetails
		filesInfos   []os.FileInfo
		fi           os.FileInfo
		f            *os.File
		depth        int
		accept       bool

		appendFilesDetails = func(fi os.FileInfo) error {
			currFilename = filepath.Join(path, fi.Name())
			if fd, err = fos.FileDetailsNew(currFilename, fi); err != nil {
				return err
			}
			if fd != nil {
				var isExist bool
				if fos.SkipDup {
					for _, f := range fos.Files {
						if f.FilenameFull == currFilename {
							isExist = true
							break
						}
					}
				}
				if !isExist {
					filesDetails = append(filesDetails, *fd)
				}
			}
			return nil
		}
	)
	// TODO
	// // "C" tests
	// fmt.Println("ldir:", C.ldir(C.CString(path)))
	// return
	if f, err = os.Open(path); err == nil {
		if fi, err = f.Stat(); err == nil {

			// If it's not a directory, return single file
			if !fi.IsDir() {
				path = filepath.Dir(path)
				err = appendFilesDetails(fi)
				if err != nil {
					log.Printf("[readDir/appendFilesDetails/!fi.IsDir]: %v", err)
				}
				return

				// fd, err = fos.FileDetailsNew(path, fi)
				// if fd != nil && err == nil {
				// 	return []FileDetails{*fd}, nil
				// }
				// errIn := fmt.Errorf("Unable to handle file: %s\n", path)
				// if err != nil {
				// 	errIn = fmt.Errorf("%v%v\n", errIn, err)
				// }
				// return nil, errIn
			}

			filesInfos, err = f.Readdir(-1)
			f.Close()
			for _, fi = range filesInfos {
				err = nil

				if !fi.IsDir() {
					if accept, err = fos.checkFileInfoMask(fi); !accept {
						if err != nil {
							log.Printf("[readDir/checkFileInfoMask]: %v", err)
							err = nil
						}
						continue
					}
					if err = appendFilesDetails(fi); err != nil {
						log.Printf("[readDir/appendFilesDetails]: %v", err)
						continue
					}
				} else {
					if fos.IncludeDir {
						if err = appendFilesDetails(fi); err != nil {
							log.Printf("[readDir/IncludeDir/appendFilesDetails]: %v", err)
							continue
						}
					}
					depth = depthRecurse
					if depth != 0 {
						depth--
						if tmpFD, err = fos.readDir(filepath.Join(path, fi.Name()), depth); err != nil {
							log.Printf("[readDir/recursive]: %v", err)
						}
						filesDetails = append(filesDetails, tmpFD...)
					}
				}
			}
		}
	}
	return
}

// GetFilesDetailsReadDir: Retrieve files details starting at "root",
// and limit depth scanning directory using 'ReadDirDepth' struct variable,
// details will be stored as 'FileDetails' structure in fos.Files slice.
func (fos *FilesOpStruct) GetFilesDetailsReadDir(root string) (err error) {

	fi, err := os.Lstat(root)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		fos.SrceDir = root
	} else {
		fos.SrceDir = filepath.Dir(root)
	}
	tmpFilesDetails, err := fos.readDir(root, fos.ReadDirDepth)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	// // Backup previous set of files
	// fos.backupFilesDetails()

	if !fos.AppendToFiles {
		fos.Files = tmpFilesDetails
	} else {
		fos.Files = append(fos.Files, tmpFilesDetails...)
	}

	return
}

// backupFilesDetails: backup previous set of 'Files', usefull when append is used
func (fos *FilesOpStruct) BackupFilesDetails() {
	fos.filesUndo = make([]FileDetails, len(fos.Files))
	copy(fos.filesUndo, fos.Files)
}

// UndoLastFilesChanges: restore previous 'Files' state, usefull when append is used
func (fos *FilesOpStruct) UndoLastFilesChanges() {
	fos.Files = make([]FileDetails, len(fos.filesUndo))
	copy(fos.Files, fos.filesUndo)
}

// checkFileInfoMask: Check for Pattern exclusion using 'os.FileInfo' parameter.
func (fos *FilesOpStruct) checkFileInfoMask(fi os.FileInfo) (accept bool, err error) {

	if len(fos.Masks) != 0 {
		for _, mask := range fos.Masks {
			if accept, err = filepath.Match(mask, fi.Name()); err != nil {
				log.Printf("[checkFileInfoMask]: %v", err)
				return false, err
			} else if accept {
				break
			}
		}
		switch {
		case accept && !fos.ExcludeMasks:
			return true, nil
		case accept && fos.ExcludeMasks:
			return false, nil
		case !accept && !fos.ExcludeMasks:
			return false, nil
		}
	}
	return true, nil
}

// GetFilesDetails: Retrieve files details starting at "root", they
// will be stored as 'FileDetails' structure in fos.Files slice.
func (fos *FilesOpStruct) GetFilesDetails(root string) (err error) {
	var (
		ok   bool
		stat os.FileInfo

		appendFilesDetails = func(path string, fi os.FileInfo) error {
			fd, err := fos.FileDetailsNew(path, fi)
			if err != nil {
				return err
			}
			if fd != nil {
				var isExist bool
				if fos.SkipDup {
					for _, f := range fos.Files {
						if f.FilenameFull == path {
							isExist = true
							break
						}
					}
				}
				if !isExist {
					fos.Files = append(fos.Files, *fd)
				}
			}
			return nil
		}
	)
	// // Backup previous set of files
	// fos.backupFilesDetails()
	// Clear 'Files' entries whether it needed
	if !fos.AppendToFiles {
		fos.Files = fos.Files[:0]
	}

	fi, err := os.Lstat(root)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		fos.SrceDir = root
	} else {
		fos.SrceDir = filepath.Dir(root)
	}
	if stat, err = os.Lstat(root); err == nil {

		switch {
		case stat.IsDir():

			return filepath.Walk(fos.SrceDir,
				func(path string, fi os.FileInfo, err error) error {

					switch {
					case err != nil:
						return err
					case fi == nil:
						return fmt.Errorf("Unable to stat: %v", path)
					case fi.IsDir() && fos.IncludeDir:
						return appendFilesDetails(path, stat)
					case !fi.IsDir():
						if ok, err = fos.checkFileInfoMask(fi); err != nil {
							log.Printf("[GetFilesDetails/Walk/checkFileInfoMask]: %v", err)
							return err
						}
						if ok {
							return appendFilesDetails(path, fi)
						}
					}
					return nil
				})
		default:
			fos.SrceDir = filepath.Dir(root)
			if ok, err = fos.checkFileInfoMask(stat); err != nil {
				log.Printf("[GetFilesDetails/checkFileInfoMask]: %v", err)
				return err
			}
			if ok {
				return appendFilesDetails(root, stat)
			}
		}
	}
	return
}

// TODO Remove this one
// GetFilesDetails: Retrieve files details starting at "root", they
// will be stored as 'FileDetails' structure in fos.Files slice.
func (fos *FilesOpStruct) GetFilesDetailsOld(root string) (err error) {
	var (
		ok   bool
		stat os.FileInfo
	)

	// // Backup previous set of files
	// fos.backupFilesDetails()
	// Clear 'Files' entries whether it needed
	if !fos.AppendToFiles {
		fos.Files = fos.Files[:0]
	}

	var checkFile = func(inFileDetails *FileDetails) (err error) {

		if inFileDetails != nil {

			// Check for pattern exclusion.
			if len(fos.Masks) != 0 {
				for _, mask := range fos.Masks {
					if ok, err = filepath.Match(
						mask,
						filepath.Base(inFileDetails.FilenameFull)); err != nil {

						return
					}

					switch {
					case ok && !fos.ExcludeMasks:
						fos.Files = append(fos.Files, *inFileDetails)
					case !ok && fos.ExcludeMasks:
						fos.Files = append(fos.Files, *inFileDetails)
					}
				}
			} else {
				fos.Files = append(fos.Files, *inFileDetails)
			}
		}
		return
	}

	fos.SrceDir = root

	if stat, err = os.Lstat(fos.SrceDir); err == nil {
		switch {
		case stat.IsDir():

			if err = filepath.Walk(fos.SrceDir, func(path string, fInfo os.FileInfo, err error) error {

				switch {
				case err != nil:
					return err
				case fInfo == nil:
					return fmt.Errorf("Error stat: %v", path)
				}

				f, err := fos.FileDetailsNew(path, fInfo)
				if err != nil {
					return err
				}
				return checkFile(f)

			}); err != nil {
				return err
			}

		default:

			fos.SrceDir = filepath.Dir(root)
			f, err := fos.FileDetailsNew(root, stat)
			if err != nil {
				return err
			}
			err = checkFile(f)
			if err != nil {
				return err
			}
		}
	}
	return
}

// MkdirAll: This is a modified version of original golang function,
// this one set the defined (specific) owner for each created directories.
func (fos *FilesOpStruct) MkdirAll(path string, perm os.FileMode) error {

	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{"mkdir", path, syscall.ENOTDIR}
	}

	i := len(path) // Skip trailing path separator.
	for i > 0 && os.IsPathSeparator(path[i-1]) {
		i--
	}

	j := i // Scan backward over element.
	for j > 0 && !os.IsPathSeparator(path[j-1]) {
		j--
	}

	if j > 1 {
		err = fos.MkdirAll(path[:j-1], perm)
		if err != nil {
			return err
		}
	}

	if err = os.Mkdir(path, perm); err != nil {

		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}

	// Change dir owner if required
	if err = fos.SetOwner(path); err != nil {
		return err
	}

	return nil
}

// Sort: Descending as default behavior
func (fos *FilesOpStruct) SortFiles(ascend ...bool) {

	var ascending bool

	if len(ascend) > 0 {
		ascending = ascend[0]
	}

	switch ascending {

	case false:
		sort.SliceStable(fos.Files,
			func(i, j int) bool {
				return strings.ToLower(fos.Files[i].FilenameFull) <
					strings.ToLower(fos.Files[j].FilenameFull)
			})

	case true:
		sort.SliceStable(fos.Files,
			func(i, j int) bool {
				return strings.ToLower(fos.Files[i].FilenameFull) >
					strings.ToLower(fos.Files[j].FilenameFull)
			})
	}
}

// Read FilesOpStruct from file
func (fos *FilesOpStruct) Read(filename string) (err error) {

	var textFileBytes []byte

	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(textFileBytes, &fos)
	}
	if err != nil {
		fmt.Printf("Enable to read 'FilesOpStruct' file: %s\n", err.Error())
	}
	return
}

// Write FilesOpStruct to file
func (fos *FilesOpStruct) Write(filename string) (err error) {

	var (
		jsonData []byte
		out      bytes.Buffer
	)

	if jsonData, err = json.Marshal(&fos); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm)
		}
	}
	return err
}
