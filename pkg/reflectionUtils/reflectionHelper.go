package reflectionutils

import "reflect"

// Get pointer to any value
func ToPointer[T any](v T) *T {
	return &v
}

// Get base element value
func GetElementValue(val reflect.Value) reflect.Value {
	switch val.Kind() {
	case reflect.Interface:
		return GetElementValue(reflect.ValueOf(val.Interface()))
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
		return GetElementValue(val.Elem())
	default:
		return val
	}
}

// Get base element type
func GetElementType(tType reflect.Type) reflect.Type {
	switch tType.Kind() {
	case reflect.Interface:
		return GetElementType(reflect.TypeOf(tType))
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
		return GetElementType(tType.Elem())
	default:
		return tType
	}
}

func StringPointerType() reflect.Type {
	emptyString := ""
	return reflect.TypeOf(&emptyString)
}
