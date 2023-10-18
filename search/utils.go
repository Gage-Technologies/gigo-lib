package search

import (
	"fmt"
	"reflect"
	"time"
)

// formatInterfaceToFilterString
//
//		Formats an interface value to a string of the base type
//	 of its value. Pointers are formatted as their underlying
//	 value. Custom types are formatted as their underlying type.
//	 Strings can be optionally wrapped in single quotes.
//
// Args:
//
//		x (interface): The value to be formatted
//	 wrapString(boo): Whether to wrap strings in single quotes
//
// Returns:
//
//	(string) The string representation of the value
func formatInterfaceToString(x interface{}, wrapString bool) string {
	// return empty string for nil interface
	if x == nil {
		return ""
	}

	// create value for interface
	val := reflect.ValueOf(x)

	// get the underlying type of the value
	valType := val.Type()
	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
		val = val.Elem()
	}

	// handle each macro type from Go and then handle the rest with default
	switch valType.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("%v", val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", val.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", val.Float())
	case reflect.String:
		// conditionally wrap string in single quotes
		if wrapString {
			return fmt.Sprintf("'%s'", val.String())
		}
		return fmt.Sprintf("%s", val.String())
	case reflect.Struct:
		// handle time
		// NOTE: if anyone sees this and knows a better way to do this please fix this
		t, ok := x.(time.Time)
		if ok {
			return fmt.Sprintf("%d", t.Unix())
		}
		tp, ok := x.(*time.Time)
		if ok {
			return fmt.Sprintf("%d", tp.Unix())
		}
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
