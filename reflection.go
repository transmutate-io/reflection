package reflection

import (
	"errors"
	"reflect"
)

var (
	// ErrNilPointer is returned when a nil pointer is found
	ErrNilPointer = errors.New("nil pointer")
	// ErrNotAPointer is returned when passing an argument that is not a pointer
	ErrNotAPointer = errors.New("not a pointer")
	// ErrNotAStruct is returned when passing an argument that is not a struct
	ErrNotAStruct = errors.New("not a struct")
	// ErrFieldNotFound is returned when a field is not found
	ErrFieldNotFound = errors.New("field not found")
)

// Any empty interface
type Any interface{}

// MustCopyFields calls CopyFields and panics on error
func MustCopyFields(src, dst Any) {
	if err := CopyFields(src, dst); err != nil {
		panic(err)
	}
}

func typeOfStruct(t Any, mustPointer bool) (reflect.Type, bool, error) {
	r, p, err := valueOfStruct(t, mustPointer)
	if err != nil {
		return nil, false, err
	}
	return r.Type(), p, nil
}

func valueOfStruct(t Any, mustPointer bool) (reflect.Value, bool, error) {
	// is t nil
	if t == nil {
		return reflect.Value{}, false, ErrNilPointer
	}
	// is t a pointer
	tv := reflect.ValueOf(t)
	if mustPointer && tv.Kind() != reflect.Ptr {
		return reflect.Value{}, false, ErrNotAPointer
	}
	var isPointer bool
	if tv.Kind() == reflect.Ptr {
		isPointer = true
		if tv.IsNil() {
			return reflect.Value{}, false, ErrNilPointer
		}
		tv = tv.Elem()
	}
	// is t a struct
	if tv.Kind() != reflect.Struct {
		return reflect.Value{}, false, ErrNotAStruct
	}
	return tv, isPointer, nil
}

// CopyFields copies any existing fields of the same type or a type that
// implements an interface, from src to dst
func CopyFields(src, dst Any) error {
	// is src or dst nil
	if src == nil || dst == nil {
		return ErrNilPointer
	}
	// source
	vs, _, err := valueOfStruct(src, false)
	if err != nil {
		return err
	}
	vd, _, err := valueOfStruct(dst, true)
	if err != nil {
		return err
	}
	// copy only fields existing on dst
	vst := vs.Type()
	vdt := vd.Type()
	for i := 0; i < vdt.NumField(); i++ {
		// get dst field
		dfld := vdt.Field(i)
		// get src field by name
		sfld, ok := vst.FieldByName(dfld.Name)
		if !ok {
			continue
		}
		// not the same type
		if sfld.Type != dfld.Type {
			// is dst field an interface
			if dfld.Type.Kind() != reflect.Interface {
				continue
			}
			// does src field implement dst field
			if !sfld.Type.Implements(dfld.Type) {
				continue
			}
		}
		// set value
		vd.FieldByName(dfld.Name).Set(vs.Field(i))
	}
	return nil
}

// FieldReplacementMap represents a set of substitutions to be performed on a struct type
type FieldReplacementMap = map[string]Any

// MustReplaceTypeFields calls ReplaceFieldsType and panics on error
func MustReplaceTypeFields(t Any, rm FieldReplacementMap) Any {
	r, err := ReplaceFieldsType(t, rm)
	if err != nil {
		panic(err)
	}
	return r
}

// ReplaceFieldsType generates a struct using t as a model and performing the
// substitutions provided in replMap
func ReplaceFieldsType(t Any, replMap FieldReplacementMap) (Any, error) {
	tt, isPointer, err := typeOfStruct(t, false)
	if err != nil {
		return nil, err
	}
	// new fields
	newFields := make([]reflect.StructField, 0, tt.NumField())
	for i := 0; i < tt.NumField(); i++ {
		fld := tt.Field(i)
		if newType, ok := replMap[fld.Name]; ok {
			fld.Type = reflect.TypeOf(newType)
		}
		newFields = append(newFields, fld)
	}
	tr := reflect.New(reflect.StructOf(newFields))
	if !isPointer {
		tr = tr.Elem()
	}
	return tr.Interface(), nil
}

// MustFilterFields calls FilterFields and panics on error
func MustFilterFields(t Any, flds ...string) Any {
	r, err := FilterFields(t, flds...)
	if err != nil {
		panic(err)
	}
	return r
}

// FilterFields generates a struct using t as a model with only the fields provided
func FilterFields(t Any, flds ...string) (Any, error) {
	tt, isPointer, err := typeOfStruct(t, false)
	if err != nil {
		return nil, err
	}
	// new fields
	newFields := make([]reflect.StructField, 0, len(flds))
	for _, i := range flds {
		fld, ok := tt.FieldByName(i)
		if !ok {
			return nil, ErrFieldNotFound
		}
		newFields = append(newFields, fld)
	}
	tr := reflect.New(reflect.StructOf(newFields))
	// return a pointer only if t is a pointer
	if !isPointer {
		tr = tr.Elem()
	}
	return tr.Interface(), nil
}

// MustHasField calls HasField and panics on error
func MustHasField(t Any, name string) bool {
	r, err := HasField(t, name)
	if err != nil {
		panic(err)
	}
	return r
}

// HasField checks if a field with a the provided name exists
func HasField(t Any, name string) (bool, error) {
	if _, err := Field(t, name); err != nil {
		return false, err
	}
	return true, nil
}

func isType(t reflect.Type, _type Any) bool {
	return t != nil &&
		_type != nil &&
		t.String() == reflect.TypeOf(_type).String()
}

// IsType compares the type of t with _type
func IsType(t Any, _type Any) bool {
	return isType(reflect.TypeOf(t), _type)
}

// MustField calls Field and panics on error
func MustField(t Any, name string) Any {
	r, err := Field(t, name)
	if err != nil {
		panic(err)
	}
	return r
}

// Field returns a field by name
func Field(t Any, name string) (Any, error) {
	tv, _, err := valueOfStruct(t, false)
	if err != nil {
		return nil, err
	}
	// get field
	fld := tv.FieldByName(name)
	if !fld.IsValid() {
		return nil, ErrFieldNotFound
	}
	return fld.Interface(), nil
}

// MustFieldIsType calls FieldIsType and panics on error
func MustFieldIsType(t Any, name string, _type Any) (bool, error) {
	r, err := FieldIsType(t, name, _type)
	if err != nil {
		return false, err
	}
	return r, nil
}

// FieldIsType compares the type of a field  with _type
func FieldIsType(t Any, name string, _type Any) (bool, error) {
	tt, _, err := typeOfStruct(t, false)
	if err != nil {
		return false, err
	}
	fld, ok := tt.FieldByName(name)
	if !ok {
		return false, err
	}
	return isType(fld.Type, _type), nil
}
