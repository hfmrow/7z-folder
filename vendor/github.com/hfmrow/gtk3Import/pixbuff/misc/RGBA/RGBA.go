// RGBA.go

/*
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2020 H.F.M - RGBA package github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT
	License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package RGBA

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/gotk3/gotk3/gdk"
)

// Rgba: Contains RGBA values to facilitate transfer between local
// storage intended to be backed up such as in an 'option' file
// and the GdkRGBA object.
type Rgba struct {
	Red,
	Green,
	Blue int
	Alpha float64
}

func RgbaNew(r, g, b int, a float64) *Rgba {
	return &Rgba{
		Red:   r,
		Green: g,
		Blue:  b,
		Alpha: a}
}

// ToGdkRGBA: Set values to a new GdkRGBA object
// 'options' means 1st = invert, 2nd = preserve alpha
func (rgba *Rgba) ToGdkRGBA(options ...bool) *gdk.RGBA {
	i := 0
	ia := 0.0
	if len(options) > 0 && options[0] {
		i = 255
	}
	if len(options) > 1 && !options[1] {
		ia = 1.0
	}
	r := math.Abs(float64(i-rgba.Red)) / 255.0
	g := math.Abs(float64(i-rgba.Green)) / 255.0
	b := math.Abs(float64(i-rgba.Blue)) / 255.0
	a := math.Abs(ia - rgba.Alpha)

	return gdk.NewRGBA(r, g, b, a)
}

// FromGdkRGBA: Set values from GdkRGBA object
func (rgba *Rgba) FromGdkRGBA(RGBA *gdk.RGBA) {
	rgba.Red = int(RGBA.GetRed() * 255)
	rgba.Green = int(RGBA.GetGreen() * 255)
	rgba.Blue = int(RGBA.GetBlue() * 255)
	rgba.Alpha = RGBA.GetAlpha()
}

// RGBAInvert: do a color video inversion preserving original GdkRBGA
// func RGBAInvert(rgba *gdk.RGBA, invert bool, preserveAlpha ...bool) *gdk.RGBA {

// 	var i, ia float64 = 1, 1
// 	if !invert {
// 		i = 0
// 		ia = 0
// 	}
// 	if len(preserveAlpha) > 0 && preserveAlpha[0] {
// 		ia = 0
// 	}
// 	return gdk.NewRGBA(
// 		math.Abs(i-rgba.GetRed()),
// 		math.Abs(i-rgba.GetGreen()),
// 		math.Abs(i-rgba.GetBlue()),
// 		math.Abs(ia-rgba.GetAlpha()))
// }

// RGBAParseFromCss: Convert CSS RGBA synbtax to GdkRGBA structure.
func RGBAParseFromCss(r, g, b int, a float64) *gdk.RGBA {
	return gdk.NewRGBA(float64(r)/255.0, float64(g)/255.0, float64(b)/255.0, a)
}

// RGBAParseToCss: Convert GdkRGBA to string that match required format
// by CSS syntax. '255, 255, 255, 0.5'. 'invert' indicate that result
// will be inversed.
func RGBAParseToCss(rgba *gdk.RGBA, invert ...bool) string {

	var inv float64
	if len(invert) > 0 && invert[0] {
		inv = 1
	}
	return fmt.Sprintf("%d, %d, %d, %1.3f",
		int(math.Abs(inv-rgba.Floats()[0])*255),
		int(math.Abs(inv-rgba.Floats()[1])*255),
		int(math.Abs(inv-rgba.Floats()[2])*255),
		math.Abs(inv-rgba.Floats()[3]))
}

// TODO Checkout for ALPHA chanel
// Convert gdk.RGBA to Hex string value: "#EF2929".
func ColGdkRGBA2Hex(value *gdk.RGBA) string {
	// Convert int string value to Hex with 2 digits.
	var build2DigitsHex = func(intValueStr string) string {
		xi, _ := strconv.Atoi(intValueStr)
		xs := fmt.Sprintf("%X", xi)
		if len(xs) == 1 {
			xs = "0" + xs
		}
		return xs
	}

	regOne := regexp.MustCompile(`[()]`)
	regTwo := regexp.MustCompile(`[,]`)
	tmpStrSl := regOne.Split(value.String(), -1)
	tmpStrSl = regTwo.Split(tmpStrSl[1], -1)
	rr := build2DigitsHex(tmpStrSl[0])
	gg := build2DigitsHex(tmpStrSl[1])
	bb := build2DigitsHex(tmpStrSl[2])
	aa := "FF"
	if len(tmpStrSl) > 3 {
		aa = build2DigitsHex(tmpStrSl[3])
	}
	return fmt.Sprintf("#%s%s%s%s", rr, gg, bb, aa)
}
