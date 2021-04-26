// pixbufByte.go

/*
	Source file auto-generated on Sun, 20 Oct 2019 13:50:31 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	©2019-21 H.F.M pixbuf byte library
	- provide simplified function to handle gotk3 pixbuff objects...

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/*************************************************/
/* Images functions, used to initialize objects */
/* You can use it to load your own embedded    */
/* images, icons ...                          */
/*********************************************/

// TryStockIcon: Convenience function to combine a stock-icons-theme
// and an asset/filename icon to send to "GetPixBuf()" function.
func TryStockIcon(iconLookupFlag gtk.IconLookupFlags, stockName, assetName interface{}) (names []interface{}) {
	names = append(names, stockName)
	names = append(names, iconLookupFlag)
	names = append(names, assetName)
	return
}

// OptPict: hold options for SetPict(), you can use OptPictStructNew()
// to set defaults values.
type OptPict struct {
	Size int
	// 0 = limit height, 1 = limit width, 2 = percents,
	// 3 = limit height or width depending on MaxSize value.
	ResizeMethode int
	// Rotate image / icon
	Rotate gdk.PixbufRotation
	// SpinButton / Entry ... icon positions
	Position gtk.EntryIconPosition
	// Opacity [0.01 to 1.0] work with: Box, ToggleButton, ToolButton, MenuButton, Button
	// a value of 0 means 1.0; 0.01 means 0
	Opacity float64
	// Ask for a GdkPixbufAnimation image
	Animated bool
	// Max size accepted relative to "ResizeMethode" = 3
	MaxSize int
}

// OptPictStructNew: 'values' options can be nil, ther is no obligation
// to use this function to create the structure because functions are
// designed to handle defaults options automatically, this is a convenient use.
// The options (each values can be nil) order (if used) must be:
// Size, ResizeMethode, Rotate, Position, Opacity, Animated, MaxSize.
func OptPictStructNew(values ...interface{}) OptPict {

	// Set default options
	opt := OptPict{
		Size:          0,
		ResizeMethode: 0,
		Rotate:        gdk.PIXBUF_ROTATE_NONE,
		Position:      gtk.ENTRY_ICON_PRIMARY,
		Opacity:       1.0,
		Animated:      false,
		MaxSize:       256,
	}

	var fillOptInt = func(i, index int) {
		switch index {

		case 0:
			opt.Size = i
		case 1:
			opt.ResizeMethode = i
		case 6:
			opt.MaxSize = i
		}
	}
	// Fill opt with given values
	if values != nil {
		for idx, val := range values {
			switch val.(type) {

			case nil:
				continue
			case gdk.PixbufRotation:
				opt.Rotate = val.(gdk.PixbufRotation)
			case gtk.EntryIconPosition:
				opt.Position = val.(gtk.EntryIconPosition)
			case float64:
				opt.Opacity = val.(float64)
			case bool:
				opt.Animated = val.(bool)
			case int:
				fillOptInt(val.(int), idx)
			}
		}
	}
	return opt
}

// getOptPict: Handling OptPict for many usages
func _getOptPict(options ...interface{}) OptPict {

	opt := OptPictStructNew()

	if len(options) > 0 {
		switch options[0].(type) {
		case int:
			opt.Size = options[0].(int)
		case interface{}:
			opt = options[0].(OptPict)
		}
	}
	return opt
}

