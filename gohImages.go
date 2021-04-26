// gohImages.go

/*
	Source file auto-generated on Fri, 16 Apr 2021 22:07:08 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019-21 hfmrow - 7z folder v1.6 github.com/hfmrow/7z-folder
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

/**********************************************************/
/* This section preserve user modifications on update.   */
/* Images declarations, used to initialize objects with */
/* The SetPict() func, accept both kind of variables:  */
/* filename or []byte content in case of using        */
/* embedded binary data. The variables names are the */
/* same. "assetsDeclarationsUseEmbedded(bool)" func */
/* could be used to toggle between filenames and   */
/* embedded binary type. See SetPict()            */
/* declaration to learn more on how to use it.   */
/************************************************/
func assignImages() {
	SetPict(obj.ButtonCancel, signCancel20)
	SetPict(obj.ButtonCompress, zipYellow20)
	SetPict(obj.ButtonExit, logOut20)
	SetPict(obj.ImageMainTop, icon7zFolder500x48)
	SetPict(obj.MainWindow, icon7zip48)
}

/**********************************************************/
/* This section is rewritten on assets update.           */
/* Assets var declarations, this step permit to make a  */
/* bridge between the differents types used, string or */
/* []byte, and to simply switch from one to another.  */
/*****************************************************/
var mainGlade interface{}          // assets/glade/main.glade
var icon7zFolder500x48 interface{} // assets/images/icon-7zFolder-500x48.png
var icon7zip48 interface{}         // assets/images/icon-7zip-48.png
var logOut20 interface{}           // assets/images/Log-Out-20.png
var signCancel20 interface{}       // assets/images/Sign-cancel-20.png
var signSelect20 interface{}       // assets/images/Sign-Select-20.png
var zipYellow20 interface{}        // assets/images/zip-yellow-20.png
