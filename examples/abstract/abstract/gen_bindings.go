// Code generated by tsgo. DO NOT EDIT.
package main

/*
#include <stdlib.h>
#include "helpers.h"
#include <stdint.h>

static inline size_t float32Size() {
  return sizeof(float);
}

static inline size_t float64Size() {
  return sizeof(double);
}

static inline size_t int32Size() {
  return sizeof(int32_t);
}

static inline size_t int64Size() {
  return sizeof(int64_t);
}

static inline size_t uint32Size() {
  return sizeof(uint32_t);
}

static inline size_t uint64Size() {
  return sizeof(uint64_t);
}
*/
import "C"

import(
  "github.com/cryptodeal/tsgo/examples/abstract"
  "unsafe"
  "fmt"
  "runtime/cgo"
  "encoding/json"
)

var ptrTrckr = make(map[unsafe.Pointer]C.size_t)

//export dispose
func dispose(ptr unsafe.Pointer, ctx unsafe.Pointer) {
  if _, ok := ptrTrckr[ptr]; ok {
    delete(ptrTrckr, ptr)
    defer C.free(ptr)
  } else {
    panic(fmt.Sprintf("Error: `%#v` not found in ptrTrckr", ptr))
  }
}

//export genDisposePtr
func genDisposePtr() unsafe.Pointer {
  return C.disposePtr
}

//export arraySize
func arraySize(ptr unsafe.Pointer) C.size_t {
  if val, ok := ptrTrckr[ptr]; ok {
    return val
  }
  panic(fmt.Sprintf("Error: `%#v` not found in ptrTrckr", ptr))
}

func CFloat32(b []float32) unsafe.Pointer {
  arr_len := len(b)
  p := C.malloc(C.size_t(arr_len) * C.float32Size())
  sliceHeader := struct {
    p   unsafe.Pointer
    len int
    cap int
  }{p, arr_len, arr_len}
  s := *(*[]float32)(unsafe.Pointer(&sliceHeader))
  copy(s, b)
  ptrTrckr[p] = C.size_t(arr_len)
  return p
}

func CFloat64(b []float64) unsafe.Pointer {
  arr_len := len(b)
  p := C.malloc(C.size_t(arr_len) * C.float64Size())
  sliceHeader := struct {
    p   unsafe.Pointer
    len int
    cap int
  }{p, arr_len, arr_len}
  s := *(*[]float64)(unsafe.Pointer(&sliceHeader))
  copy(s, b)
  ptrTrckr[p] = C.size_t(arr_len)
  return p
}

func CInt32(b []int32) unsafe.Pointer {
  arr_len := len(b)
  p := C.malloc(C.size_t(arr_len) * C.int32Size())
  sliceHeader := struct {
    p   unsafe.Pointer
    len int
    cap int
  }{p, arr_len, arr_len}
  s := *(*[]int32)(unsafe.Pointer(&sliceHeader))
  copy(s, b)
  ptrTrckr[p] = C.size_t(arr_len)
  return p
}

func CInt64(b []int64) unsafe.Pointer {
  arr_len := len(b)
  p := C.malloc(C.size_t(arr_len) * C.int64Size())
  sliceHeader := struct {
    p   unsafe.Pointer
    len int
    cap int
  }{p, arr_len, arr_len}
  s := *(*[]int64)(unsafe.Pointer(&sliceHeader))
  copy(s, b)
  ptrTrckr[p] = C.size_t(arr_len)
  return p
}

func CUint32(b []uint32) unsafe.Pointer {
  arr_len := len(b)
  p := C.malloc(C.size_t(arr_len) * C.uint32Size())
  sliceHeader := struct {
    p   unsafe.Pointer
    len int
    cap int
  }{p, arr_len, arr_len}
  s := *(*[]uint32)(unsafe.Pointer(&sliceHeader))
  copy(s, b)
  ptrTrckr[p] = C.size_t(arr_len)
  return p
}

func CUint64(b []uint64) unsafe.Pointer {
  arr_len := len(b)
  p := C.malloc(C.size_t(arr_len) * C.uint64Size())
  sliceHeader := struct {
    p   unsafe.Pointer
    len int
    cap int
  }{p, arr_len, arr_len}
  s := *(*[]uint64)(unsafe.Pointer(&sliceHeader))
  copy(s, b)
  ptrTrckr[p] = C.size_t(arr_len)
  return p
}

