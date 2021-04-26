// popupMenu.go

/*
	©2019-21 H.F.M. MIT license

	A simple builder for popup menu that can handle icons or not.
*/

package gtk3Import

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	gipf "github.com/hfmrow/gtk3Import/pixbuff"
)

var ( // Lib mapping
	HexToBytes = gipf.HexToBytes
	GetPixBuf  = gipf.GetPixBuf
)

type ItemOptionsType int

const (
	// Like GtkMenuItem
	ITEM_OPT_NORMAL ItemOptionsType = 1 << 0
	// Like GtkCheckMenuItem
	ITEM_OPT_CHECK ItemOptionsType = 1 << 1
	// Like GtkRadioMenuItem
	ITEM_OPT_RADIO ItemOptionsType = 1 << 2
	// Like GtkSeparatorMenuItem
	ITEM_OPT_SEPARATOR ItemOptionsType = 1 << 3
	// Add icon before GtkMenuItem
	ITEM_OPT_ICON ItemOptionsType = 1 << 4
	// Disable the use of mnemonic
	ITEM_OPT_NO_MNEMONIC ItemOptionsType = 1 << 5
	// Start a new group for GtkRadioButton(s)
	ITEM_OPT_RADIO_NEW_GRP ItemOptionsType = 1 << 6
	// align label to other (adding blank icon)
	ITEM_OPT_ALGN_LBL ItemOptionsType = 1 << 7
	// Box container alignement centered
	ITEM_OPT_ALGN_BOX_CNTR ItemOptionsType = 1 << 8
)

// PopupMenuStruct: means that only standard widgets will be used
// If you want to add icons, prefere using next structure 'PopupMenuIconStruct'
type PopupMenuStruct struct {
	Menu *gtk.Menu
	// Indicate what type of function the window has.
	WindowTypeHint gdk.WindowTypeHint // default (WINDOW_TYPE_HINT_POPUP_MENU)
	// left mouse button instead of right
	LMB,
	// Reserving space to toggle (gtkCheckMenuItem/GtkRadioMenuItem)
	ReserveToggleSize bool
	// space separator for box (used when image is present)
	BoxSpacing,
	// Contain the index of last GtkMenuItem created
	nextMenuItemIdx,
	// When a callback is defined (treeview RMB popup for example),
	// set the minimum selected row(s) to allow showing the popup.
	TreeViewMinSelectedRows int
	// Used to link radio button within a group
	lastRadioMenuItmGroup *gtk.RadioMenuItem
	// Where menuitems are stored
	items []combObj
	// Options that can be used with designed function
	OPT_NORMAL,
	OPT_CHECK,
	OPT_RADIO,
	OPT_SEPARATOR,
	OPT_NO_MNEMONIC,
	OPT_RADIO_NEW_GRP ItemOptionsType
}

type combObj struct {
	widget      interface{}
	tType       ItemOptionsType
	combined    bool
	xMenuItem   *gtk.MenuItem
	hasCallback bool
}

func combObjNew() combObj {
	c := new(combObj)
	return *c
}

/*
 * Classic version
 */
// Retrieve current underlayed Gtk MenuItem
func (c *combObj) getMenuItem() *gtk.MenuItem {
	obj := c.widget
	menuItem, ok := obj.(*gtk.MenuItem)
	if !ok {
		log.Printf("[Warning!]: Unable to convert %T to (*gtk.MenuItem)", obj)
	}
	// }
	return menuItem
}

// Retrieve current underlayed GtkCheckMenuItem
func (c *combObj) getCheckMenuItem() *gtk.CheckMenuItem {
	obj := c.widget
	menuItem, ok := obj.(*gtk.CheckMenuItem)
	if !ok {
		log.Printf("[Warning!]: Unable to convert %T to (*gtk.CheckMenuItem)", obj)
	}
	return menuItem
}

// Retrieve current underlayed GtkRadioMenuItem
func (c *combObj) getRadioMenuItem() *gtk.RadioMenuItem {
	obj := c.widget
	menuItem, ok := obj.(*gtk.RadioMenuItem)
	if !ok {
		log.Printf("[Warning!]: Unable to convert %T to (*gtk.RadioMenuItem)", obj)
	}
	return menuItem
}

// Retrieve current underlayed GtkSeparatorMenuItem
func (c *combObj) getSeparatorMenuItem() *gtk.SeparatorMenuItem {
	obj := c.widget
	menuItem, ok := obj.(*gtk.SeparatorMenuItem)
	if !ok {
		log.Printf("[Warning!]: Unable to convert %T to (*gtk.SeparatorMenuItem)", obj)
	}
	return menuItem
}

// Retrieve current underlayed GtkMenuItem combined version (Icon)
func (c *combObj) getCombMenuItem() *gtk.MenuItem {
	return c.xMenuItem
}

/*
 * Icon version
 */
// Retrieve current underlayed GtkLabel
func (c *combObj) getLabel() *gtk.Label {
	obj := c.widget
	menuItem, ok := obj.(*gtk.Label)
	if !ok {
		log.Printf("[Warning!]: Unable to convert %T to (*gtk.Label)", obj)
	}
	return menuItem
}

