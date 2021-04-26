// bench.go

/*
	Copyright ©2018-19 H.F.M
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	Measuring lapse (may be multiples) between operations.

Usage:
		bench := BenchNew(true) // true means we want display results at Stop().
		bench.Lapse("first lapse")
		// Doing something ...
		bench.Lapse("Nth lapse")
		// Doing another thing ...
		bench.Stop                 // Display results
*/

package bench

import (
	"fmt"
	"strings"
	"time"
)

type Bench struct {
	Lapses  []lapse
	Average lapse
	Total   lapse

	lbl     string
	display bool
}

// BenchNew: default showResult = true
func BenchNew(showResults ...bool) (bench *Bench) {
	bench = new(Bench)
	bench.display = true
	if len(showResults) > 0 {
		bench.display = showResults[0]
	}
	return
}
func (b *Bench) Start(label ...string) {
	b.Lapse(label...)
}
func (b *Bench) Lapse(label ...string) {
	b.lbl = "Start"

	if len(label) == 0 {
		b.lbl = fmt.Sprintf("%d", len(b.Lapses))
	} else {
		b.lbl = strings.Join(label, " ")
	}

	b.Lapses = append(b.Lapses, lapse{
		Time:  time.Now(),
		Label: b.lbl})
}

func (b *Bench) NanoConv(nano int64) (min, sec, ms, ns int64) {
	min = nano / 60000000000
	minr := nano % 60000000000
	sec = minr / 1000000000
	secr := minr % 1000000000
	ms = secr / 1000000
	ns = secr % 1000000
	return
}

func (b *Bench) Stop() {
	b.Lapse()
	var means, nano int64
	nextLbl, tmpNextLbl := b.Lapses[0].Label, ""
	for t := 1; t < len(b.Lapses); t++ { // Lapses calculation
		tmpNextLbl = b.Lapses[t].Label //
		b.Lapses[t].Label = nextLbl    // Switching labels
		nextLbl = tmpNextLbl           //
		nano = b.Lapses[t].Time.Sub(b.Lapses[t-1].Time).Nanoseconds()
		means += nano
		b.Lapses[t].toLapse(nano)
	}

	b.Lapses = b.Lapses[1:]
	if len(b.Lapses)-1 > 0 { // Average calculation

		b.Average.toLapse(means / int64(len(b.Lapses)))
		b.Average.Label = "Average"

		b.Total.toLapse(means)
		b.Total.Label = "Total"
	}

	if b.display {
		b.disp()
	}
}

func (b *Bench) Reset() {
	b.Lapses = b.Lapses[:0]
}

func (b *Bench) disp() {
	if len(b.Lapses) > 1 {
		for _, lps := range b.Lapses {
			fmt.Printf("%s: %s\n", lps.Label, lps.String)
		}
	} else {
		fmt.Printf("%s\n", b.Lapses[0].String)
	}
}

type lapse struct {
	Time                       time.Time
	Elapsed                    time.Time
	Min, Sec, Ms, Ns           int64
	StringShort, String, Label string
}

func (l *lapse) toLapse(nano int64) {
	b := new(Bench)
	l.Elapsed = time.Unix(0, nano)
	l.Min, l.Sec, l.Ms, l.Ns = b.NanoConv(nano)
	l.StringShort = fmt.Sprintf("%dm %d.%ds", l.Min, l.Sec, l.Ms)
	l.String = fmt.Sprintf("%dm %ds %dms %dns", l.Min, l.Sec, l.Ms, l.Ns)
}
