package files

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	glsg "github.com/hfmrow/genLib/strings"
)

// Retrieve current realpath and options filename. Options/Config path
// depend on devMode value, true means current directory, false means
// current/real user directory '~/.config/Creat/appName'
// 'realUser' if provided will set the current user (used for root access).
func GetAbsRealPath(optionDir string, devMode bool, realUser ...*user.User) (absoluteRealPath, optFilename string) {

	var (
		err                    error
		cUser                  *user.User
		base, absoluteBaseName string

		// Set wanted extension
		setExt = func(filename, ext string) (out string) {
			return filename[:len(filename)-len(path.Ext(filename))] + ext
		}

		// Owning directories recursively until '$HOME/.config'
		ownDirs = func(u *user.User, path string) (errIn error) {
			var uid, gid int
			if uid, errIn = strconv.Atoi(u.Uid); errIn == nil {
				if gid, errIn = strconv.Atoi(u.Gid); errIn == nil {
					if path, errIn = filepath.Rel(u.HomeDir, path); errIn == nil {

						for path != ".config" {
							if errIn = os.Chown(filepath.Join(u.HomeDir, path), uid, gid); errIn != nil {
								return
							}
							splitted := strings.Split(path, string(os.PathSeparator))
							path = filepath.Join(splitted[:len(splitted)-1]...)
						}
					}
				}
			}
			return
		}
	)

	if len(realUser) > 0 {
		cUser = realUser[0]
	} else {
		if cUser, err = user.Current(); err != nil {
			log.Fatal(err)
		}
	}

	if absoluteBaseName, err = os.Executable(); err == nil {
		absoluteRealPath, base = filepath.Split(absoluteBaseName)
		configPath := absoluteRealPath
		baseNoExt := setExt(base, "")

		if !devMode {
			configPath = filepath.Join(cUser.HomeDir, ".config", glsg.ToCamel(optionDir), baseNoExt)
			if _, err = os.Stat(configPath); os.IsNotExist(err) {
				if err = os.MkdirAll(configPath, 0755); err == nil {
					err = ownDirs(cUser, configPath)
				}
			}
		}
		if err == nil {
			optFilename = setExt(filepath.Join(configPath, baseNoExt), ".opt")
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return
}

// DirKeepSame: returns 'same' which corresponds to the equal parts of the
// two files and 'diff' which corresponds to the change between the two,
// the 'f1' argument must be the shortest path.
func DirKeepSame(f1, f2 string) (same, diff []string) {

	var (
		one = strings.Split(f1, string(os.PathSeparator))
		two = strings.Split(f2, string(os.PathSeparator))
	)
	for idx, dir1 := range one {
		for i := idx; i < len(two); i++ {
			t := two[i]
			if dir1 == t {
				same = append(same, dir1)
				break
			} else {
				diff = append(diff, t)
			}
		}
	}
	return append([]string{string(os.PathSeparator)}, same...), diff
}

// CheckOutsideRoot: Check if 'filePath' is outside the current scope of the
// 'rootDir' (must be a DIRECTORY).
// count = 0 means 'filePath' is at the same level than 'rootDir'.
// count > 0 means 'filePath' is in a sub-directory of 'rootDir' and the value
//           indicate how many sub-directories.
// count < 0 means how many parent directories relative to base directory.
func CheckOutsideRoot(rootDir, filePath string) (count int) {
	rel, err := filepath.Rel(rootDir, filePath)
	if err != nil {
		log.Printf("[CheckOutsideRoot], Unable to find relative path: %v\n", err)
		return -1
	}
	for strings.HasPrefix(rel, "..") {
		rel = strings.TrimPrefix(strings.TrimPrefix(rel, ".."), string(os.PathSeparator))
		count--
	}
	if count == 0 {
		count = len(strings.Split(rel, string(os.PathSeparator))) - 1
	}
	return count
}

// CheckOutsideDir: Check if 'filenameNew' is in the same directory as
// 'filenameOrig'. returns 'false' if the 'root' directory is the same.
func CheckOutsideDir(baseFilename, newFilename string) bool {

	log.Printf("[CheckOutsideDir], This function is outdated, please use '%s' instead\n", "CheckOutsideRoot")

	re := regexp.MustCompile(`^(\.\.[/\\]|\.)`)

	if re.MatchString(baseFilename) || re.MatchString(newFilename) {
		return true
	}

	if rel, err := filepath.Rel(
		filepath.Dir(baseFilename),
		filepath.Dir(newFilename)); err == nil {
		if re.MatchString(rel) {
			return true
		}
	}
	return false
}

// FileMatch: If a pattern contained in "patterns" match,
// true is returned
func ExtFileMatch(name string, patterns []string) bool {
	for _, pattern := range patterns {
		if match, err := filepath.Match(pattern, name); match {
			return match
		} else if err != nil {
			fmt.Printf("FileMatch: %s, %s, - %s\n", pattern, name, err.Error())
		}
	}
	return false
}

// sizeToBytes: convert uint32 size to bytes representation (32bits) "000001f1"
func SizeToBytes(size uint32) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, size); err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

// ReducePath: Reduce path length by preserving count element from the end
func TruncatePath(fullpath string, count ...int) (reduced string) {
	elemCnt := 2
	if len(count) != 0 {
		elemCnt = count[0]
	}
	splited := strings.Split(fullpath, string(os.PathSeparator))
	if len(splited) > elemCnt+1 {
		return "..." + string(os.PathSeparator) + filepath.Join(splited[len(splited)-elemCnt:]...)
	}
	return fullpath
}

// GenFileName: Generate a randomized file name
func GenFileName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(r.Int63n(time.Now().UnixNano())))))
}