// Retrieve current underlayed GtkCheckButton
func (c *combObj) getCheckButton() *gtk.CheckButton {
	obj := c.widget
	menuItem, ok := obj.(*gtk.CheckButton)
	if !ok {
		log.Printf("[Warning!]: Unable to convert %T to (*gtk.CheckButton)", obj)
	}
	return menuItem
}

// Retrieve current underlayed GtkRadioButton
func (c *combObj) getRadioButton() *gtk.RadioButton {
	obj := c.widget
	menuItem, ok := obj.(*gtk.RadioButton)
	if !ok {
		log.Printf("[Warning!]: Unable to convert %T to (*gtk.RadioButton)", obj)
	}
	return menuItem
}

// PopupMenuNew: return a new PopupMenuStruct structure
func PopupMenuStructNew() (pop *PopupMenuStruct) {
	pop = new(PopupMenuStruct)

	pop.lastRadioMenuItmGroup = nil
	pop.ReserveToggleSize = false
	pop.TreeViewMinSelectedRows = 1
	pop.WindowTypeHint = gdk.WINDOW_TYPE_HINT_POPUP_MENU

	pop.OPT_NORMAL = ITEM_OPT_NORMAL
	pop.OPT_CHECK = ITEM_OPT_CHECK
	pop.OPT_RADIO = ITEM_OPT_RADIO
	pop.OPT_SEPARATOR = ITEM_OPT_SEPARATOR
	pop.OPT_NO_MNEMONIC = ITEM_OPT_NO_MNEMONIC
	pop.OPT_RADIO_NEW_GRP = ITEM_OPT_RADIO_NEW_GRP
	return
}

// Retrieve current GtkMenuItem
func (pop *PopupMenuStruct) GetMenuItemCurrent() *gtk.MenuItem {
	return pop.items[pop.nextMenuItemIdx-1].getMenuItem()
}

// Retrieve current GtkCheckMenuItem
func (pop *PopupMenuStruct) GetCheckMenuItemCurrent() *gtk.CheckMenuItem {
	return pop.items[pop.nextMenuItemIdx-1].getCheckMenuItem()
}

// Retrieve current GtkRadioMenuItem
func (pop *PopupMenuStruct) GetRadioMenuItemCurrent() *gtk.RadioMenuItem {
	return pop.items[pop.nextMenuItemIdx-1].getRadioMenuItem()
}

// Retrieve current GtkSeparatorMenuItem
func (pop *PopupMenuStruct) GetSeparatorMenuItemCurrent() *gtk.SeparatorMenuItem {
	return pop.items[pop.nextMenuItemIdx-1].getSeparatorMenuItem()
}

// Retrieve the IMenuItem at 'idx'
// NOTE: Require type assertion!!!
func (pop *PopupMenuStruct) GetIMenuItemIdx(idx int) gtk.IMenuItem {
	return pop.items[idx].widget.(gtk.IMenuItem)
}

func (pop *PopupMenuStruct) buildXMenuItem(lbl string, ico, activateFunction interface{}, opt ItemOptionsType) (*combObj, error) {
	var (
		cmbo              combObj
		err               error
		menuItem          *gtk.MenuItem
		checkMenuItem     *gtk.CheckMenuItem
		radioMenuItem     *gtk.RadioMenuItem
		separatorMenuItem *gtk.SeparatorMenuItem
	)
	switch {
	case opt&ITEM_OPT_SEPARATOR != 0:
		separatorMenuItem, err = gtk.SeparatorMenuItemNew()
		if err == nil {
			if len(lbl) > 0 {
				separatorMenuItem.SetLabel(lbl)
			}
			cmbo = combObj{widget: separatorMenuItem, tType: ITEM_OPT_SEPARATOR, combined: false}
			return &cmbo, nil
		}
	case opt&ITEM_OPT_CHECK != 0:
		if opt&ITEM_OPT_NO_MNEMONIC != 0 {
			checkMenuItem, err = gtk.CheckMenuItemNewWithLabel(lbl)
		} else {
			checkMenuItem, err = gtk.CheckMenuItemNewWithMnemonic(lbl)
		}
		if err == nil {
			switch f := activateFunction.(type) {
			case func(c *gtk.CheckMenuItem):
				checkMenuItem.Connect("activate", f)
			default:
				return nil, fmt.Errorf(
					"AddItem: callback for %T require 'func(c *gtk.CheckMenuItem)' type!",
					checkMenuItem)
			}
			cmbo = combObj{widget: checkMenuItem, tType: ITEM_OPT_CHECK, combined: false}
			return &cmbo, nil
		}
		return nil, err

	case opt&ITEM_OPT_RADIO != 0, opt&ITEM_OPT_RADIO_NEW_GRP != 0:
		if opt&ITEM_OPT_NO_MNEMONIC != 0 {
			radioMenuItem, err = gtk.RadioMenuItemNewWithLabel(nil, lbl)
		} else {
			radioMenuItem, err = gtk.RadioMenuItemNewWithMnemonic(nil, lbl)
		}
		if err == nil {
			switch f := activateFunction.(type) {
			case func(r *gtk.RadioMenuItem):
				radioMenuItem.Connect("activate", f)
			default:
				return nil, fmt.Errorf(
					"AddItem: callback for %T require 'func(r *gtk.RadioMenuItem)' type!",
					radioMenuItem)
			}

			cmbo = combObj{widget: radioMenuItem, tType: ITEM_OPT_RADIO, combined: false}
			return &cmbo, nil
		}
	case opt&ITEM_OPT_NORMAL != 0:
		if opt&ITEM_OPT_NO_MNEMONIC != 0 {
			menuItem, err = gtk.MenuItemNewWithLabel(lbl)
		} else {
			menuItem, err = gtk.MenuItemNewWithMnemonic(lbl)
		}
		if err == nil {
			if activateFunction != nil {
				switch f := activateFunction.(type) {
				case func():
					menuItem.Connect("activate", f)
				}
			}
			cmbo = combObj{widget: menuItem, tType: ITEM_OPT_NORMAL, combined: false}
			return &cmbo, nil
		}
	default:
		return nil, fmt.Errorf("AddItemWithOptions: Unable to define choosen option(s)")
	}

	return nil, err
}

