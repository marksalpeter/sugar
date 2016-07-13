package sugar

import (
	"encoding"
	"fmt"
	"reflect"
)

// Copy copies a into b
func Copy(a, b interface{}) error {
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	return copy(aValue, bValue)
}

func copy(a, b reflect.Value) error {
	// make sure we're always copying the same type of thing
	if a.Kind() != b.Kind() {
		return fmt.Errorf("tried to copy incompatible types")
	}

	// marshal and unmarshal data via the binary marshaler/unmarshaler
	if a.CanInterface() && b.CanInterface() {
		// make sure we're not copying to or from a nil value
		if a.Kind() == reflect.Ptr && a.IsNil() {
			return nil
		} else if b.Kind() == reflect.Ptr && b.IsNil() && b.CanSet() {
			b.Set(reflect.New(b.Type().Elem()))
		}

		aMarshaler, aOK := a.Interface().(encoding.BinaryMarshaler)
		bUnmarshaler, bOK := b.Interface().(encoding.BinaryUnmarshaler)
		if aOK && bOK {
			if bs, err := aMarshaler.MarshalBinary(); err == nil {
				if err := bUnmarshaler.UnmarshalBinary(bs); err == nil {
					return nil
				}
			}
		}
	}

	switch a.Kind() {
	case reflect.Ptr:
		if a.IsNil() {
			return nil
		} else if b.IsNil() && b.CanSet() {
			b.Set(reflect.New(b.Type().Elem()))
		}
		return copy(a.Elem(), b.Elem())
	case reflect.Slice:
		b.Set(reflect.MakeSlice(a.Type(), a.Len(), a.Cap()))
		for i, l := 0, a.Len(); i < l; i++ {
			if err := copy(a.Index(i), b.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Array:
		b.Set(reflect.New(reflect.ArrayOf(a.Len(), a.Type())))
		for i, l := 0, a.Len(); i < l; i++ {
			if err := copy(a.Index(i), b.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Struct:
		for i, l := 0, a.Type().NumField(); i < l; i++ {
			copy(a.Field(i), b.Field(i))
		}
	case reflect.Map:
		for _, key := range a.MapKeys() {
			b.SetMapIndex(key, a.MapIndex(key))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.String,
		reflect.Bool:
		if b.CanSet() {
			b.Set(a)
		}
	}

	return nil
}
