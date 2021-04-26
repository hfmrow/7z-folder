// stringOperations.go

/*
	Part of H.F.M personal GenLib p	akage (not published)
	Copyright ©2020 H.F.M github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package strings

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type HtmlEscOpt int

const (
	HE_UNESCAPE HtmlEscOpt = 1 << iota
	HE_PANGO
)

// HtmlEscOrUnEsc: Escape or unescape to/from html, usefull for pango too.
func HtmlEscOrUnEsc(in string, options ...HtmlEscOpt) string {

	var (
		opt         HtmlEscOpt
		unesc       bool
		htmlEscaper = map[string]string{
			`&`: "&amp;",
			`<`: "&lt;",
			`>`: "&gt;",
		}
	)
	if len(options) > 0 {
		opt = options[0]
	}
	if opt&HE_PANGO == 0 {
		htmlEscaper[`'`] = "&#39;" // "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
		htmlEscaper[`"`] = "&#34;" // "&#34;" is shorter than "&quot;".
	}
	unesc = opt&HE_UNESCAPE != 0

	for unes, escp := range htmlEscaper {
		switch unesc {
		case true:
			in = strings.ReplaceAll(in, escp, unes)
		case false:
			in = strings.ReplaceAll(in, unes, escp)
		}
	}
	return in
}

// UnescapeToUtf8: Convert raw string that contain escaped values to
// utf-8 literal string.
func UnescapeToUtf8(inStr string) string {

	var escRawConv = map[string]string{
		`\n`: "\u000A", // line feed or newline
		`\t`: "\u0009", // horizontal tab
		`\"`: "\u0022", // double quote
		`\r`: "\u000D", // carriage return
		`\\`: "\u005c", // backslash
		`\f`: "\u000C", // form feed
		`\v`: "\u000b", // vertical tab
		`\b`: "\u0008", // backspace
		`\a`: "\u0007", // alert or bell (little joke)
	}
	for esc, uni := range escRawConv {
		inStr = strings.ReplaceAll(inStr, esc, uni)
	}

	return inStr
}

// UnEscapedStr: Convert raw string that contain escaped values to
// literal string. Same result as below.
func UnEscapedStr(in string) string {
	return fmt.Sprint(UnescapeToUtf8(in))
}

// UnEscapeString: Convert raw string that contain escaped values to
// literal string.
func UnEscapeString(inString string) string {
	inString = strings.ReplaceAll(inString, `"`, `\"`)
	outString, err := strconv.Unquote(`"` + inString + `"`)
	if err != nil {
		log.Fatalf("UnEscapeString/Unquote: %v [%v]\n", inString, err)
	}
	return outString
}

// LowercaseAtFirst: true if 1st char is lowercase
func LowercaseAtFirst(inString string) bool {

	if len(inString) != 0 {
		charType, _ := regexp.Compile("[[:lower:]]")
		return charType.MatchString(inString[:1])
	}
	return true
}

// toCamel: Turn string into camelCase or PascalCase style (Go).
func ToCamel(inString string, lowerAtFirst ...bool) string {

	lAtFirst := false
	regNonAlNum := regexp.MustCompile("[^[:alnum:]]+")

	if len(lowerAtFirst) > 0 && lowerAtFirst[0] {
		lAtFirst = lowerAtFirst[0]
	}
	tmpString := regNonAlNum.Split(SeparateUpper(inString, " "), -1)
	sl := make([]string, len(tmpString))

	for idx, word := range tmpString {
		if lAtFirst && idx < 1 {
			sl[idx] = strings.ToLower(word)
		} else {
			sl[idx] = strings.Title(word)
		}
	}
	return strings.Join(sl, "")
}

// toSnake: Turn string into snake_case style (C).
func ToSnake(inString string, allCaps ...bool) string {
	return toSnakeOrKebab(inString, "_", allCaps...)
}

// toKebab: Turn string into kebab-case style (URLs).
func ToKebab(inString string, allCaps ...bool) string {
	return toSnakeOrKebab(inString, "-", allCaps...)
}

// toSnakeOrKebab:
func toSnakeOrKebab(inString, sep string, allCaps ...bool) string {

	allC := false
	regNonAlNum := regexp.MustCompile("[^[:alnum:]]+")

	if len(allCaps) > 0 {
		allC = allCaps[0]
	}
	tmpString := regNonAlNum.Split(SeparateUpper(strings.TrimSpace(inString), " "), -1)
	var sl []string
	for _, word := range tmpString {
		if len(word) > 0 {
			if !allC {
				sl = append(sl, strings.ToLower(word))
			} else {
				sl = append(sl, strings.ToUpper(word))
			}
		}
	}
	return strings.Join(sl, sep)
}

// SeparateUpper: Add a 'sep' before each upper case char
// except the first.
func SeparateUpper(inString, sep string) string {

	var words []string
	l := 0
	for s := inString; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		words = append(words, s[:l])
	}
	if len(words) > 0 {
		return strings.Join(words, sep)
	}
	return inString
}

// RemoveNonAlNum: Remove all non alpha-numeric char
func RemoveNonAlNum(inString string) string {
	nonAlNum := regexp.MustCompile(`[[:punct:]]`)
	return nonAlNum.ReplaceAllString(inString, "")
}

// RemoveNonNum: Remove all non numeric char
func RemoveNonNum(inString string) string {
	nonAlNum := regexp.MustCompile(`[[:punct:][:alpha:]]`)
	return nonAlNum.ReplaceAllString(inString, "")
}

// ReplaceSpace: replace all [[:space::]] with 'repl'
func ReplaceSpace(inString, repl string) string {
	spaceRegex := regexp.MustCompile(`[[:space:]]`)
	return spaceRegex.ReplaceAllString(inString, repl)
}

// RemoveDupSpace: Remove duplicated space/tab in string
func RemoveDupSpace(inString string) string {
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`) //	to match 2 or more whitespace symbols inside a string
	return strings.TrimSpace(remInside.ReplaceAllString(inString, " "))
}

// RemoveSpace: remove all [[:space::]]
func RemoveSpace(inString string) string {
	spaceRegex := regexp.MustCompile(`[[:space:]]`)
	return spaceRegex.ReplaceAllString(inString, "")
}

// ReplacePunct: replace all [[:punct::]] with 'repl'
func ReplacePunct(inString, repl string) string {
	spaceRegex := regexp.MustCompile(`[[:punct:]]`)
	return spaceRegex.ReplaceAllString(inString, repl)
}

// SplitNumeric: Split and keep all numeric values in a string
func SplitNumeric(inString string) (outText []string, err error) {
	toSplit := regexp.MustCompile(`[[:alpha:][:punct:]]`)
	spaceSepared := toSplit.ReplaceAllString(inString, " ")
	outText = strings.Split(RemoveDupSpace(spaceSepared), " ")
	return outText, err
}

// // TrimSpace: Some multiple way to trim strings. cmds is optionnal or accept multiples args
// func TrimSpace(inputString string, cmds ...string) (newstring string, err error) {

// 	osForbiden := regexp.MustCompile(`[<>:"/\\|?*]`)
// 	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`)    //	to match 2 or more whitespace symbols inside a string
// 	remInsideNoTab := regexp.MustCompile(`[\p{Zs}]{2,}`) //	(preserve \t) to match 2 or more space symbols inside a string

// 	if len(cmds) != 0 {
// 		for _, command := range cmds {
// 			switch command {
// 			case "+h": //	Escape html
// 				inputString = html.EscapeString(inputString)
// 			case "-h": //	UnEscape html
// 				inputString = html.UnescapeString(inputString)
// 			case "+e": //	Escape specials chars
// 				inputString = fmt.Sprintf("%q", inputString)
// 			case "-e": //	Un-Escape specials chars
// 				tmpString, err := strconv.Unquote(`"` + inputString + `"`)
// 				if err != nil {
// 					return inputString, err
// 				}
// 				inputString = tmpString
// 			case "-w": //	Change all illegals chars (for path in linux and windows) into "-"
// 				inputString = osForbiden.ReplaceAllString(inputString, "-")
// 			case "+w": //	clean all illegals chars (for path in linux and windows)
// 				inputString = osForbiden.ReplaceAllString(inputString, "")
// 			case "-c": //	Trim [[:space:]] and clean multi [[:space:]] inside
// 				inputString = strings.TrimSpace(remInside.ReplaceAllString(inputString, " "))
// 			case "-ct": //	Trim [[:space:]] and clean multi [[:space:]] inside (preserve TAB)
// 				inputString = strings.Trim(remInsideNoTab.ReplaceAllString(inputString, " "), " ")
// 			case "-s": //	To match 2 or more whitespace leading/ending/inside a string (include \t, \n)
// 				inputString = strings.Join(strings.Fields(inputString), " ")
// 			case "-&": //	Replace ampersand CHAR with ampersand HTML code
// 				inputString = strings.Replace(inputString, "&", "&amp;", -1)
// 			case "+&": //	Replace ampersand HTML code with ampersand CHAR
// 				inputString = strings.Replace(inputString, "&amp;", "&", -1)
// 			default:
// 				return inputString, errors.New("TrimSpace, " + command + ", does not exist")
// 			}
// 		}
// 	}
// 	return inputString, nil
// }