// AddItem: Add item to menu. This version handle
// checkbutton/radiobutton with groups and classic labels
func (pop *PopupMenuStruct) AddItem(
	lbl string, activateFunction interface{}, options ...interface{}) (widgetIdx int, err error) {

	var (
		opt  ItemOptionsType = ITEM_OPT_NORMAL // Default Option
		icon interface{}
	)
	switch len(options) {
	case 1:
		opt = opt | options[0].(ItemOptionsType)
		icon = nil
	}

	cmbo, err := pop.buildXMenuItem(lbl, icon, activateFunction, opt)
	if err != nil {
		return -1, err
	}
	pop.items = append(pop.items, *cmbo)

	// Handling RadioMenuItem groups
	if opt&ITEM_OPT_RADIO != 0 || opt&ITEM_OPT_RADIO_NEW_GRP != 0 {
		if opt&ITEM_OPT_RADIO_NEW_GRP != 0 {
			pop.lastRadioMenuItmGroup = cmbo.getRadioMenuItem()
		}
		if pop.lastRadioMenuItmGroup == nil {
			pop.lastRadioMenuItmGroup = cmbo.getRadioMenuItem()
		}
		cmbo.getRadioMenuItem().JoinGroup(pop.lastRadioMenuItmGroup)
	}

	pop.nextMenuItemIdx++
	return pop.nextMenuItemIdx - 1, err
}

// MenuBuild: Build popupmenu.
func (pop *PopupMenuStruct) MenuBuild() *gtk.Menu {
	var err error
	if pop.Menu, err = gtk.MenuNew(); err == nil {
		pop.Menu = appendToExistingMenu(pop.items, pop.Menu)
	} else {
		log.Println("Popup menu creation error !")
		return nil
	}
	pop.Menu.SetProperty("reserve_toggle_size", pop.ReserveToggleSize)
	pop.Menu.SetProperty("menu_type_hint", pop.WindowTypeHint)
	return pop.Menu
}

// CheckRMB: Check whether an event comes from the right button of the
// mouse and display the popup if it is the case at the mouse position.
func (pop *PopupMenuStruct) CheckRMB(obj interface{}, event *gdk.Event) bool {
	return checkRMB(pop.LMB, pop.Menu, event)
}

// CheckRMBFromTreeView: May be directly used as callback function for
// TreeView' "button-press-event" signal, considere to initialize the
// popup menu before setting this function as callback. Otherwise, the
// call will generate error "nil pointer ..."
// NOTE: "button-press-event" signal instead of "button-release-event"
// must be used to avoid deselecting rows on right click.
func (pop *PopupMenuStruct) CheckRMBFromTreeView(tv *gtk.TreeView, event *gdk.Event) bool {
	return checkRMBFromTreeView(pop.LMB, pop.Menu, tv, event, pop.TreeViewMinSelectedRows)
}

// AppendToExistingMenu: append "MenuItems" to an existing "*gtk.Menu"
// Useful when you want to just add some entries to the context menu that
// already exist in a gtk.TextView or gtk.Entry by using "populate-popup"
// signal. Notice: GtkWidget > GtkMenu:
// menu := &gtk.Menu{gtk.MenuShell{gtk.Container{*w}}}
// Each connection to signal need to re-create the entire menu.
func (pop *PopupMenuStruct) AppendToExistingMenu(menu *gtk.Menu) *gtk.Menu {
	return appendToExistingMenu(pop.items, menu)
}

func appendToExistingMenu(menuItems []combObj, menu *gtk.Menu) *gtk.Menu {
	for _, menuItem := range menuItems {
		switch menuItem.tType {
		case ITEM_OPT_SEPARATOR:
			menu.Append(menuItem.getSeparatorMenuItem())
		default:
			if menuItem.combined {
				menu.Append(menuItem.xMenuItem)
			} else {
				switch menuItem.tType {
				case ITEM_OPT_NORMAL:
					menu.Append(menuItem.getMenuItem())
				case ITEM_OPT_CHECK:
					menu.Append(menuItem.getCheckMenuItem())
				case ITEM_OPT_RADIO:
					menu.Append(menuItem.getRadioMenuItem())
				}
			}
		}
	}
	menu.ShowAll()
	return menu
}

