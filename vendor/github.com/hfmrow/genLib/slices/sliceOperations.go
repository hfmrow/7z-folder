// sliceOperations.go

package slices

import (
	"fmt"
	"log"
	"reflect"

	gst "github.com/hfmrow/genLib/strings"
)

// RemDupSlIface: Designed to accept all types. Remove duplicate entries
// Returns true if something has changed, false otherwise.
// NOTE: 'slice' argument MUST be given as pointer '&'
func RemDupSlIface(slice1, slice2 interface{}) bool {

	var idxStore []int
	switch s := slice1.(type) {
	case *[]int:
		for idx1, val1 := range *s {
			for _, val2 := range *slice2.(*[]int) {
				if val1 == val2 {
					idxStore = append(idxStore, idx1)
				}
			}
		}
		if len(idxStore) == 0 {
			return false
		}
		ss := *s
		for _, idx := range idxStore {
			ss = append(ss[:idx], ss[idx+1:]...)
		}
		*slice1.(*[]int) = ss
	case *[]string:
		for idx1, val1 := range *s {
			for _, val2 := range *slice2.(*[]string) {
				if val1 == val2 {
					idxStore = append(idxStore, idx1)
				}
			}
		}
		if len(idxStore) == 0 {
			return false
		}
		ss := *s
		for _, idx := range idxStore {
			ss = append(ss[:idx], ss[idx+1:]...)
		}
		*slice1.(*[]string) = ss
	default:
		log.Printf("IsExistSlIface: Type [%T] or [%T], is not actually handled.\nRemember that the arrays must be passed as &pointer\n", slice1, slice2)
	}
	return true
}

// IsExistSlIface: Same as above but designed to accept all types.
// Return the position whether it found or -1 if not.
func IsExistSlIface(slice interface{}, item interface{}) int {

	switch s := slice.(type) {
	case []int:
		for idx, val := range s {
			if val == item.(int) {
				return idx
			}
		}
	case []string:
		for idx, val := range s {
			if val == item.(string) {
				return idx
			}
		}
	default:
		log.Printf("IsExistSlIface: Type [%T], is not actually handled.\n", slice)
	}
	return -1
}

// DeleteSlIface: designed to accept all types. Return true on success.
// NOTE: 'slice' argument MUST be given as pointer '&'
func DeleteSlIface(slice interface{}, idx int) bool {

	var length int
	switch s := slice.(type) {
	case *[]int:
		length = len(*s)
		if length > idx {
			ss := *s
			ss = append(ss[:idx], ss[idx+1:]...)
			*slice.(*[]int) = ss
			return true
		}
	case *[]string:
		length = len(*s)
		if len(*s) > idx {
			ss := *s
			ss = append(ss[:idx], ss[idx+1:]...)
			*slice.(*[]string) = ss
			return true
		}
	default:
		log.Printf("DeleteSlIface: Type [%T], is not actually handled.\nRemember that the array must be passed as &pointer\n", slice)
		return false
	}
	log.Printf("DeleteSlIface: Targeted index [%d] is out of range, slice length [%d]\n", idx, length)
	return false
}

// GetStrIndex: Get index of a string in a slice, Return -1 if no entry found ...
func GetStrIndex(slice []string, item string) int {

	log.Printf("GetStrIndex: function out of date, must be replaced with [%s]\n", "IsExistSlIface")

	for idx, row := range slice {
		if row == item {
			return idx
		}
	}
	return -1
}

// IsExistSl: if exist then  ...
func IsExistSl(slice []string, item string) bool {

	log.Printf("IsExistSl: function out of date, must be replaced with [%s]\n", "IsExistSlIface")

	for _, mainRow := range slice {
		if mainRow == item {
			return true
		}
	}
	return false
}

// GetStrIndex2dCol: Search in 2d string slice if a column's value exist and return row number.
func GetStrIndex2dCol(slice [][]string, value string, col int) int {
	for idx, row := range slice {
		if row[col] == value {
			return idx
		}
	}
	return -1
}

// IsExist2dCol: Search in 2d string slice if a column's value exist.
func IsExist2dCol(slice [][]string, value string, col int) bool {
	for _, mainRow := range slice {
		if mainRow[col] == value {
			return true
		}
	}
	return false
}

// IsExist2d Search in 2d string slice if a row exist (deepequal row).
func IsExist2d(slice [][]string, cmpRow []string) bool {
	for _, mainRow := range slice {
		if reflect.DeepEqual(mainRow, cmpRow) {
			return true
		}
	}
	return false
}

// CmpRemSl2d: Compare slice1 and slice2 if row exist on both,
// the raw is removed from slice2 and result returned.
func CmpRemSl2d(sl1, sl2 [][]string) (outSlice [][]string) {
	var skip bool
	for _, row2 := range sl2 {
		for _, row1 := range sl1 {
			if reflect.DeepEqual(row1, row2) {
				skip = true
				break
			}
		}
		if !skip {
			outSlice = append(outSlice, row2)
		}
		skip = false
	}
	return outSlice
}

// CheckForDupSl2d: Return true if duplicated row is found
func CheckForDupSl2d(sl [][]string) bool {
	for idx := 0; idx < len(sl)-1; idx++ {
		for secIdx := idx + 1; secIdx < len(sl); secIdx++ {
			if reflect.DeepEqual(sl[idx], sl[secIdx]) {
				return true
			}
		}
	}
	return false
}

// RemoveDupSl: Remove duplicate entry in a string slice
func RemoveDupSl(slice []string) []string {

	log.Printf("GetStrIndex: function out of date, must be replaced with [%s]\n", "RemDupSlIface")

	var isExist bool
	tmpSlice := make([]string, 0)
	for _, inValue := range slice {
		isExist = false
		inValue = gst.RemoveNonAlNum(inValue)
		for _, outValue := range tmpSlice {
			if outValue == inValue {
				isExist = true
				break
			}
		}
		if !isExist {
			tmpSlice = append(tmpSlice, inValue)
		}
	}
	fmt.Println(len(tmpSlice))
	return tmpSlice
}

// RemoveDupSl2d: Remove duplicate entry in a 2d string slice based on column number content.
func RemoveDupSl2d(slice [][]string, col int) (outSlice [][]string) {
	if len(slice) == 0 {
		return slice
	}
	var dupFlag bool
	outSlice = append(outSlice, slice[0])
	for primIdx := len(slice) - 1; primIdx >= 0; primIdx-- {
		dupFlag = false
		for secIdx := len(outSlice) - 1; secIdx >= 0; secIdx-- {
			if outSlice[secIdx][col] == slice[primIdx][col] {
				dupFlag = true
			}
		}
		if !dupFlag {
			outSlice = append(outSlice, slice[primIdx])
		}
	}
	return outSlice
}

// Preppend: Add data at the begining of a string slice
func Preppend(slice []string, prepend ...string) []string {
	return append(prepend, slice...)
}

// AppendAt: Add data at a specified position in slice of a string
func AppendAt(slice []string, pos int, insert ...string) []string {
	if pos > len(slice) {
		pos = len(slice)
	} else if pos < 0 {
		pos = 0
	}
	return append(slice[:pos], append(insert, slice[pos:]...)...)
}
