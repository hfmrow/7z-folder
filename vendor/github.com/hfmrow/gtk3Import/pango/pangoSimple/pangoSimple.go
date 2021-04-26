// PangoMarkupBinder.go

/*
*	©2019 H.F.M. MIT license
*	Handle Pango markup functions.

*	This is a pango markup binder to gotk3 pango library ...
*	I have made it to make more simple markup handling when i'm working on gtk objects.
*	It can be used in treeview, dialog, label, ... Each object where you can set "markup" content

 */

package gtk3Import

import (
	"fmt"
	"regexp"
	"strings"
)

type PangoSimple struct {
	pangoEscapeChar [][]string
}

func PangoSimpleNew() (ps *PangoSimple) {
	ps = new(PangoSimple)
	ps.Init()
	return
}

// Init: Init [PangoSimple] structure
func (pm *PangoSimple) Init() {
	// pm.pangoEscapeChar = [][]string{{"<", "&lt;", string([]byte{0x15})}, {"&", "&amp;", string([]byte{0x16})}}
	pm.pangoEscapeChar = [][]string{{"<", "&lt;", "lOwErThAnTmPrEpLaCeMeNt"}, {"&", "&amp;", "aMpErSaNdTmPrEpLaCeMeNt"}}
}

// MarkupHttpClickable: Search for http adresses in the input text and make http adress as clickable links
func (pm *PangoSimple) MarkupHttpClickable(inString string, keepFromEnd ...int) (outString string) {
	var adrToDisp string
	outString = pm.Prepare(inString) // In case nothing is found, returned value is same as entered value
	reg := regexp.MustCompile(`(http|https|ftp|ftps)\:\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(\/\S*)?`)
	indexes := reg.FindAllIndex([]byte(outString), -1)

	removePrefix := regexp.MustCompile(`(http|https|ftp|ftps)\:\/\/`)

	for idx := len(indexes) - 1; idx >= 0; idx-- {
		inLeft := outString[:indexes[idx][0]]
		inRight := outString[indexes[idx][1]:]
		url := outString[indexes[idx][0]:indexes[idx][1]]

		adrToDisp = removePrefix.ReplaceAllString(url, "")
		if len(keepFromEnd) > 0 {
			tmpUrl := strings.Split(url, "/")
			adrToDisp = strings.Join(tmpUrl[len(tmpUrl)-keepFromEnd[0]:], "/")
		}

		outString = inLeft + `<a href="` + url + `">` + adrToDisp + `</a>` + inRight
	}
	return pm.Finalize(outString)
}

// ApplyMarkup: This method add automatically the closer at end.
// ie: ApplyMarkup(`<span foreground="#FFFFFF">`, "hello")
func (pm *PangoSimple) ApplyMarkup(markup ...string) string {
	return pm.Finalize(markup[0] + pm.Prepare(markup[1]) + pm.getClose(markup[0]))
}

func (pm *PangoSimple) ApplyMarkupAgain(markup ...string) string {
	return markup[0] + markup[1] + pm.getClose(markup[0])
}

// getClose: return the last argument derived from the first to close the pango string
func (pm *PangoSimple) getClose(in string) string {
	if strings.HasPrefix(in, `<span `) {
		return `</span>`
	}
	if strings.HasPrefix(in, `<a `) {
		return `</a>`
	}
	return fmt.Sprintf("</%s>", strings.Split(in, `>`)[0][1:])
}

// prepare: sanitize input string to safely use with pango
func (pm *PangoSimple) Prepare(inString string) string {
	inString = strings.ReplaceAll(inString, pm.pangoEscapeChar[1][0], pm.pangoEscapeChar[1][2])
	return strings.ReplaceAll(inString, pm.pangoEscapeChar[0][0], pm.pangoEscapeChar[0][2])
}

// finalize: restore originals characters using markup replacement
func (pm *PangoSimple) Finalize(inString string) string {
	inString = strings.ReplaceAll(inString, pm.pangoEscapeChar[1][2], pm.pangoEscapeChar[1][1])
	return strings.ReplaceAll(inString, pm.pangoEscapeChar[0][2], pm.pangoEscapeChar[0][1])
}

type PangoColor struct {
	Black     string
	Brown     string
	White     string
	Red       string
	Green     string
	Blue      string
	Cyan      string
	Magenta   string
	Purple    string
	Turquoise string
	Violet    string

	Darkred   string
	Darkgreen string
	Darkblue  string
	Darkgray  string
	Darkcyan  string

	Lightblue      string
	Lightgray      string
	Lightgreen     string
	Lightturquoise string
	Lightred       string
	Lightyellow    string
}

func PangoColorNew() (pc *PangoColor) {
	pc = new(PangoColor)
	pc.Init()
	return
}
func (pc *PangoColor) Init() {
	//	Colors initialisation
	pc.Black = "#000000"
	pc.Brown = "#7C2020"
	pc.White = "#FFFFFF"
	pc.Red = "#FF2222"
	pc.Green = "#22BB22"
	pc.Blue = "#0044FF"
	pc.Cyan = "#14FFFA"
	pc.Magenta = "#D72D6C"
	pc.Purple = "#8B0037"
	pc.Turquoise = "#009187"
	pc.Violet = "#7F00FF"
	pc.Darkred = "#300000"
	pc.Darkgreen = "#003000"
	pc.Darkblue = "#000030"
	pc.Darkcyan = "#003333"
	pc.Darkgray = "#303030"
	pc.Lightturquoise = "#80FFE7"
	pc.Lightblue = "#ADD8E6"
	pc.Lightgray = "#E4DDDD"
	pc.Lightgreen = "#87FF87"
	pc.Lightred = "#FF6666"
	pc.Lightyellow = "#FFFF6F"
}