// SetPic: Assign image to an Object depending on type, accept stock-icons-theme,
// filename or []byte. Use OptPict{} struct to set options, size can be applied
// directly as int or be nil.
// The options are described below:
// - Load a gif animated image and specify using animation (true), resizing not
//   allowed with animations.
//     SetPict(GtkImage, "linearProgressHorzBlue.gif", OptPict{Animated: true})
// - Resize to 32 pixels height, keep porportions & assign image to GtkButton.
//     SetPict(GtkButton, "stillImage.png", 32) or
//     SetPict(GtkButton, "stillImage.png", OptPict{Size: 32})
// - With default size, resizing not allowed for GtkSpinButton position at right-end.
//     SetPict(GtkSpinButton, "stillImage.png", OptPict{Position: gtk.ENTRY_ICON_SECONDARY})
// - Load a stock-icons-theme or an asset/filename icon to a GtkMenuButton where stock not available.
//     SetPict(gtk.MenuButton, TryStockIcon("open-menu-symbolic", "stillImage.png"), 18) or
//     SetPict(gtk.MenuButton, TryStockIcon("open-menu-symbolic", "stillImage.png"), OptPict{Size: 18})
// - Rotate upsidedown, resize 14px:
//     SetPict(gtk.MenuButton, TryStockIcon("open-menu-symbolic", "stillImage.png"), OptPict{Size: 14, Rotate: gdk.PIXBUF_ROTATE_UPSIDEDOWN})
func SetPict(iObject, varPath interface{}, options ...interface{}) (err error) {
	var inPixbuf *gdk.Pixbuf
	var inPixbufAnimation *gdk.PixbufAnimation
	var image *gtk.Image

	opt := _getOptPict(options...)

	// Since the default initialization value is set to 0,
	// we need to change it to show a visible image ...
	if opt.Opacity == 0 {
		opt.Opacity = 1
	}

	// Get pixbuff type (normal or animation for GtkImage)
	if opt.Animated {
		if inPixbufAnimation, err = GetPixBufAnimation(varPath); err == nil {
			image, err = gtk.ImageNewFromAnimation(inPixbufAnimation)
		}
	} else {
		if inPixbuf, err = GetPixBuf(varPath, opt); err == nil {
			image, err = gtk.ImageNewFromPixbuf(inPixbuf)
		}
	}

	// Handling the error if there is
	if err != nil {
		var filename string
		// Look for a stock icon or a filename or an ambedded asset
		switch path := varPath.(type) {
		case []interface{}:
			filename = path[2].(string)
		default:
			filename = path.(string)
		}
		if _, err = os.Stat(filename); len(filename) > 0 /*os.IsNotExist(err)*/ {
			err = fmt.Errorf("SetPict: [%v] %v", varPath, err)
			fmt.Println(err.Error()) // double output
		}
		return
	}

	// Objects parsing
	if image != nil {

		image.SetOpacity(opt.Opacity)

		switch object := iObject.(type) {

		case *gtk.Button: // Set Image to GtkButton
			object.SetImage(image)
			object.SetAlwaysShowImage(true)

		case *gtk.MenuButton: // Set Image to GtkMenuButton
			object.SetImage(image)
			object.SetAlwaysShowImage(true)

		case *gtk.ToggleButton: // Set Image to GtkToggleButton
			object.SetImage(image)
			object.SetAlwaysShowImage(true)

		case *gtk.ToolButton: // Set Image to GtkToolButton
			object.SetIconWidget(image)

		case *gtk.Box: // Add Image to GtkBox
			object.Add(image)

		case *gtk.SpinButton: // Set Icon to GtkSpinButton, No resize, No Animate.
			object.SetIconFromPixbuf(opt.Position, inPixbuf)

		case *gtk.ApplicationWindow: // Set Icon to GtkApplicationWindow, No Animate.
			object.SetIcon(inPixbuf)

		case *gtk.Window: // Set Icon to GtkWindow, No Animate.
			object.SetIcon(inPixbuf)

		case *gtk.Image: // Set Image to GtkImage
			if opt.Animated {
				object.SetFromAnimation(inPixbufAnimation)
			} else {
				object.SetFromPixbuf(inPixbuf)
			}
		}
	} else {
		err = fmt.Errorf("Image error: %v", err)
	}
	return
}

// GetPixBuf: Get gdk.PixBuf from stock, filename or []byte, depending
// on type. For an explanation of the options, see above, OptPict{}.
// note: If stock-icons-theme does not exist, the asset/filename icon is used.
func GetPixBuf(varPath interface{}, options ...interface{}) (outPixbuf *gdk.Pixbuf, err error) {
	var pixbufLoader *gdk.PixbufLoader

	opt := _getOptPict(options...)

	// Look for a stock icon and a filename or ambedded asset
	switch varPath.(type) {
	case []interface{}:
		var iconTheme *gtk.IconTheme
		if iconTheme, err = gtk.IconThemeGetDefault(); err == nil {
			stock := varPath.([]interface{})[0].(string)
			iconLookupFlags := varPath.([]interface{})[1].(gtk.IconLookupFlags)
			outPixbuf, err = iconTheme.LoadIcon(stock, opt.Size, iconLookupFlags)
		}
		if err != nil {
			varPath = varPath.([]interface{})[2]
			fmt.Println(err.Error() + ". try to load the secondary choice.")
		}
	}

	if outPixbuf == nil { // in this case, no stock icon was found or was not requested
		switch varPath.(type) {
		case string: // Its a filename
			outPixbuf, err = gdk.PixbufNewFromFile(varPath.(string))
		case []uint8: // Its a binary data (embedded asset)
			if pixbufLoader, err = gdk.PixbufLoaderNew(); err == nil {
				outPixbuf, err = pixbufLoader.WriteAndReturnPixbuf(varPath.([]byte))
			}
		}

		if err == nil && opt.Size != 0 {
			newWidth, newHeight := NormalizeSize(outPixbuf.GetWidth(), outPixbuf.GetHeight(), opt)
			if outPixbuf, err = outPixbuf.ScaleSimple(newWidth, newHeight, gdk.INTERP_HYPER); err == nil {
				outPixbuf, err = outPixbuf.RotateSimple(opt.Rotate)
			}
		}
	}
	return
}

