package reflection

import (
	"reflect"
)

type structBuilderField struct {
	Type interface{}
	Tag  string
}

type structBuilder struct {
	Fields map[string]*structBuilderField
}

func NewStructBuilder() *structBuilder {
	return &structBuilder{Fields: make(map[string]*structBuilderField, 16)}
}

func (sb *structBuilder) WithField(name string, _type interface{}, tag string) *structBuilder {
	sb.Fields[name] = &structBuilderField{Type: _type, Tag: tag}
	return sb
}

type StructField struct {
	Name string
	Type interface{}
	Tag  string
}

func (sb *structBuilder) WithFields(flds ...*StructField) *structBuilder {
	for _, i := range flds {
		sb.WithField(i.Name, i.Type, i.Tag)
	}
	return sb
}

func (sb *structBuilder) Build() interface{} {
	return sb.buildPointer().Elem().Interface()
}

func (sb *structBuilder) BuildPointer() interface{} {
	return sb.buildPointer().Interface()
}

func (sb *structBuilder) buildPointer() reflect.Value {
	f := make([]reflect.StructField, 0, len(sb.Fields))
	for name, fld := range sb.Fields {
		f = append(f, reflect.StructField{
			Name: name,
			Type: reflect.TypeOf(fld.Type),
			Tag:  reflect.StructTag(fld.Tag),
		})
	}
	return reflect.New(reflect.StructOf(f))
}
