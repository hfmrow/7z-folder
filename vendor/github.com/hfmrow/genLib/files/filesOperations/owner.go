// owner.go

/*
	Copyright ©2020-21 hfmrow - File operations v1.0 https://github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package fileOp

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

type OwnerStruct struct {
	Uid,
	Gid int

	Root bool

	RealUser,
	CurrentUser *user.User
}

// OwnerStructNew: This structure gathers file permissions, the files
// owner change methods and Real, Current user information.
func OwnerStructNew() (owns *OwnerStruct, err error) {

	owns = new(OwnerStruct)
	if owns.CurrentUser, owns.RealUser, owns.Root, err = GetRootCurrRealUser(); err == nil {
		if owns.Root {
			if owns.Uid, err = strconv.Atoi(owns.RealUser.Uid); err == nil {
				owns.Gid, err = strconv.Atoi(owns.RealUser.Gid)
			}
		} else {
			if owns.Uid, err = strconv.Atoi(owns.CurrentUser.Uid); err == nil {
				owns.Gid, err = strconv.Atoi(owns.CurrentUser.Gid)
			}

		}
	}
	return
}

// ChangeFileOwner: set the owner of the file to user defined
// in 	Uid and Gid. If not modified, values target current realuser.
// Obtaining root rights differs when sudo or pkexec is
// used, this function checks which command was used and acts to
// perform the job correctly.
func (owns *OwnerStruct) ChangeFileOwner(filename string) (err error) {
	return os.Chown(filename, owns.Uid, owns.Gid)
}

// GetFileOwner: Work only for linux users.
func (owns *OwnerStruct) GetFileOwner(filename string) (uid, gid int, ok bool, err error) {
	return GetFileOwner(filename)
}

// GetFileOwner: Work only for linux users.
func GetFileOwner(filename string) (uid, gid int, ok bool, err error) {

	var info os.FileInfo
	if info, err = os.Lstat(filename); err == nil {
		if stat, yes := info.Sys().(*syscall.Stat_t); ok {
			ok = yes
			uid = int(stat.Uid)
			gid = int(stat.Gid)
		} else {
			// If not linux, return current uid and gid
			ok = false
			uid = os.Getuid()
			gid = os.Getgid()
		}
	}
	return
}

// ChangeFileOwner: set the owner of the file to real user instead
// of root. Obtaining root rights differs when sudo or pkexec is
// used, this function checks which command was used and acts to
// perform the job correctly.
func ChangeFileOwner(filename string) (err error) {
	var (
		uid,
		gid int
		root     bool
		realUser *user.User
	)
	if _, realUser, root, err = GetRootCurrRealUser(); err == nil {
		if root {
			if uid, err = strconv.Atoi(realUser.Uid); err == nil {
				if gid, err = strconv.Atoi(realUser.Gid); err == nil {
					err = os.Chown(filename, uid, gid)
				}
			}
		}
	}
	return
}

// GetRootCurrRealUser: Retrieve information about root state
// and current and real user. Handle different behavior between
// sudo and pkexec commands.
func GetRootCurrRealUser() (currentUser, realUser *user.User, root bool, err error) {
	realUser = new(user.User)
	if currentUser, err = user.Current(); err == nil {
		*(realUser) = *(currentUser)                                  // Duplicate content
		if root = (currentUser.Uid + currentUser.Gid) == "00"; root { // root rights acquired ?
			if realUser, err = user.LookupId(os.Getenv("PKEXEC_UID")); err != nil {
				// fmt.Println("pkexec used")
				if realUser, err = user.Lookup(os.Getenv("SUDO_USER")); err != nil {
					// fmt.Println("sudo used")
					return
				}
			}
		}
	}
	return
}