/*
 * Icon version, this structure hold a home-made version of GtkMenu with Icon
 * handling. Since Gtk3 does not include icon possibility anymore, i've made
 *  this implementation to continue using icons in my menu.
 * - GtkMenuItem is replaced by a GtkBox which embeds a GtkLabel, everything will be embedded in a GtkMenuItem.
 * - GtkCheckMenuItem is replaced by a GtkBox which embeds a GtkCheckButton, everything will be embedded in a GtkMenuItem.
 * - GtkRadioMenuItem is replaced by a GtkBox which embeds a GtkRadioButton, everything will be embedded in a GtkMenuItem.
 * - GtkSeparatorMenuItem still untouched.
 * Note:	for Icon, each GtkBox may embed/or not a GtkImage at first.
 * 			The "signal" emitted by the objects is processed transparently.
 *			All objects are accessible using they own method just after
*			creation or later, using they indexes and respective methods.
*/
// PopupMenuIconStruct: Structure that hold popup menu options and methods.
// A simple builder for popup menu that may handle icons.
// Instead of classics GtkMenuItem, this structure use:
// - GtkLabel as GtkMenuItem, it will be embedded in a GtkMenuItem
// - GtkCheckButton as GtkCheckMenuItem, it will be embedded in a GtkMenuItem
// - GtkRadioButton as GtkRadioMenuItem, it will be embedded in a GtkMenuItem
// To get the active widget, you can use 'GetXXXCurrent' or 'GetXXXIdx' method
// Note: the "toggled" signal is handled as transparent as possible.
type PopupMenuIconStruct struct {
	Menu *gtk.Menu
	// Indicate what type of function the window has.
	WindowTypeHint gdk.WindowTypeHint // default (WINDOW_TYPE_HINT_POPUP_MENU)
	// left mouse button instead of right
	LMB,
	// Reserving space to toggle (gtkCheckMenuItem/GtkRadioMenuItem)
	ReserveToggleSize bool
	// Define icons size
	IconsSize,
	// space separator for box (used for image)
	BoxSpacing,
	// Contain the last GtkMenuItme created
	nextMenuItemIdx,
	// Used to link radio button within a group
	lastRadioButtonGroupIdx,
	// When a callback is defined (treeview RMB popup for example),
	// set the minimum selected row(s) to allow showing the popup.
	TreeViewMinSelectedRows int
	// Used to link radio button within a group
	lastRadioButtonGroup *gtk.RadioButton
	// Where menuitems are stored
	items []combObj
	// Options that can be used with designed function
	OPT_NORMAL,
	OPT_CHECK,
	OPT_RADIO,
	OPT_SEPARATOR,
	OPT_ICON,
	OPT_NO_MNEMONIC,
	OPT_RADIO_NEW_GRP,
	OPT_ALGN_LBL,
	OPT_ALGN_BOX_CNTR ItemOptionsType
}

// PopupMenuNew: return a new PopupMenuIconStruct structure
func PopupMenuIconStructNew() (pop *PopupMenuIconStruct) {
	pop = new(PopupMenuIconStruct)

	pop.IconsSize = 18
	pop.BoxSpacing = 4

	pop.lastRadioButtonGroup = nil
	pop.ReserveToggleSize = false
	pop.TreeViewMinSelectedRows = 1
	pop.WindowTypeHint = gdk.WINDOW_TYPE_HINT_POPUP_MENU

	pop.OPT_NORMAL = ITEM_OPT_NORMAL
	pop.OPT_CHECK = ITEM_OPT_CHECK
	pop.OPT_RADIO = ITEM_OPT_RADIO
	pop.OPT_SEPARATOR = ITEM_OPT_SEPARATOR
	pop.OPT_ICON = ITEM_OPT_ICON
	pop.OPT_NO_MNEMONIC = ITEM_OPT_NO_MNEMONIC
	pop.OPT_RADIO_NEW_GRP = ITEM_OPT_RADIO_NEW_GRP
	pop.OPT_ALGN_LBL = ITEM_OPT_ALGN_LBL
	pop.OPT_ALGN_BOX_CNTR = ITEM_OPT_ALGN_BOX_CNTR
	return
}

// Retrieve current GtkMenuItem
func (pop *PopupMenuIconStruct) GetMenuItemCurrent() *gtk.MenuItem {
	return pop.items[pop.nextMenuItemIdx-1].getCombMenuItem()
}

// Retrieve current undferlayed GtkLabel
func (pop *PopupMenuIconStruct) GetLabelCurrent() *gtk.Label {
	return pop.items[pop.nextMenuItemIdx-1].getLabel()
}

// Retrieve current undferlayed GtkCheckButton
func (pop *PopupMenuIconStruct) GetCheckButtonCurrent() *gtk.CheckButton {
	return pop.items[pop.nextMenuItemIdx-1].getCheckButton()
}

// Retrieve current undferlayed GtkRadioButton
func (pop *PopupMenuIconStruct) GetRadioButtonCurrent() *gtk.RadioButton {
	return pop.items[pop.nextMenuItemIdx-1].getRadioButton()
}

// Retrieve current GtkSeparatorMenuItem
func (pop *PopupMenuIconStruct) GetSeparatorMenuItemCurrent() *gtk.SeparatorMenuItem {
	return pop.items[pop.nextMenuItemIdx-1].getSeparatorMenuItem()
}

// Get the MenuItem at 'idx', which contains underlying widgets
func (pop *PopupMenuIconStruct) GetMenuItemIdx(idx int) *gtk.MenuItem {
	return pop.items[idx].xMenuItem
}