func encodeJSON(x interface{}) []byte {
  res, err := json.Marshal(x)
  if err != nil {
    fmt.Println(err)
    panic(err)
  }
  return res
}

//export _IntTest
func _IntTest(foo *C.char) C.int {
  _foo := C.GoString(foo)
  _returned_value := C.int(abstract.IntTest(_foo))
  return _returned_value
}

//export _Float32ArrayTest
func _Float32ArrayTest(foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CFloat32(abstract.Float32ArrayTest(_foo))
  return _returned_value
}

//export _Float64ArrayTest
func _Float64ArrayTest(foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CFloat64(abstract.Float64ArrayTest(_foo))
  return _returned_value
}

//export _Int32ArrayTest
func _Int32ArrayTest(foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CInt32(abstract.Int32ArrayTest(_foo))
  return _returned_value
}

//export _Int64ArrayTest
func _Int64ArrayTest(foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CInt64(abstract.Int64ArrayTest(_foo))
  return _returned_value
}

//export _Uint32ArrayTest
func _Uint32ArrayTest(foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CUint32(abstract.Uint32ArrayTest(_foo))
  return _returned_value
}

//export _Uint64ArrayTest
func _Uint64ArrayTest(foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CUint64(abstract.Uint64ArrayTest(_foo))
  return _returned_value
}

//export _StringTest
func _StringTest() *C.char {
  _returned_value := C.CString(abstract.StringTest())
  defer C.free(unsafe.Pointer(_returned_value))
  return _returned_value
}

//export _Float32ArgTest
func _Float32ArgTest(foo unsafe.Pointer, foo_len uint64) unsafe.Pointer {
  _foo := unsafe.Slice((*float32)(foo), foo_len)
  _returned_value := CFloat32(abstract.Float32ArgTest(_foo))
  return _returned_value
}

//export _Float64ArgTest
func _Float64ArgTest(foo unsafe.Pointer, foo_len uint64) unsafe.Pointer {
  _foo := unsafe.Slice((*float64)(foo), foo_len)
  _returned_value := CFloat64(abstract.Float64ArgTest(_foo))
  return _returned_value
}

//export _Int32ArgTest
func _Int32ArgTest(foo unsafe.Pointer, foo_len uint64) unsafe.Pointer {
  _foo := unsafe.Slice((*int32)(foo), foo_len)
  _returned_value := CInt32(abstract.Int32ArgTest(_foo))
  return _returned_value
}

//export _Int64ArgTest
func _Int64ArgTest(foo unsafe.Pointer, foo_len uint64) unsafe.Pointer {
  _foo := unsafe.Slice((*int64)(foo), foo_len)
  _returned_value := CInt64(abstract.Int64ArgTest(_foo))
  return _returned_value
}

//export _Uint32ArgTest
func _Uint32ArgTest(foo unsafe.Pointer, foo_len uint64) unsafe.Pointer {
  _foo := unsafe.Slice((*uint32)(foo), foo_len)
  _returned_value := CUint32(abstract.Uint32ArgTest(_foo))
  return _returned_value
}

//export _Uint64ArgTest
func _Uint64ArgTest(foo unsafe.Pointer, foo_len uint64) unsafe.Pointer {
  _foo := unsafe.Slice((*uint64)(foo), foo_len)
  _returned_value := CUint64(abstract.Uint64ArgTest(_foo))
  return _returned_value
}

//export _TestStruct
func _TestStruct() unsafe.Pointer {
  return C.hackyHandle(C.uintptr_t(cgo.NewHandle(abstract.TestStruct())))
}

//export _GET_StructBar_Field
func _GET_StructBar_Field(handle C.uintptr_t) *C.char {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.StructBar)
  _returned_value := C.CString(string(s.Field))
  defer C.free(unsafe.Pointer(_returned_value))
  return _returned_value
}

//export _GET_StructBar_FieldWithWeirdJSONTag
func _GET_StructBar_FieldWithWeirdJSONTag(handle C.uintptr_t) C.int64_t {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.StructBar)
  _returned_value := C.int64_t(int64(s.FieldWithWeirdJSONTag))
  return _returned_value
}

