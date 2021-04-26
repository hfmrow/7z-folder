// statusBar.go

/*
	This library use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019 H.F.M - Handle Statusbar messages & Window Title
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

/*************/
/* Titlebar */
/***********/

type TitleBar struct {
	MainTitle string
	Separator string

	appendTitle  string
	prependTitle string
	window       *gtk.Window
	titleNew     string
}

func TitleBarStructureNew(window *gtk.Window, title ...string) *TitleBar {
	tb := new(TitleBar)
	tb.Init(window, title...)
	return tb
}

// * Alternative use of TitleBarNew() * in the case where we need
// to import the structure rather than using it via the imported library.
func (tb *TitleBar) Init(window *gtk.Window, title ...string) {
	var err error
	tb.window = window
	tb.Separator = " - "

	if len(title) > 0 {
		tb.MainTitle = strings.Join(title, " ")
	} else {
		if tb.MainTitle, err = tb.window.GetTitle(); err != nil {
			log.Fatalf("Unable to acquire window title: %s\n", err.Error())
		}
	}
}

func (tb *TitleBar) Reset() string {
	tb.window.SetTitle(tb.MainTitle)
	return tb.MainTitle
}

func (tb *TitleBar) Update(toAdd []string, appendTo ...bool) string {
	var appendToTitle bool
	if len(appendTo) > 0 {
		appendToTitle = appendTo[0]
	}
	if len(toAdd) > 0 {
		if appendToTitle && len(toAdd) != 0 {
			tb.titleNew = tb.MainTitle + tb.Separator + strings.Join(toAdd, tb.Separator)

		} else if len(toAdd) != 0 {
			tb.titleNew = strings.Join(toAdd, tb.Separator) + tb.Separator + tb.MainTitle
		}
	} else {
		tb.titleNew = tb.MainTitle
	}
	tb.window.SetTitle(tb.titleNew)
	return tb.titleNew
}

/**************/
/* Statusbar */
/************/

type StatusBar struct {
	Messages  []string /* Each row contain associated strings refere to contextId number */
	Separator string
	statusbar *gtk.Statusbar
	contextId uint
	Prefix    []string
}

/* Init: Initialise structure to handle elements to be displayed. */
func StatusBarStructureNew(originStatusbar *gtk.Statusbar, prefix []string, stackId ...int) (bar *StatusBar) {
	bar = new(StatusBar)
	bar.Init(originStatusbar, prefix, stackId...)
	return
}

// Show: display / hide statusbar
func (bar *StatusBar) Show(state bool) {
	bar.statusbar.SetVisible(state)
}

/* Init: Initialise structure to handle elements to be displayed. */
func (bar *StatusBar) Init(originStatusbar *gtk.Statusbar, prefix []string, stackId ...int) {
	var stack int
	if len(stackId) == 0 {
		stack = 0
	} else {
		stack = stackId[0]
	}
	bar.Separator = " | "
	bar.statusbar = originStatusbar
	bar.contextId = bar.statusbar.GetContextId(strconv.Itoa(stack)) /* get contextId of stack */
	bar.Messages = make([]string, len(prefix))
	bar.Prefix = prefix
}

/* Add: add new element and return his own position. */
func (bar *StatusBar) Add(prefix, inString string) (position int) {
	bar.Prefix = append(bar.Prefix, prefix)
	bar.Messages = append(bar.Messages, prefix+" "+inString)
	bar.Disp()
	return len(bar.Messages) - 1
}

/* Add: set element at desired position. */
func (bar *StatusBar) Set(inString string, pos int) {
	if pos > len(bar.Messages)-1 || pos < 0 {
		inString = "Statusbar error: Invalid range to setting this message -> " + inString
		pos = len(bar.Messages) - 1
	}
	bar.Messages[pos] = inString
	bar.Disp()
}

/* Del: remove element at defined position and get the new length of elements. */
func (bar *StatusBar) Del(pos int) (newLength int) {
	copy(bar.Messages[pos:], bar.Messages[pos+1:])
	bar.Messages = bar.Messages[:cap(bar.Messages)-2]
	copy(bar.Prefix[pos:], bar.Prefix[pos+1:])
	bar.Prefix = bar.Prefix[:cap(bar.Prefix)-2]
	bar.Disp()
	return len(bar.Messages)
}

/* CleanAll: remove all elements (set to empty string) from the messages list. */
func (bar *StatusBar) CleanAll() {
	for idx, _ := range bar.Messages {
		bar.Messages[idx] = ""
	}
	bar.Disp()
}

/* Disp: display content of stored elements into statusbar */
func (bar *StatusBar) Disp() {
	var dispMessages []string
	for idxMessage, message := range bar.Messages {
		if len(message) != 0 {
			dispMessages = append(dispMessages, bar.Prefix[idxMessage]+" "+message)
		}
	}
	bar.statusbar.Push(bar.contextId, strings.Join(dispMessages, bar.Separator))
}
