// widgets.go

/*
	Source file auto-generated on Sat, 19 Oct 2019 22:06:16 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019 H.F.M (github.com/hfmrow)
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"reflect"

	"github.com/gotk3/gotk3/gtk"
)

// WidgetProperties: A convenient structure that implements
// methods to add multiple properties to an object at a time,
// considering a property restriction list that contains
// unauthorized properties for certain types of objects.
type WidgetProperties struct {
	properties       []widgetProperty
	propRestrictions [][]interface{} // List of Property by object that not usable
	initialized      bool
}
type widgetProperty struct {
	Name  string
	Value interface{}
}

// PropsToWidget: Apply the properties to the object after filtering
// unauthorized ones. Can be used for independent objects, such as those
// built into an independant box that will be added as a container, to
// the Dialog's list of the main objects.
func (wp *WidgetProperties) PropsToWidget(wdg gtk.IWidget) {
	wp.init()
	for _, prop := range wp.properties {
		for _, obj := range wp.propRestrictions {
			if reflect.TypeOf(wdg) == reflect.TypeOf(obj[0]) {
				for _, avoidProp := range obj[1:] {
					if prop.Name == avoidProp.(string) {
						return
					}
				}
			}
		}
		wdg.Set(prop.Name, prop.Value)
	}
}

// AddPropertyUnallowed: add new unallowed property to the list
// i.e: []interface{}{new(gtk.Image), "justify", "wrap", "pattern", "relief"}
// means that the gtk.Image object cannot handle:
// "justify", "wrap", "pattern", "relief" properties.
func (wp *WidgetProperties) AddPropertyUnallowed(value []interface{}) {
	wp.init()
	// Add new property for an object to the list that will be applied during control' creation.
	wp.propRestrictions = append(wp.propRestrictions, value)
}

// AddProperty: add new property to the list, will be applied to all parent objects
func (wp *WidgetProperties) AddProperty(name string, value interface{}) {
	wp.init()
	// Add new property for an object to the list that will be applied during control' creation.
	wp.properties = append(wp.properties, widgetProperty{Name: name, Value: value})
}

// Reset: user defined lists
func (wp *WidgetProperties) Reset(wdg gtk.IWidget) {
	wp.properties = wp.properties[:0]
	wp.propRestrictions = wp.propRestrictions[:0]
	wp.initialized = false
	wp.init()
}

// init: default unallowed properties for objects
func (wp *WidgetProperties) init() {
	if !wp.initialized {
		wp.propRestrictions = [][]interface{}{
			{new(gtk.Image), "justify", "wrap", "pattern", "relief"},
			{new(gtk.Separator), "halign", "justify", "wrap", "pattern", "relief"},
			{new(gtk.Calendar), "wrap", "pattern", "relief"},
			{new(gtk.Box), "wrap", "pattern", "relief"}}
		wp.initialized = true
	}
}