// Retrieve the underlayed IWidget at 'idx'
func (pop *PopupMenuIconStruct) GetIWidgetIdx(idx int) gtk.IWidget {
	return pop.items[idx].widget.(gtk.IWidget)
}

// AddItemWithOptions: Add items to menu. This version handle
// checkbutton/radiobutton with groups and classic labels with
// or without icon and more options ...
func (pop *PopupMenuIconStruct) buildXMenuItem(lbl string, ico, activateFunction interface{}, opt ItemOptionsType) (*combObj, error) {
	var (
		cmbo        combObj
		err         error
		menuItem    *gtk.MenuItem
		pixbuf      *gdk.Pixbuf
		image       *gtk.Image
		label       *gtk.Label
		checkbutton *gtk.CheckButton
		radiobutton *gtk.RadioButton
		box         *gtk.Box

		addImage = func(ico interface{}) error {
			// The function below is a part of personal gotk3 library that
			// allow to load image with some facilities. May handle a full
			// filename or an embedded binary data (hex/zip compressed).
			if pixbuf, err = GetPixBuf(ico, pop.IconsSize); err == nil {
				if image, err = gtk.ImageNewFromPixbuf(pixbuf); err == nil {
					box.PackStart(image, false, false, 0)
					return nil
				}
			}
			return err
		}
	)
	if box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, pop.BoxSpacing); err == nil {
		if opt&ITEM_OPT_ALGN_BOX_CNTR == 0 {
			box.SetHAlign(gtk.ALIGN_START)
		}
		if menuItem, err = gtk.MenuItemNew(); err == nil {
			menuItem.Container.Add(box)

			if opt&ITEM_OPT_ICON != 0 {
				if ico == nil {
					return nil, fmt.Errorf("AddItemWithOptions: Icon requested but not provided")
				}
				err = addImage(ico)
			} else if opt&ITEM_OPT_ALGN_LBL != 0 {
				err = addImage(blankIcon)
				image.SetOpacity(0)
			}
			if err != nil {
				return nil, err
			}

			switch {
			case opt&ITEM_OPT_SEPARATOR != 0:
				separatorItem, err := gtk.SeparatorMenuItemNew()
				if err != nil {
					return nil, err
				}
				// If empty, do not cast SetLabel, permit to draw a 1 pixel line
				// instead of a plain (full height) mark as separator.
				if len(lbl) > 0 {
					separatorItem.SetLabel(lbl)
				}
				cmbo = combObj{widget: separatorItem, tType: ITEM_OPT_SEPARATOR, combined: false}

			case opt&ITEM_OPT_CHECK != 0:
				if opt&ITEM_OPT_NO_MNEMONIC != 0 {
					checkbutton, err = gtk.CheckButtonNewWithLabel(lbl)
				} else {
					checkbutton, err = gtk.CheckButtonNewWithMnemonic(lbl)
				}
				if err != nil {
					return nil, err
				}
				box.PackEnd(checkbutton, true, true, 0)
				switch f := activateFunction.(type) {
				case func(chk *gtk.CheckButton):
					menuItem.Connect("activate", func() {
						checkbutton.SetActive(!checkbutton.GetActive())
						f(checkbutton)
					})
				default:
					return nil, fmt.Errorf("AddItemWithOptions: callback type 'func(chk *gtk.CheckButton)' required!")
				}
				cmbo = combObj{widget: checkbutton, tType: ITEM_OPT_CHECK, combined: true, xMenuItem: menuItem}

			case opt&ITEM_OPT_RADIO != 0, opt&ITEM_OPT_RADIO_NEW_GRP != 0:
				if opt&ITEM_OPT_NO_MNEMONIC != 0 {
					radiobutton, err = gtk.RadioButtonNewWithLabel(nil, lbl)
				} else {
					radiobutton, err = gtk.RadioButtonNewWithMnemonic(nil, lbl)
				}
				if err != nil {
					return nil, err
				}
				box.PackEnd(radiobutton, true, true, 0)
				switch f := activateFunction.(type) {
				case func(chk *gtk.RadioButton):
					menuItem.Connect("activate", func() {
						radiobutton.SetActive(!radiobutton.GetActive())
						f(radiobutton)
					})
				default:
					return nil, fmt.Errorf("AddItemWithOptions: callback type 'func(rbn *gtk.RadioButton)' required!")
				}
				cmbo = combObj{widget: radiobutton, tType: ITEM_OPT_RADIO, combined: true, xMenuItem: menuItem}

			case opt&ITEM_OPT_NORMAL != 0:
				if opt&ITEM_OPT_NO_MNEMONIC != 0 {
					label, err = gtk.LabelNew(lbl)
				} else {
					label, err = gtk.LabelNewWithMnemonic(lbl)
				}
				if err != nil {
					return nil, err
				}
				box.PackEnd(label, true, true, 0)
				switch f := activateFunction.(type) {
				case func():
					menuItem.Connect("activate", f)
				}
				cmbo = combObj{widget: label, tType: ITEM_OPT_NORMAL, combined: true, xMenuItem: menuItem}

			default:
				return nil, fmt.Errorf("AddItemWithOptions: Unable to define choosen option(s)")
			}
			menuItem.ShowAll()
			return &cmbo, nil
		}
	}

	return nil, err
}

