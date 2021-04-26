// unit-hr.go

/*
	Copyright ©2020 H.F.M - Unit human readable v1.0 https://github.com/hfmrow
	This library convert raw units to human readable version.
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package unit_hr

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

// options: UNIT_HIDE, UNIT_SHORTEN, UNIT_DECIMAL, UNIT_LOWER
type HROptions int

const (
	UNIT_DEFAULT HROptions = 1 << 0
	UNIT_HIDE    HROptions = 1 << 1
	UNIT_SHORTEN HROptions = 1 << 2
	UNIT_DECIMAL HROptions = 1 << 3
	UNIT_LOWER   HROptions = 1 << 4
)

// humanReadableSize: Convert file size (octets/bytes) to human readable
// version. Note: 'options' in order means => 'useDecimal', 'hideUnit'.
// 'useDecimal' argument define kilo = 1000 instead of 1024.
func HumanReadableSize(size interface{}, opt ...HROptions) string {

	var (
		options                HROptions
		hideUnit               bool
		val                    float64
		kilo                   float64 = 1024
		sP, sT, sG, sM, sK, sb string  = "PiB", "TiB", "GiB", "MiB", "KiB", "B"
	)
	if len(opt) > 0 {
		options = opt[0]
	}
	if options&UNIT_DECIMAL != 0 {
		options = options | UNIT_SHORTEN
		kilo = 1000
	}
	if options&UNIT_SHORTEN != 0 {
		sP, sT, sG, sM, sK, sb = "PB", "TB", "GB", "MB", "kB", "B"
	}
	if options&UNIT_HIDE != 0 {
		hideUnit = true
	}
	switch v := size.(type) {
	case uint64:
		val = float64(v)
	case uint32:
		val = float64(v)
	case uint:
		val = float64(v)
	case int64:
		val = float64(v)
	case int32:
		val = float64(v)
	case int:
		val = float64(v)
	case float64:
		val = float64(v)
	case float32:
		val = float64(v)
	default:
		log.Printf("Unable to define type of: %v\n", size)
	}
	unit := sb
	switch {
	case val < kilo:
		val = val
		return fmt.Sprintf("%.0f%s", val, unit)
	case val < math.Pow(kilo, 2):
		val = val / kilo
		unit = sK
	case val < math.Pow(kilo, 3):
		val = val / math.Pow(kilo, 2)
		unit = sM
	case val < math.Pow(kilo, 4):
		val = val / math.Pow(kilo, 3)
		unit = sG
	case val < math.Pow(kilo, 5):
		val = val / math.Pow(kilo, 4)
		unit = sT
	case val < math.Pow(kilo, 6):
		val = val / math.Pow(kilo, 5)
		unit = sP
	}
	if options&UNIT_LOWER != 0 {
		unit = strings.ToLower(unit)
	}
	if hideUnit {
		unit = ""
	}
	return fmt.Sprintf("%2.2f%s", val, unit)
}

// HumanReadableSizeAlt: Same as above with an alternative method.
// Format byte size to human readable format.
func HumanReadableSizeAlt(b int64, options ...HROptions) string {
	var (
		unit    int64 = 1024
		suffix        = "iB"
		unitStr       = "KMGTPE"
	)
	if len(options) > 0 {
		opt := options[0]
		if opt&UNIT_DECIMAL != 0 {
			unit = 1000
		}
		if opt&UNIT_SHORTEN != 0 || opt&UNIT_DECIMAL != 0 {
			suffix = "B"
		}
		if opt&UNIT_LOWER != 0 {
			suffix = strings.ToLower(suffix)
			unitStr = strings.ToLower(unitStr)
		}
		if opt&UNIT_HIDE != 0 {
			suffix = ""
			unitStr = ""
		}
	}
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%s%c%s", float64(b)/float64(div), " ", unitStr[exp], suffix)
}

// String: Convert frenquency value int64 to string human readable version.
func HumanReadableFreq(value interface{}, hideUnit ...bool) string {

	var (
		val  float64
		unit string
	)

	switch v := value.(type) {
	case int64:
		val = float64(v)
	case int32:
		val = float64(v)
	case int:
		val = float64(v)
	default:
		val = value.(float64)
	}

	switch {
	case val < 1000:
		if len(hideUnit) > 0 && hideUnit[0] {
			unit = ""
		} else {
			unit = "hz"
		}
		fmt.Sprintf("%2.0f%s", val, unit)
	case val < 1000e3:
		val = val / 1000
		unit = "Mhz"
	case val < float64(1000e6):
		val = val / 1000e3
		unit = "Ghz"
	case val < float64(1000e9):
		val = val / 1000e6
		unit = "Thz"
	}

	if len(hideUnit) > 0 && hideUnit[0] {
		unit = ""
	}
	return fmt.Sprintf("%2.2f%s", val, unit)
}

/*
 * form: https://stackoverflow.com/questions/36530251/time-since-with-months-and-years
 */
func DaysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func Elapsed(from, to time.Time) (inverted bool, years, months, days, hours, minutes, seconds, nanoseconds int) {
	if from.Location() != to.Location() {
		to = to.In(to.Location())
	}

	inverted = false
	if from.After(to) {
		inverted = true
		from, to = to, from
	}

	y1, M1, d1 := from.Date()
	y2, M2, d2 := to.Date()

	h1, m1, s1 := from.Clock()
	h2, m2, s2 := to.Clock()

	ns1, ns2 := from.Nanosecond(), to.Nanosecond()

	years = y2 - y1
	months = int(M2 - M1)
	days = d2 - d1

	hours = h2 - h1
	minutes = m2 - m1
	seconds = s2 - s1
	nanoseconds = ns2 - ns1

	if nanoseconds < 0 {
		nanoseconds += 1e9
		seconds--
	}
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	if days < 0 {
		days += DaysIn(y2, M2-1)
		months--
	}
	if months < 0 {
		months += 12
		years--
	}
	return
}

// TODO
// humanReadableDate: Convert time.Duration to string human readable version.
func HumanReadableDate(duration time.Duration, options ...string) string {

	var (
		out  string
		min  = time.Second * 60
		hour = min * 60
		day  = hour * 24
		week = day * 7
		d    int64
	)

	fmt.Println("DONT'T USE !!! not usable")

	switch {

	case duration >= week: // week
		d = int64(duration / week)
		out += fmt.Sprintf("%d week ago", d)
		duration = duration - (time.Duration(d) * (week))

	case duration >= day: // day
		d = int64(duration / day)
		out += fmt.Sprintf("%d day ago", d)
		duration = duration - (time.Duration(d) * (day))

	case duration >= hour: // hour
		d = int64(duration / hour)
		out += fmt.Sprintf("%d hour ago", d)
		duration = duration - (time.Duration(d) * (hour))

	case duration >= min: // minute
		d = int64(duration / min)
		out += fmt.Sprintf("%d minute ago", d)
		duration = duration - (time.Duration(d) * (min))

	case int64(duration) > 0: // seconds
		out += fmt.Sprintf("%06.3fs", duration.Seconds())
	}

	return out
}

// humanReadableTime: Convert time.Duration to string human readable version.
func HumanReadableTime(duration time.Duration, unitSep ...string) string {

	var (
		out,
		sep string
		stop bool
		min  = time.Second * 60
		hour = min * 60
		day  = hour * 24
		week = day * 7
		d    int64
	)

	if len(unitSep) > 0 {
		sep = unitSep[0]
	}

	for !stop {
		switch {

		case duration >= week: // week
			d = int64(duration / week)
			out += fmt.Sprintf("%dw%s", d, sep)
			stop = duration == week
			duration = duration - (time.Duration(d) * (week))

		case duration >= day: // day
			d = int64(duration / day)
			out += fmt.Sprintf("%dd%s", d, sep)
			stop = duration == day
			duration = duration - (time.Duration(d) * (day))

		case duration >= hour: // hour
			d = int64(duration / hour)
			out += fmt.Sprintf("%dh%s", d, sep)
			stop = duration == hour
			duration = duration - (time.Duration(d) * (hour))

		case duration >= min: // minute
			d = int64(duration / min)
			out += fmt.Sprintf("%dm%s", d, sep)
			stop = duration == min
			duration = duration - (time.Duration(d) * (min))

		case int64(duration) > 0: // seconds
			out += fmt.Sprintf("%06.3fs%s", duration.Seconds(), sep)
			stop = true

		default:
			stop = true
		}
	}
	return out
}

// humanReadableTime: Convert bytes value to unit/sec string using 'delta'
// to specify the time interval. 'ext' could be 'iB' -> "KiB/s" or B -> "KB/s"
func HumanReadableBandwidth(bytes uint64, delta time.Duration, ext string) string {

	var (
		tmp       = float64(bytes) / delta.Seconds()
		dataBytes = uint64(tmp)
		unit      string
	)

	switch {
	case dataBytes < uint64(2<<9):
		return fmt.Sprintf("%.0f%sB/s", tmp, unit)

	case dataBytes < uint64(2<<19):
		tmp = tmp / float64(2<<9)
		unit = "K"

	case dataBytes < uint64(2<<29):
		tmp = tmp / float64(2<<19)
		unit = "M"

	case dataBytes < uint64(2<<39):
		tmp = tmp / float64(2<<29)
		unit = "G"

	case dataBytes < uint64(2<<49):
		tmp = tmp / float64(2<<39)
		unit = "T"
	}
	return fmt.Sprintf("%2.2f%s%s/s", tmp, unit, ext)
}
