package reflection

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReflections(t *testing.T) {
	a := &struct {
		A int64
		B float64
		C string
	}{}
	// check type
	require.True(t, IsType(a.A, int64(42)), "type error")
	// check invalid type
	require.False(t, IsType(42, nil), "type error")
	// field type
	ok, err := FieldIsType(a, "B", float64(3.14))
	require.NoError(t, err, "unexpected error")
	require.True(t, ok, "type error")
	ok, err = FieldIsType(*a, "C", string(""))
	require.NoError(t, err, "unexpected error")
	require.True(t, ok, "type error")
	// replace fields
	b, err := ReplaceFieldsType(a, FieldReplacementMap{
		"A": uint64(42),
		"B": float32(3.14),
	})
	require.NoError(t, err, "can't replace types")
	// check field types
	ok, err = FieldIsType(b, "A", uint64(42))
	require.NoError(t, err, "unexpected error")
	require.True(t, ok, "type error")
	ok, err = FieldIsType(b, "B", float32(3.14))
	require.NoError(t, err, "unexpected error")
	require.True(t, ok, "type error")
	ok, err = FieldIsType(b, "C", "")
	require.NoError(t, err, "unexpected error")
	require.True(t, ok, "type error")
	// filter fields
	c, err := FilterFields(a, "A", "C")
	require.NoError(t, err, "can't filter")
	// check field existence
	ok, err = HasField(c, "A")
	require.NoError(t, err, "unexpected error")
	require.True(t, ok, "A is missing")
	ok, err = HasField(c, "C")
	require.NoError(t, err, "unexpected error")
	require.True(t, ok, "C is missing")
	ok, err = HasField(c, "B")
	require.Error(t, err, "expecting an error")
	require.False(t, ok, "B is present")
	// build d
	d := NewStructBuilder().WithFields(
		&StructField{Name: "A", Type: int64(42), Tag: `hello:"world"`},
		&StructField{Name: "B", Type: float32(3.14), Tag: `ola:"mundo"`},
	).BuildPointer()
	err = CopyFields(c, d)
	require.NoError(t, err, "can't copy fields")
	// check field values
	fld1, err := Field(c, "A")
	require.NoError(t, err, "unexpected error")
	fld2, err := Field(d, "A")
	require.NoError(t, err, "unexpected error")
	require.Equal(t, fld1, fld2, "field mismatch")
}