// AddItem: Add item to menu. This version handle
// checkbutton/radiobutton with groups and classic labels with
// or without icon and more options ...
func (pop *PopupMenuIconStruct) AddItem(
	lbl string, activateFunction interface{}, options ...interface{}) (widgetIdx int, err error) {

	var (
		opt  ItemOptionsType = ITEM_OPT_NORMAL // Default Option
		icon interface{}     = nil
	)
	switch len(options) {
	case 1:
		opt = opt | options[0].(ItemOptionsType)
	case 2:
		opt = opt | options[0].(ItemOptionsType)
		icon = options[1]
	}

	cmbo, err := pop.buildXMenuItem(lbl, icon, activateFunction, opt)
	if err != nil {
		return -1, err
	}
	pop.items = append(pop.items, *cmbo)

	// Handling RadioMenuItem groups
	if opt&ITEM_OPT_RADIO != 0 || opt&ITEM_OPT_RADIO_NEW_GRP != 0 {
		if opt&ITEM_OPT_RADIO_NEW_GRP != 0 {
			pop.lastRadioButtonGroup = cmbo.getRadioButton()
		}
		if pop.lastRadioButtonGroup == nil {
			pop.lastRadioButtonGroup = cmbo.getRadioButton()
		}
		cmbo.getRadioButton().JoinGroup(pop.lastRadioButtonGroup)
	}

	pop.nextMenuItemIdx++
	return pop.nextMenuItemIdx - 1, err
}

// MenuBuild: Build popupmenu.
func (pop *PopupMenuIconStruct) MenuBuild() *gtk.Menu {
	var err error
	if pop.Menu, err = gtk.MenuNew(); err == nil {
		pop.Menu = appendToExistingMenu(pop.items, pop.Menu)
	} else {
		log.Println("Popup menu creation error !")
		return nil
	}
	pop.Menu.SetProperty("reserve_toggle_size", pop.ReserveToggleSize)
	pop.Menu.SetProperty("menu_type_hint", pop.WindowTypeHint)
	return pop.Menu
}

// CheckRMB: Check whether an event comes from the right button of the
// mouse and display the popup if it is the case at the mouse position.
func (pop *PopupMenuIconStruct) CheckRMB(obj interface{}, event *gdk.Event) bool {
	return checkRMB(pop.LMB, pop.Menu, event)
}

// CheckRMBFromTreeView: May be directly used as callback function for
// TreeView' "button-press-event" signal, considere to initialize the
// popup menu before setting this function as callback. Otherwise, the
// call will generate error "nil pointer ..."
func (pop *PopupMenuIconStruct) CheckRMBFromTreeView(tv *gtk.TreeView, event *gdk.Event) bool {
	return checkRMBFromTreeView(pop.LMB, pop.Menu, tv, event, pop.TreeViewMinSelectedRows)
}

func (pop *PopupMenuIconStruct) AppendToExistingMenu(menu *gtk.Menu) *gtk.Menu {
	return appendToExistingMenu(pop.items, menu)
}

func checkRMB(popLMB bool, menu *gtk.Menu, event *gdk.Event) bool {
	eventButton := gdk.EventButtonNewFromEvent(event)
	if uint(eventButton.Button()) == mouseBtn(popLMB) {
		menu.PopupAtPointer(event)
		return true /* we handled this */
	}
	return false /* we did not handle this */
}

func checkRMBFromTreeView(popLMB bool, menu *gtk.Menu, tv *gtk.TreeView, event *gdk.Event, minSelRows int) bool {
	if selection, err := tv.GetSelection(); err == nil {

		eventButton := gdk.EventButtonNewFromEvent(event)
		if uint(eventButton.Button()) == mouseBtn(popLMB) {
			// If right click is not over a selected row then
			// unselect all and select the row under the cursor
			eventMotion := gdk.EventMotionNewFromEvent(event)
			x, y := eventMotion.MotionVal()
			if path, _, _, _, ok := tv.GetPathAtPos(int(x), int(y)); ok {
				if !selection.PathIsSelected(path) {
					selection.UnselectAll()
					selection.SelectPath(path)
				}
			}
			if selection.CountSelectedRows() >= minSelRows {
				// Display popup menu
				menu.PopupAtPointer(event)
				return true /* we handled this */
			}
		}
	}
	return false /* we did not handle this */
}

// mouseBtn: get uint value of specified button to match
func mouseBtn(popLMB bool) uint {
	if popLMB {
		return 1 // LMB
	}
	return 3 // RMB
}

