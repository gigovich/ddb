package schema

import (
	"fmt"
	"reflect"
)

var (
	// ErrObjNotPointer ...
	ErrObjNotPointer = fmt.Errorf("scheme object should be pointer")

	// ErrFieldNotPointer ...
	ErrFieldNotPointer = fmt.Errorf("filed should be pointer")

	// ErrNoSuchField ...
	ErrNoSuchField = fmt.Errorf("no such field")
)

// Table ...
type Table struct {
	PK      []string
	Remaps  map[string]string
	Getters map[string]struct{}

	ot reflect.Type
	ov reflect.Value

	ft map[string]reflect.StructField
	fv map[string]reflect.Value
	fa map[reflect.Value]string

	obj interface{}
	err error
}

// Default scheme of database table
func Default(obj interface{}) Table {
	t := Table{}

	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		t.err = ErrObjNotPointer
		return t
	}

	t.PK = []string{"ID"}
	t.Remaps = make(map[string]string)
	t.Getters = map[string]struct{}{"ID": struct{}{}}

	t.obj = obj
	t.ot = reflect.TypeOf(obj).Elem()
	t.ov = reflect.ValueOf(obj).Elem()
	t.ft = make(map[string]reflect.StructField)
	t.fv = make(map[string]reflect.Value)
	t.fa = make(map[reflect.Value]string)

	for i := 0; i < t.ot.NumField(); i++ {
		t.ft[t.ot.Field(i).Name] = t.ot.Field(i)
		t.fv[t.ot.Field(i).Name] = t.ov.Field(i)
		t.fa[t.ov.Field(i).Addr()] = t.ot.Field(i).Name
	}

	return t
}

// WithFieldName set alternative name for field
func (t Table) WithFieldName(field interface{}, rename string) Table {
	if t.err != nil {
		return t
	}

	if reflect.TypeOf(field).Kind() != reflect.Ptr {
		t.err = ErrFieldNotPointer
		return t
	}

	name, ok := t.fa[reflect.ValueOf(field).Elem().Addr()]
	if !ok {
		t.err = ErrFieldNotPointer
		return t
	}

	t.Remaps[name] = rename
	return t
}

// WithPK can set primary keys
func (t Table) WithPK(fields ...interface{}) Table {
	if t.err != nil {
		return t
	}
	t.PK = t.PK[:0]

	for _, f := range fields {
		if reflect.TypeOf(f).Kind() != reflect.Ptr {
			t.err = ErrFieldNotPointer
			return t
		}

		name, ok := t.fa[reflect.ValueOf(f).Elem().Addr()]
		if !ok {
			t.err = ErrNoSuchField
			return t
		}

		t.PK = append(t.PK, name)
	}

	return t
}

// WithGetter adds get object query with limit to one.
func (t Table) WithGetter(field interface{}) Table {
	if t.err != nil {
		return t
	}

	if reflect.TypeOf(field).Kind() != reflect.Ptr {
		t.err = ErrFieldNotPointer
		return t
	}

	name, ok := t.fa[reflect.ValueOf(field).Elem().Addr()]
	if !ok {
		t.err = ErrFieldNotPointer
		return t
	}

	t.Getters[name] = struct{}{}
	return t
}

// Err return first error in the chain
func (t Table) Err() error {
	return t.err
}
