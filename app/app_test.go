package app

import (
	"testing"
	_ "reflect"
	deepcopy "github.com/mitchellh/copystructure"
	"reflect"
	"time"
)

// Iface is an alias to Copy; this exists for backwards compatibility reasons.
func Iface(iface interface{}) interface{} {
	return Copy(iface)
}

// Copy creates a deep copy of whatever is passed to it and returns the copy
// in an interface{}.  The returned value will need to be asserted to the
// correct type.
func Copy(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	// Make the interface a reflect.Value
	original := reflect.ValueOf(src)

	// Make a copy of the same type as the original.
	cpy := reflect.New(original.Type()).Elem()

	// Recursively copy the original.
	copyRecursive(original, cpy)

	// Return theb copy as an interface.
	return cpy.Interface()
}

// copyRecursive does the actual copying of the interface. It currently has
// limited support for what it can handle. Add as needed.
func copyRecursive(original, cpy reflect.Value) {
	// handle according to original's Kind
	switch original.Kind() {
	case reflect.Ptr:
		// Get the actual value being pointed to.
		originalValue := original.Elem()

		// if  it isn't valid, return.
		if !originalValue.IsValid() {
			return
		}
		cpy.Set(reflect.New(originalValue.Type()))
		copyRecursive(originalValue, cpy.Elem())

	case reflect.Interface:
		// If this is a nil, don't do anything
		if original.IsNil() {
			return
		}
		// Get the value for the interface, not the pointer.
		originalValue := original.Elem()

		// Get the value by calling Elem().
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue)
		cpy.Set(copyValue)

	case reflect.Struct:
		t, ok := original.Interface().(time.Time)
		if ok {
			b, err := t.MarshalBinary()
			if err != nil {
				panic(err)
			}
			tcpy := time.Time{}
			tcpy.UnmarshalBinary(b)
			if err != nil {
				panic(err)
			}
			cpy.Set(reflect.ValueOf(tcpy))
			return
		}
		// Go through each field of the struct and copy it.
		for i := 0; i < original.NumField(); i++ {
			// The Type's StructField for a given field is checked to see if StructField.PkgPath
			// is set to determine if the field is exported or not because CanSet() returns false
			// for settable fields.  I'm not sure why.  -mohae
			if original.Type().Field(i).PkgPath != "" {
				continue
			}
			copyRecursive(original.Field(i), cpy.Field(i))
		}

	case reflect.Slice:
		if original.IsNil() {
			return
		}
		// Make a new slice and copy each element.
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), cpy.Index(i))
		}

	case reflect.Map:
		if original.IsNil() {
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue)
			cpy.SetMapIndex(key, copyValue)
		}

	default:
		cpy.Set(original)
	}
}


type obj struct {
	i int
	c string
	M map[string]string
	m map[string]string
}

func TestDeepCopy(t *testing.T)  {
	new1 := obj{i:1, c:"cccc", M:map[string]string{"xx": "xxx"}}
	new1.M["xxx"] = "xxxxx"

	//v := reflect.New(reflect.TypeOf(new1))

	//var new2i interface{}
	//if reflect.TypeOf(new1).Kind() == reflect.Ptr{
	//	new2i, _ = deepcopy.Copy(*new1)
	//}else{
	//	new2i, _ = deepcopy.Copy(new1)
	//}

	//t.Log(reflect.TypeOf(new1).Kind())
	//
	new2i, _ := deepcopy.Copy(new1)
	new2 := new2i.(obj)
	//new2.c = "ooooo"
	new2.i = 1111
	println(new2.M)
	new2.M["ttt"] = "xxxx"

	new3 := Copy(new1)

	t.Log("new1:", new1)
	t.Log("new2:", new2)
	t.Log("new3:", new3)
}