// Icon used as transparent to fill space before underlayed widget when
// ITEM_OPT_ALGN_LBL is used.
var blankIcon = HexToBytes("blank", []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x00\x3a\x02\xc5\xfd\x89\x50\x4e\x47\x0d\x0a\x1a\x0a\x00\x00\x00\x0d\x49\x48\x44\x52\x00\x00\x00\x10\x00\x00\x00\x10\x08\x06\x00\x00\x00\x1f\xf3\xff\x61\x00\x00\x01\x84\x69\x43\x43\x50\x49\x43\x43\x20\x70\x72\x6f\x66\x69\x6c\x65\x00\x00\x28\x91\x7d\x91\x3d\x48\xc3\x40\x1c\xc5\x5f\x53\xa5\x22\x15\x41\x0b\x8a\x38\x64\xa8\x4e\x16\x44\x45\x74\xd3\x2a\x14\xa1\x42\xa8\x15\x5a\x75\x30\xb9\xf4\x43\x68\xd2\x90\xa4\xb8\x38\x0a\xae\x05\x07\x3f\x16\xab\x0e\x2e\xce\xba\x3a\xb8\x0a\x82\xe0\x07\x88\x9b\x9b\x93\xa2\x8b\x94\xf8\xbf\xa4\xd0\x22\xc6\x83\xe3\x7e\xbc\xbb\xf7\xb8\x7b\x07\x08\xb5\x12\xd3\xac\xb6\x51\x40\xd3\x6d\x33\x95\x88\x8b\x99\xec\x8a\x18\x7a\x85\x80\x1e\xf4\x61\x1a\x82\xcc\x2c\x63\x56\x92\x92\xf0\x1d\x5f\xf7\x08\xf0\xf5\x2e\xc6\xb3\xfc\xcf\xfd\x39\xba\xd4\x9c\xc5\x80\x80\x48\x3c\xc3\x0c\xd3\x26\x5e\x27\x9e\xdc\xb4\x0d\xce\xfb\xc4\x11\x56\x94\x55\xe2\x73\xe2\x11\x93\x2e\x48\xfc\xc8\x75\xc5\xe3\x37\xce\x05\x97\x05\x9e\x19\x31\xd3\xa9\x39\xe2\x08\xb1\x58\x68\x61\xa5\x85\x59\xd1\xd4\x88\x27\x88\xa3\xaa\xa6\x53\xbe\x90\xf1\x58\xe5\xbc\xc5\x59\x2b\x55\x58\xe3\x9e\xfc\x85\xe1\x9c\xbe\xbc\xc4\x75\x9a\x83\x48\x60\x01\x8b\x90\x20\x42\x41\x05\x1b\x28\xc1\x46\x8c\x56\x9d\x14\x0b\x29\xda\x8f\xfb\xf8\x07\x5c\xbf\x44\x2e\x85\x5c\x1b\x60\xe4\x98\x47\x19\x1a\x64\xd7\x0f\xfe\x07\xbf\xbb\xb5\xf2\xe3\x63\x5e\x52\x38\x0e\xb4\xbf\x38\xce\xc7\x10\x10\xda\x05\xea\x55\xc7\xf9\x3e\x76\x9c\xfa\x09\x10\x7c\x06\xae\xf4\xa6\xbf\x5c\x03\xa6\x3e\x49\xaf\x36\xb5\xe8\x11\xd0\xbd\x0d\x5c\x5c\x37\x35\x65\x0f\xb8\xdc\x01\xfa\x9f\x0c\xd9\x94\x5d\x29\x48\x53\xc8\xe7\x81\xf7\x33\xfa\xa6\x2c\xd0\x7b\x0b\x74\xae\x7a\xbd\x35\xf6\x71\xfa\x00\xa4\xa9\xab\xe4\x0d\x70\x70\x08\x0c\x17\x28\x7b\xcd\xe7\xdd\x1d\xad\xbd\xfd\x7b\xa6\xd1\xdf\x0f\x6b\xbe\x72\xa4\xea\xad\xb2\x9e\x00\x00\x00\x06\x62\x4b\x47\x44\x00\xff\x00\xff\x00\xff\xa0\xbd\xa7\x93\x00\x00\x00\x09\x70\x48\x59\x73\x00\x00\x2e\x23\x00\x00\x2e\x23\x01\x78\xa5\x3f\x76\x00\x00\x00\x07\x74\x49\x4d\x45\x07\xe5\x02\x12\x17\x3b\x0b\xb8\x84\x6f\xa0\x00\x00\x00\x19\x74\x45\x58\x74\x43\x6f\x6d\x6d\x65\x6e\x74\x00\x43\x72\x65\x61\x74\x65\x64\x20\x77\x69\x74\x68\x20\x47\x49\x4d\x50\x57\x81\x0e\x17\x00\x00\x00\x12\x49\x44\x41\x54\x38\xcb\x63\x60\x18\x05\xa3\x60\x14\x8c\x02\x08\x00\x00\x04\x10\x00\x01\x85\x3f\xaa\x72\x00\x00\x00\x00\x49\x45\x4e\x44\xae\x42\x60\x82\x01\x00\x00\xff\xff\x1a\x3e\x56\xc7\x3a\x02\x00\x00"))

/******************** OUTDATED START **********************
 * Methods below MUST NOT be used anymore to display icons
 * Otherwise, without icon they still available
 */
// AddItem: Add items to menu.
// func (pop *PopupMenuIconStruct) AddItem(lbl string, activateFunction interface{},
// 	icon ...interface{}) (err error) {

// 	fmt.Printf("Warning!: POPUP/MENU implementation must be changed using new version\n")

// 	var menuItem *gtk.MenuItem
// 	var pixbuf *gdk.Pixbuf

// 	if len(icon) != 0 {
// 		// The function below is a part of personal gotk3 library
// 		// allow to load image with some facilities. May handle
// 		// filename or embedded binary data (hex/zip compressed).
// 		// pixbuf, err = gdk.PixbufNewFromFile(filename)
// 		pixbuf, err = GetPixBuf(icon[0], pop.IconsSize)
// 	}

// 	if pop.WithIcons {
// 		menuItem, err = pop.menuItemNewWithImage(lbl, pixbuf)
// 	} else {
// 		menuItem, err = gtk.MenuItemNewWithMnemonic(lbl)
// 	}

// 	// Handle the "activate" signal from the related item.
// 	if err == nil {
// 		menuItem.SetUseUnderline(true)
// 		if activateFunction != nil {
// 			menuItem.Connect("activate", activateFunction.(func()))
// 		}
// 		pop.Items = append(pop.Items, menuItem)
// 		// pop.separators = append(pop.separators, nil)
// 		pop.nextMenuItemIdx++
// 	}
// 	return err
// }

// // AddCheckMenuItem: Add items to menu.
// func (pop *PopupMenuIconStruct) AddCheckMenuItem(lbl string, activateFunction interface{}) (err error) {

// 	fmt.Printf("Warning!: POPUP/MENU implementation must be changed using new version\n")

// 	var menuItem *gtk.CheckMenuItem
// 	menuItem, err = gtk.CheckMenuItemNewWithMnemonic(lbl)

// 	// Handle the "activate" signal from the related item.
// 	if err == nil {
// 		menuItem.SetUseUnderline(true)
// 		if activateFunction != nil {
// 			menuItem.Connect("activate", activateFunction.(func()))
// 		}
// 		pop.Items = append(pop.Items, menuItem)
// 		// pop.separators = append(pop.separators, nil)
// 		pop.nextMenuItemIdx++
// 	}
// 	return err
// }

// // menuItemNewWithImage: Build an item containing an image.
// func (pop *PopupMenuIconStruct) menuItemNewWithImage(label string,
// 	pixbuf *gdk.Pixbuf) (menuItem *gtk.MenuItem, err error) {

// 	fmt.Printf("Warning!: POPUP/MENU implementation must be changed using new version\n")

// 	var image *gtk.Image
// 	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, pop.BoxSpacing)
// 	if err == nil {
// 		if image, err = gtk.ImageNewFromPixbuf(pixbuf); err == nil {
// 			label, err := gtk.LabelNewWithMnemonic(label)
// 			if err == nil {
// 				menuItem, err = gtk.MenuItemNew()
// 				if err == nil {
// 					label.SetHAlign(gtk.ALIGN_START)
// 					box.Add(image)
// 					box.PackEnd(label, true, true, 0)
// 					box.SetHAlign(gtk.ALIGN_START)
// 					menuItem.Container.Add(box)
// 					menuItem.ShowAll()
// 				}
// 			}
// 		}
// 	}
// 	return menuItem, err
// }

// AddSeparator: Add separator to menu.
// func (pop *PopupMenuIconStruct) AddSeparator(label ...string) (err error) {

// 	if separatorItem, err := gtk.SeparatorMenuItemNew(); err == nil {
// 		if len(label) > 0 {
// 			separatorItem.SetLabel(label[0])
// 		}
// 		pop.Items = append(pop.Items, separatorItem)
// 		pop.nextMenuItemIdx++
// 	}
// 	return err
// }

// MenuBuild: Build popupmenu.
// func (pop *PopupMenuIconStruct) MenuBuild() *gtk.Menu {
// 	var err error
// 	if pop.Menu, err = gtk.MenuNew(); err == nil {
// 		for _, menuItem := range pop.Items {
// 			switch m := menuItem.(type) {
// 			case *gtk.SeparatorMenuItem:
// 				pop.Menu.Append(m)
// 			case *gtk.MenuItem:
// 				pop.Menu.Append(m)
// 			case *gtk.CheckMenuItem:
// 				pop.Menu.Append(m)
// 			case *gtk.RadioMenuItem:
// 				pop.Menu.Append(m)
// 			}
// 		}
// 		pop.Menu.ShowAll()
// 	} else {
// 		log.Println("Popup menu creation error !")
// 		return nil
// 	}

// 	pop.Menu.SetReserveToggleSize(pop.SetReserveToggleSize)
// 	return pop.Menu
// }

// AppendToExistingMenu: append "MenuItems" to an existing "*gtk.Menu"
// Useful when you want to just add some entries to the context menu that
// already exist in a gtk.TextView or gtk.Entry by using "populate-popup"
// signal. Notice: GtkWidget > GtkMenu:
// menu := &gtk.Menu{gtk.MenuShell{gtk.Container{*w}}}
// Each connection to signal need to re-create the entire menu.
// func (pop *PopupMenuIconStruct) AppendToExistingMenu(menu *gtk.Menu) *gtk.Menu {
// 	for _, menuItem := range pop.Items {
// 		switch m := menuItem.(type) {
// 		case *gtk.SeparatorMenuItem:
// 			pop.Menu.Append(m)
// 		case *gtk.MenuItem:
// 			pop.Menu.Append(m)
// 		case *gtk.CheckMenuItem:
// 			pop.Menu.Append(m)
// 		case *gtk.RadioMenuItem:
// 			pop.Menu.Append(m)
// 		}
// 	}
// 	menu.ShowAll()
// 	return menu
// }

/******************** OUTDATED END **********************/