//export _GET_StructBar_FieldThatShouldBeOptional
func _GET_StructBar_FieldThatShouldBeOptional(handle C.uintptr_t) *C.char {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.StructBar)
  if s.FieldThatShouldBeOptional == nil {
    return nil
  }
  _returned_value := C.CString(string(*s.FieldThatShouldBeOptional))
  defer C.free(unsafe.Pointer(_returned_value))
  return _returned_value
}

//export _GET_StructBar_FieldThatShouldNotBeOptional
func _GET_StructBar_FieldThatShouldNotBeOptional(handle C.uintptr_t) *C.char {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.StructBar)
  _returned_value := C.CString(string(*s.FieldThatShouldNotBeOptional))
  defer C.free(unsafe.Pointer(_returned_value))
  return _returned_value
}

//export _GET_StructBar_FieldThatShouldBeReadonly
func _GET_StructBar_FieldThatShouldBeReadonly(handle C.uintptr_t) *C.char {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.StructBar)
  _returned_value := C.CString(string(s.FieldThatShouldBeReadonly))
  defer C.free(unsafe.Pointer(_returned_value))
  return _returned_value
}

//export _GET_StructBar_ArrayField
func _GET_StructBar_ArrayField(handle C.uintptr_t) unsafe.Pointer {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.StructBar)
  if s.ArrayField == nil {
    return nil
  }
  _returned_value := unsafe.Pointer(CFloat32(s.ArrayField))
  return _returned_value
}

//export _GET_StructBar_StructField
func _GET_StructBar_StructField(handle C.uintptr_t) unsafe.Pointer {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.StructBar)
  if s.StructField == nil {
    return nil
  }
  return C.hackyHandle(C.uintptr_t(cgo.NewHandle(*s.StructField)))
}

//export _GET_DemoStruct_ArrayField
func _GET_DemoStruct_ArrayField(handle C.uintptr_t) unsafe.Pointer {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.DemoStruct)
  if s.ArrayField == nil {
    return nil
  }
  _returned_value := unsafe.Pointer(CFloat32(*s.ArrayField))
  return _returned_value
}

//export _GET_DemoStruct_FieldToAnotherStruct
func _GET_DemoStruct_FieldToAnotherStruct(handle C.uintptr_t) unsafe.Pointer {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.DemoStruct)
  if s.FieldToAnotherStruct == nil {
    return nil
  }
  return C.hackyHandle(C.uintptr_t(cgo.NewHandle(*s.FieldToAnotherStruct)))
}

//export _GET_DemoStruct2_AnotherArray
func _GET_DemoStruct2_AnotherArray(handle C.uintptr_t) unsafe.Pointer {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.DemoStruct2)
  if s.AnotherArray == nil {
    return nil
  }
  _returned_value := unsafe.Pointer(CFloat64(*s.AnotherArray))
  return _returned_value
}

//export _GET_DemoStruct2_BacktoAnotherStruct
func _GET_DemoStruct2_BacktoAnotherStruct(handle C.uintptr_t) unsafe.Pointer {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.DemoStruct2)
  if s.BacktoAnotherStruct == nil {
    return nil
  }
  return C.hackyHandle(C.uintptr_t(cgo.NewHandle(*s.BacktoAnotherStruct)))
}

//export _GET_DemoStruct3_AnotherArray
func _GET_DemoStruct3_AnotherArray(handle C.uintptr_t) unsafe.Pointer {
  h := cgo.Handle(handle)
  s := h.Value().(abstract.DemoStruct3)
  if s.AnotherArray == nil {
    return nil
  }
  _returned_value := unsafe.Pointer(CFloat32(*s.AnotherArray))
  return _returned_value
}

//export _DISPOSE_Struct
func _DISPOSE_Struct(handle C.uintptr_t) {
  h := cgo.Handle(handle)
  fmt.Println("deleted handle @ uintptr:", handle)
  h.Delete()
}

//export _TestStruct2
func _TestStruct2() unsafe.Pointer {
  return C.hackyHandle(C.uintptr_t(cgo.NewHandle(*abstract.TestStruct2())))
}

//export _TestMap
func _TestMap() *C.char {
  _temp_res_val := encodeJSON(*abstract.TestMap())
  _returned_value := C.CString(string(_temp_res_val))
  defer C.free(unsafe.Pointer(_returned_value))
  return _returned_value
}

func main() {} // Required but ignored