// perms.go

/*
	Copyright ©2020-21 hfmrow - File operations v1.0 https://github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package fileOp

import (
	"fmt"
	"os"
)

/*
* Os Perms
 */

/*
Usage to Create any directories needed to put this file in them:

func DirTree(filename string){
    var dir_file_mode os.FileMode
	OS := OsPermsStructNew()

    dir_file_mode = os.ModeDir | (OS.USER_RWX | OS.ALL_R)
    os.MkdirAll(filename, dir_file_mode)
}

func TestPerms(){
	OS := OsPermsStructNew()

	permsFile := os.ModePerm & (OS.USER_RW | OS.GROUP_R | OS.OTH_R)     // 0644
	permsFileX := os.ModePerm & (OS.USER_RWX | OS.GROUP_RX | OS.OTH_RX) // 0755
	permsDir := os.ModePerm & (OS.USER_RWX | OS.GROUP_RX | OS.OTH_RX)   // 0755
	fmt.Printf("%#o, %#o, %#o, \n", permsFile, permsFileX, permsDir)
}

Memo:
   0 : - - - (aucun droit)
   1 : - - x (exécution)
   2 : - w - (écriture)
   3 : - w x (écriture et exécution)
   4 : r - - (lecture seule)
   5 : r - x (lecture et exécution)
   6 : r w - (lecture et écriture)
   7 : r w x (lecture, écriture et exécution)
*/

// Structure that hold file permissions
type PermsStruct struct {
	USER_R,
	USER_W,
	USER_X,
	USER_RW,
	USER_RX,
	USER_RWX,

	GROUP_R,
	GROUP_W,
	GROUP_X,
	GROUP_RW,
	GROUP_RX,
	GROUP_RWX,

	OTH_R,
	OTH_W,
	OTH_X,
	OTH_RW,
	OTH_RX,
	OTH_RWX,

	ALL_R,
	ALL_W,
	ALL_X,
	ALL_RW,
	ALL_RX,
	ALL_RWX,

	rEAD,
	wRITE,
	eX,
	uSER_SHIFT,
	gROUP_SHIFT,
	oTH_SHIFT,

	File,
	FileX,
	Dir os.FileMode
}

func PermsStructNew() (osp *PermsStruct) {
	osp = new(PermsStruct)

	osp.rEAD = 04
	osp.wRITE = 02
	osp.eX = 01
	osp.uSER_SHIFT = 6
	osp.gROUP_SHIFT = 3
	osp.oTH_SHIFT = 0

	osp.USER_R = osp.rEAD << osp.uSER_SHIFT
	osp.USER_W = osp.wRITE << osp.uSER_SHIFT
	osp.USER_X = osp.eX << osp.uSER_SHIFT
	osp.USER_RW = osp.USER_R | osp.USER_W
	osp.USER_RX = osp.USER_R | osp.USER_X
	osp.USER_RWX = osp.USER_RW | osp.USER_X

	osp.GROUP_R = osp.rEAD << osp.gROUP_SHIFT
	osp.GROUP_W = osp.wRITE << osp.gROUP_SHIFT
	osp.GROUP_X = osp.eX << osp.gROUP_SHIFT
	osp.GROUP_RW = osp.GROUP_R | osp.GROUP_W
	osp.GROUP_RX = osp.GROUP_R | osp.GROUP_X
	osp.GROUP_RWX = osp.GROUP_RW | osp.GROUP_X

	osp.OTH_R = osp.rEAD << osp.oTH_SHIFT
	osp.OTH_W = osp.wRITE << osp.oTH_SHIFT
	osp.OTH_X = osp.eX << osp.oTH_SHIFT
	osp.OTH_RW = osp.OTH_R | osp.OTH_W
	osp.OTH_RX = osp.OTH_R | osp.OTH_X
	osp.OTH_RWX = osp.OTH_RW | osp.OTH_X

	osp.ALL_R = osp.USER_R | osp.GROUP_R | osp.OTH_R
	osp.ALL_W = osp.USER_W | osp.GROUP_W | osp.OTH_W
	osp.ALL_X = osp.USER_X | osp.GROUP_X | osp.OTH_X
	osp.ALL_RW = osp.ALL_R | osp.ALL_W
	osp.ALL_RX = osp.ALL_R | osp.ALL_X
	osp.ALL_RWX = osp.ALL_RW | osp.ALL_X

	// Set files/Dirs permissions
	osp.File = osp.GetFileMode(osp.USER_RW | osp.GROUP_R | osp.OTH_R)     // 0644
	osp.FileX = osp.GetFileMode(osp.USER_RWX | osp.GROUP_RX | osp.OTH_RX) // 0755
	osp.Dir = osp.GetFileMode(osp.USER_RWX | osp.GROUP_RWX | osp.OTH_RX)  // 0775

	return
}

// GetFileMode: add os.ModePerm to os.FileMode, that permit to use
// it directly in read or write golang functions.
func (osp *PermsStruct) GetFileMode(perm os.FileMode) os.FileMode {
	return os.ModePerm & perm
}

// DispPerms: display right.
// i.e: DispPerms(OS.GROUP_RW | OS.USER_RW | OS.ALL_R) -> 664
func (osp *PermsStruct) DispPerms(perm os.FileMode) {
	fmt.Printf("os.ModePerm & %#o\n", perm)
}
