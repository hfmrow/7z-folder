// strNum.go

/*
	©2019 H.F.M. MIT license
*/

package main

import (
	"fmt"
	"regexp"
	"strings"
)

// // BaseNoExt: get only the name without ext.
// func BaseNoExt(filename string) (outFilename string) {
// 	outFilename = filepath.Base(filename)
// 	ext := filepath.Ext(outFilename)
// 	return outFilename[:len(outFilename)-len(ext)]
// }

type WordWithDigit struct {
	maxLength, maxLengthLeft int
	zeroMask                 string
	ForceRightDigit          bool
}

func (w *WordWithDigit) Init(words []string) {
	for _, word := range words {
		if len(word) > w.maxLength {
			w.maxLength = len(word)
			if digitsPosition(word) == 0 {
				digits := getDigits(word)
				if len(digits) > w.maxLengthLeft {
					w.maxLengthLeft = len(digits)
				}
			}
		}
	}
}

// FillWordToMatchMaxLength: Convert word(s) into numbered one, like "label1" -> "label000001" etc...
// results are based on list of words that determine max length of them to determinate
// the final length of the modified word. This is used in case of sorting list
// of words that contains numeric value to avoid the disorder result
// like "1label", "10label", "2label" etc ...
func (w *WordWithDigit) FillWordToMatchMaxLength(inString string, removeExt ...bool) (outString string) {
	var word string
	if len(removeExt) != 0 {
		if removeExt[0] {
			inString = BaseNoExt(inString)
		}
	}
	inString = strings.ToLower(strings.TrimSpace(inString))
	zeroCount := w.maxLength - len(inString)
	for idx := 0; idx < zeroCount; idx++ {
		w.zeroMask += "0"
	}
	wordPos := digitsPosition(inString)
	digits := getDigits(inString)
	switch wordPos {
	case 0: // Left
		word = inString[len(digits):len(inString)]
		outString = word + w.zeroMask + digits
	case 1: // Right
		word = inString[:len(inString)-len(digits)]
		outString = word + w.zeroMask + digits
	case -1: // None
		outString = inString + w.zeroMask
	}
	w.zeroMask = ""
	return outString
}

// numPosition: detect position of digit part: 0=left, 1=right, -1=none
func digitsPosition(inString string) int {
	digitS := regexp.MustCompile(`^(\d)`)
	digitE := regexp.MustCompile(`(\d)$`)
	switch {
	case digitS.MatchString(inString):
		return 0 // Left
	case digitE.MatchString(inString):
		return 1 // Right
	}
	return -1 // None
}

// getDigits: return digit part of string prior at start or at end, -1=none
func getDigits(inString string) (value string) {
	digitS := regexp.MustCompile("(^[0-9]*)")
	digitE := regexp.MustCompile("([0-9]*$)")
	start := digitS.FindString(inString)
	end := digitE.FindString(inString)
	switch {
	case len(start) != 0: // Left
		value = start
	case len(end) != 0: // Right
		value = end
	}
	return value
}

// Increment: Add incrementation to filename
func Increment(inSL []string, separator string, startAt int, atLeft ...bool) (outSL []string) {
	var tmpStr string
	var left bool
	if len(atLeft) != 0 {
		left = atLeft[0]
	}
	digits := fmt.Sprintf("%d", len(inSL)+startAt)
	regRepl := regexp.MustCompile(`[[:digit:]]`)
	digits = regRepl.ReplaceAllString(digits, "0")
	for idx, line := range inSL {
		inc := fmt.Sprintf("%d", idx+startAt)
		inc = digits[len(inc):] + inc
		switch left {
		case true:
			tmpStr = inc + separator + line
		case false:
			tmpStr = line + separator + inc
		}
		outSL = append(outSL, tmpStr)
	}
	return outSL
}