// GetAnimationImage: Retrieve *gtk.Image animation 'gif' from []byte or filename.
func GetAnimationImage(varPath interface{}) (gifImage *gtk.Image, err error) {
	if anim, err := GetPixBufAnimation(varPath); err == nil {
		if gifImage, err = gtk.ImageNewFromAnimation(anim); err == nil {
			return gifImage, nil
		}
	}
	return nil, err
}

// GetPixBufAnimation: Get gdk.PixBufAnimation from filename or []byte, depending on type
func GetPixBufAnimation(varPath interface{}) (outPixbufAnimation *gdk.PixbufAnimation, err error) {
	var pixbufLoader *gdk.PixbufLoader
	switch varPath.(type) {
	case string:
		outPixbufAnimation, err = gdk.PixbufAnimationNewFromFile(varPath.(string))
	case []uint8:
		if pixbufLoader, err = gdk.PixbufLoaderNew(); err == nil {
			outPixbufAnimation, err = pixbufLoader.WriteAndReturnPixbufAnimation(varPath.([]byte))
		}
	}
	return
}

// NormalizeSize: compute new size with kept proportions based on defined format.
// formats: 0 = limit height, 1 = limit width, 2 = percents,
// 3 = limit height or width depending on MaxSize value.
func NormalizeSize(oldWidth, oldHeight int, options ...interface{}) (outWidth, outHeight int) {

	opt := _getOptPict(options...)

	switch opt.ResizeMethode {
	case 0: // limit Height
		outWidth = int(float64(oldWidth) * (float64(opt.Size) / float64(oldHeight)))
		outHeight = opt.Size
	case 1: // limit Width
		outWidth = opt.Size
		outHeight = int(float64(oldHeight) * (float64(opt.Size) / float64(oldWidth)))
	case 2: // percent
		outWidth = int((float64(oldWidth) * float64(opt.Size)) / 100)
		outHeight = int((float64(oldHeight) * float64(opt.Size)) / 100)
	case 3: // limit Height or Width
		switch {
		case oldWidth >= opt.MaxSize:
			opt.Size = opt.MaxSize
			opt.ResizeMethode = 1
			return NormalizeSize(oldWidth, oldHeight, opt)
		case oldHeight >= opt.MaxSize:
			opt.Size = opt.MaxSize
			opt.ResizeMethode = 0
			return NormalizeSize(oldWidth, oldHeight, opt)
		}
		opt.ResizeMethode = 0
		return NormalizeSize(oldWidth, oldHeight, opt)
	}
	return
}

/***************************************/
/* Embedded data conversion functions */
/* Used to make variable content     */
/* available in go-source           */
/***********************************/
// GetBytesFromVarAsset: Get []byte representation from file or asset, depending on type
func GetBytesFromVarAsset(varPath interface{}) (outBytes []byte, err error) {
	switch reflect.TypeOf(varPath).String() {
	case "string":
		return ioutil.ReadFile(varPath.(string))
	case "[]uint8":
		return varPath.([]byte), err
	}
	return
}

// HexToBytes: Convert Gzip Hex to []byte used for embedded binary in source code
func HexToBytes(varPath string, gzipData []byte) (outByte []byte) {
	r, err := gzip.NewReader(bytes.NewBuffer(gzipData))
	if err == nil {
		var bBuffer bytes.Buffer
		if _, err = io.Copy(&bBuffer, r); err == nil {
			if err = r.Close(); err == nil {
				return bBuffer.Bytes()
			}
		}
	}
	if err != nil {
		fmt.Printf("An error occurred while reading: %s\n%v\n", varPath, err.Error())
	}
	return
}
