// formatText.go

package strings

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// An unwanted behavior may occur on string where word's length > max...
func FormatTextQuoteBlankLines(str string) (out string) {
	outStrings := strings.Split(str, GetTextEOL([]byte(str)))
	for _, line := range outStrings {
		if len(line) != 0 {
			out += strings.TrimSuffix(strings.TrimPrefix(strconv.Quote(line+"\n"), `"`), `"`)
		} else {
			out += `\n`
		}
	}
	return `"` + out + `"`
}

// FormatTextParagraphFormatting:
func FormatTextParagraphFormatting(str string, max int, truncateWordOverMaxSize bool, indentFirstLinetr ...string) string {
	var tmpParag, tmpFinal, indent string
	eol := GetTextEOL([]byte(str))
	var addLine = func() {
		tmpFinal += FormatText(tmpParag, max, truncateWordOverMaxSize, indent)
	}
	if len(indentFirstLinetr) != 0 {
		indent = indentFirstLinetr[0]
	}
	lines := strings.Split(str, eol)
	for _, line := range lines {
		if len(line) != 0 {
			tmpParag += line
		} else {
			addLine()
			tmpFinal += eol + eol
			tmpParag = tmpParag[:0]
		}
	}
	addLine()
	return strings.TrimSuffix(tmpFinal, eol)
}

// FormatText: Format words text to fit (column/windows with limited width) "max"
// chars. An unwanted behavior may occur on string where word's length > max...
func FormatText(str string, max int, truncateWordOverMaxSize bool, indenT ...string) string {
	var indent string
	if len(indenT) != 0 {
		indent = indenT[0]
	}
	if len(str) > 0 {
		eol := GetTextEOL([]byte(str))
		var tmpLines, tmpWordTooLong []string
		space := regexp.MustCompile(`[[:space:]]`)
		var outText []string
		var countChar, length int
		// in case where string does not contain any [[:space:]] char
		var cutLongString = func(inStr string, inMax int, truncate bool) (outSliceText []string) {
			if !truncate {
				length = len(inStr)
				for count := 0; count < length; count = count + inMax {
					if count+inMax < length {
						outSliceText = append(outSliceText, inStr[count:count+inMax])
					} else {
						outSliceText = append(outSliceText, inStr[count:length])
					}
				}
			} else {
				outSliceText = append(outSliceText, TruncateString(inStr, "...", inMax, 1))
			}
			return outSliceText
		}
		text := space.Split(str, -1) // Split str at each blank char
		for idxWord := 0; idxWord < len(text); idxWord++ {
			length = len(text[idxWord]) + 1
			if length >= max {
				tmpWordTooLong = cutLongString(text[idxWord], max-1, truncateWordOverMaxSize)
				text = append(text[:idxWord], text[idxWord+1:]...) // Remove slice entry
				// Insert slices entries
				text = append(text[:idxWord], append(tmpWordTooLong, text[idxWord:]...)...)
				length = len(text[idxWord]) // Calculate new length
			}

			if countChar+length <= max {
				tmpLines = append(tmpLines, text[idxWord])
				countChar += length
			} else {
				outText = append(outText, indent+strings.Join(tmpLines, " "))
				tmpLines = tmpLines[:0] // Clear slice
				countChar = 0
				idxWord--
			}
		}
		// Get the rest of the text.
		outText = append(outText, indent+strings.Join(tmpLines, " "))
		return strings.Join(outText, eol)
	}
	return ""
}

// TruncateString: Reduce string length for display (prefix is separator like: "...",
// option=0 -> put separator at the begining of output string. Option=1 -> center,
// is where separation is placed. option=2 -> line feed, trunc the whole string
// using LF without shorting it. Max, is max char length of the output string.
func TruncateString(inString, prefix string, max, option int) string {
	var center, cutAt bool
	var outText string
	switch option {
	case 1:
		center = true
		cutAt = false
		max = max - len(prefix)
	case 2:
		center = false
		cutAt = true
	default:
		center = false
		cutAt = false
		max = max - len(prefix)
	}
	length := len(inString)
	if length > max {
		if cutAt {
			for count := 0; count < length; count = count + max {
				if count+max < length {
					outText += fmt.Sprintln(inString[count : count+max])
				} else {
					outText += fmt.Sprintln(inString[count:length])
				}
			}
			return outText
		} else if center {
			midLength := max / 2
			inString = inString[:midLength] + prefix + inString[length-midLength-1:]
		} else {
			inString = prefix + inString[length-max:]
		}
	}
	return inString
}
