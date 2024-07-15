package main

import (
	"fmt"
	"reflect"
)

// ConvertTo converts a value of T1 type to value of T2 type
//
// It panics if one of followings met:
//
//	If T1's Kind is not a Pointer to Struct or T2's Kind is not a struct
//	If T1 and T2 does not share the same layout memory (ref examples below for details)
//	If T2 is larger than T1
//	If T1's underlying type is T2
//	Currently if T1 or T2's underlying field Kind is not one of following: Bool, Int*, Float*, Uint, Uint8, Uint16, Uint32, Uint64, Float*
//
// Pros and cons of using this function are
//
//	Pros:
//		Avoiding allocating memory so that reducing pressure on GC
//		Make your program faster
//	Cons:
//		Dangerous to use in real world apps so please use this if you know what you're actually doing
//		No modifications on the result returned because of also impacting on the input
//
// For example:
//
//	type T1 struct {
//		t1 int8
//		t2 int32
//		t3 string
//	}
//
//	type T2 struct {
//		tt1 int8
//		tt2 int32
//	}
//
// Note that T1 and T2 share the same layout as T2 has 2 ordered properties as T1 does
// regardless of differences of property's name and we should use the func for this case,
//
// The followings are prohibited:
//
//	type T1 struct {
//		t1 int8
//		t2 int32
//		t3 string
//	}
//
//	type T2 struct {
//		tt1 int32 --> should be int8 as t1 of T1
//		tt2 int8 --> should be int32 as t2 of T1
//	}
//
// although both structs share the same layout, T2 must not be larger than T1, the below is also avoided
//
//	type T1 struct {
//		tt1 int8
//		tt2 int32
//	}
//
//	type T2 struct {
//		t1 int8
//		t2 int32
//		t3 string
//	}
//
// type T1 T2 is also avoided since it's unnecessary to use this func instead of usual castings
func ConvertTo[T1, T2 any](t1 T1) (t2 T2) {
	ttyp1, ttyp2 := reflect.TypeOf(t1), reflect.TypeOf(t2)

	if ttyp1.Kind() != reflect.Pointer || ttyp1.Elem().Kind() != reflect.Struct {
		panic("T1 must be a pointer to struct")
	}
	if ttyp2.Kind() != reflect.Struct {
		panic("T2 must be a struct")
	}

	elmTyp1, elmTyp2 := ttyp1.Elem(), ttyp2
	if elmTyp1.ConvertibleTo(elmTyp2) {
		panic("T1 must not be a type definition of T2 or vice versa")
	}
	if elmTyp1.NumField() < elmTyp1.NumField() {
		panic("T2 must be no larger than T1")
	}
	for i := 0; i < elmTyp2.NumField(); i++ {
		ftyp2 := elmTyp2.Field(i).Type
		ftyp1 := elmTyp1.Field(i).Type
		if ftyp1.Kind() != ftyp2.Kind() {
			panic("T2 and T1 must have the same field type")
		}
		if ftyp1.Kind() == reflect.Pointer {
			ftyp1 = ftyp1.Elem()
		}
		if ftyp2.Kind() == reflect.Pointer {
			ftyp2 = ftyp2.Elem()
		}
		if !isValidKind(ftyp1.Kind()) || !isValidKind(ftyp2.Kind()) {
			panic("T2 or T1 must be a valid type")
		}
	}
	t2 = *(*T2)(reflect.ValueOf(t1).UnsafePointer())
	return
}

var validKinds = [...]reflect.Kind{
	reflect.Bool,
	reflect.Int,
	reflect.Int8,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Float32,
	reflect.Float64,
}

func isValidKind(kind reflect.Kind) bool {
	for _, validKind := range validKinds {
		if kind == validKind {
			return true
		}
	}
	return false
}

type Type1 struct {
	t1 int8
	t2 *int32
	t3 string
}
type Type2 struct {
	tt1 int8
	tt2 *int32
}

func main() {
	two := int32(2)
	t1 := Type1{t1: 1, t2: &two, t3: "test"}
	fmt.Println(*t1.t2)
	t2 := ConvertTo[*Type1, Type2](&t1)
	fmt.Println(*t2.tt2)
}
