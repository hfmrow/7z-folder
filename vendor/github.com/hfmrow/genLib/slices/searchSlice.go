// searchSlice.go

package slices

import (
	"errors"
	"fmt"
	"regexp"
)

// SearchSl: Search in 2d string slice. cs=case sensitive, ww=whole word, rx=regex
func SearchSl(find string, table [][]string, caseSensitive, POSIXcharClass, POSIXstrictMode, regex, wholeWord bool) (out [][]string, err error) {
	if len(table) != 0 {
		if len(find) != 0 {
			search := find
			var outTable [][]string
			// if POSIXcharClass { // Commented to avoid library cycle ...
			// 	search = str.StringToCharacterClasses(find, caseSensitive, POSIXstrictMode)
			// }
			search = `(` + search + `)`
			if wholeWord {
				search = `\b` + search + `\b`
			}
			if !caseSensitive && !POSIXcharClass {
				search = `(?i)` + search
			}
			regX, err := regexp.Compile(search)
			if err != nil {
				return out, err
			}
			fmt.Println("regex: " + search)
			if regex {
				regX, err = regexp.Compile(find)
				if err != nil {
					return out, err
				}
			}

			for idxRow := 0; idxRow < len(table); idxRow++ {
				for _, col := range table[idxRow] {
					if regX.MatchString(col) {
						outTable = append(outTable, table[idxRow])
						break // Avoid duplicate when element found twice in same row
					}
				}
			}
			if len(outTable) == 0 {
				return out, errors.New(find + "\n\nNot found ...")
			}
			return outTable, nil // Result found then return it
		}
	}
	return [][]string{}, errors.New("Nothing to search ...")
}
