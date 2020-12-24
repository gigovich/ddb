package dsl

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

// PK fields list
type PK []interface{}

// Remaps of the field names
type Remaps map[interface{}]string

// GetBy query single object by fields names or group of fields
type GetBy [][]interface{}

// ListBy query object list by fields names or gourp of fileds
type ListBy [][]interface{}

// Table ...
type Table struct {
	PK     PK
	Remaps Remaps
	GetBy  GetBy
	ListBy ListBy

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

	t.PK = PK{}
	t.Remaps = Remaps{}
	t.GetBy = GetBy{}
	t.ListBy = ListBy{}

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