// splitPath: make a slice from a string path
func SplitPath(path string) (outSlice []string) {
	// remove leading and ending "/"
	path = strings.Trim(path, string(os.PathSeparator))
	return strings.Split(path, string(os.PathSeparator))
}

// removePathBefore: remove directories before or after the chosen name
func RemovePathBefore(path []string, at string, after ...bool) []string {
	var afterMark bool
	if len(after) > 0 {
		afterMark = after[0]
	}
	for idx := len(path) - 1; idx >= 0; idx-- {
		if path[idx] == at {
			if afterMark {
				path = path[idx+1:]
			} else {
				path = path[idx:]
			}
			break
		}
	}
	return path
}

// IsDirOrSymlinkDir: File is a directory or a symlinked directory ?
// Need: os.FileInfo - > os.Lstat
func IsDirOrSymlinkDir(slRoot string, slStat os.FileInfo) (slIsDir bool) {
	var err error
	var fName string
	if slStat.IsDir() {
		return true
	} else if slStat.Mode()&os.ModeSymlink != 0 {
		if fName, err = os.Readlink(filepath.Join(slRoot, slStat.Name())); err == nil {
			if slStat, err = os.Stat(fName); err == nil {
				if slStat.IsDir() {
					return true
				}
			}
		}
	}
	if err != nil {
		log.Printf("Unable to scan: %s\n", err.Error)
	}
	return
}

// IsDirEmpty:
func IsDirEmpty(name string) (empty bool, err error) {
	var f *os.File
	if f, err = os.Open(name); err == nil {
		defer f.Close()
		if _, err = f.Readdirnames(1); err == io.EOF {
			return true, nil
		}
	}

	return false, err
}

// GetCurrentDir: Get current directory
func GetCurrentDir() (dir string, err error) {
	return os.Getwd()
}

// TempFileName generates a temporary filename for use
// in testing or whatever
func TempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	// return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
	return prefix + hex.EncodeToString(randBytes) + suffix
}

// TempMake: Make temporary directory
func TempMake(prefix string) string {
	dir, err := ioutil.TempDir("", prefix+"-")
	if err != nil {
		log.Fatalf("Unablme to create temp directory: %s\n", err.Error())
	}
	return dir + string(os.PathSeparator)
}

// TempRemove: Remove directory recursively
func TempRemove(fName string) (err error) {
	if err = os.RemoveAll(fName); err != nil {
		return (err)
	}
	return nil
}

// EEOptions: options for 'ExtEnsure' function.
type EEOptions int

const (
	EE_DEFAULT           EEOptions = 1 << 0
	EE_REM_ALL_AFTER_DOT EEOptions = 1 << 1 // Remove all after the first dot
	EE_PRESERVE_FILENAME EEOptions = 1 << 2 // Preserve full filename, just adding extension
)

// ExtEnsure: Ensuring that the filename has the desired extension.
// 'option' correspond to 'EEOptions' declaration above.
func ExtEnsure(filename, ext string, option ...EEOptions) string {

	var (
		opt            EEOptions = EE_DEFAULT
		origExtRemoved bool
	)
	if len(option) > 0 {
		opt = option[0]
	}

	// Usually used when filename is a directory
	if opt&EE_PRESERVE_FILENAME != 0 {
		return filename + ext
	}

	// For some double extensions, i.e: "tar.gz", remove before apply
	if strings.HasSuffix(filename, ext) {
		filename = strings.TrimSuffix(filename, ext)
		origExtRemoved = true
	}

	// Whether removing all after 1st dot is requested
	if opt&EE_REM_ALL_AFTER_DOT != 0 {
		tmpSplitted := strings.Split(filename, ".")
		if len(tmpSplitted) > 0 {
			filename = tmpSplitted[0]
		}
	}

	if origExtRemoved {
		return filename + ext
	}
	return filename[:len(filename)-len(path.Ext(filename))] + ext
}

// BaseNoExt: get only the name without ext.
func BaseNoExt(filename string) (outFilename string) {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filepath.Base(filename)))
}

// magic number mime detection
var magicTable = map[string]string{
	"\x37\x7A\xBC\xAF\x27\x1C\x00\x04": "7zip",
	"\xFD\x37\x7A\x58\x5A\x00\x00":     "xz",
	"\x1F\x8B\x08\x00\x00\x09\x6E\x88": "gzip",
	"\x75\x73\x74\x61\x72":             "tar",
	"\x7F\x45\x4C\x46\x02\x01\x01":     "executable-linux",
}

// GetSimpleFileType: scan first bytes to detect type of file.
// This is a basic version for kick usage.
func GetSimpleFileType(filename string) string {
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		buffReader := bufio.NewReader(file)
		for magic, mime := range magicTable {
			if peeked, err := buffReader.Peek(len([]byte(magic))); err == nil {
				tmpMagic := []byte(magic)
				if bytes.Index(peeked, tmpMagic) == 0 {
					return mime
				}
			}
		}
	}
	return "unknown"
}

// FileInfos: As remind for available file information.
func FileInfos(filename string) error {

	fileInfo, err := os.Lstat(filename)

	if err != nil {
		return err
	}

	fmt.Println("Name : ", fileInfo.Name())

	fmt.Println("Size : ", fileInfo.Size())

	fmt.Println("Mode/permission : ", fileInfo.Mode())

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		fmt.Println("File is a symbolic link")
	}

	fmt.Println("Modification Time : ", fileInfo.ModTime())

	fmt.Println("Is a directory? : ", fileInfo.IsDir())

	fmt.Println("Is a regular file? : ", fileInfo.Mode().IsRegular())

	fmt.Println("Unix permission bits? : ", fileInfo.Mode().Perm())

	fmt.Println("Permission in string : ", fileInfo.Mode().String())

	fmt.Println("What else underneath? : ", fileInfo.Sys())

	return nil
}
