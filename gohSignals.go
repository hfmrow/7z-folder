// gohSignals.go

/*
	Source file auto-generated on Mon, 26 Apr 2021 08:45:09 using Gotk3 Objects Handler v1.7.8
	©2018-21 hfmrow https://hfmrow.github.io

	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 hfmrow - 7z folder v1.6 github.com/hfmrow/7z-folder

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/***********************************************************/
/* Signals & Property Implementations part. You may add   */
/* properties or signals as you wish, best way is to add */
/* them by grouping objects names yours modifications   */
/* will be preserved on updates.                       */
/******************************************************/
// signalsPropHandler: initialise signals used by gtk objects ...
func signalsPropHandler() {
	obj.ButtonCancel.Connect("clicked", ButtonCancelClicked)
	obj.ButtonCompress.Connect("clicked", ButtonCompressClicked)
	obj.ButtonExit.Connect("clicked", ButtonExitClicked)
	obj.CheckbuttonAutoInc.Connect("toggled", RadioButtonGrpChanged)
	obj.EvenboxTop.Connect("button-release-event", EvenboxTopButtonReleaseEvent)
	obj.RadioButtonAppend.Connect("group-changed", RadioButtonGrpChanged)
	obj.RadioButtonNew.Connect("group-changed", RadioButtonGrpChanged)
	obj.RadioButtonUpdate.Connect("group-changed", RadioButtonGrpChanged)
}
