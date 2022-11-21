// Code generated by tsgo. DO NOT EDIT.
package main

/*
#include <stdlib.h>

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

import (
  "github.com/cryptodeal/tsgo/examples/abstract"
  "unsafe"
  "fmt"
  "encoding/json"
)

var ptrTrckr = make(map[uintptr]C.size_t)

//export dispose
func dispose(ptr unsafe.Pointer, ctx unsafe.Pointer) {
  ptr_num := uintptr(ptr)
  if _, ok := ptrTrckr[ptr_num]; ok {
    delete(ptrTrckr, ptr_num)
    defer C.free(ptr)
  } else {
    fmt.Sprintf("panic(Error: pointer `%d` not found in ptrTrckr", ptr_num)
  }
}

func genDisposePtr() unsafe.Pointer {
  return C.disposePtr
}

//export ArraySize
func ArraySize(array unsafe.Pointer) C.size_t {
  ptr_num := uintptr(array)
  if val, ok := ptrTrckr[ptr_num]; ok {
    return val
  }
  fmt.Sprintf("panic(Error: pointer `%d` not found in ptrTrckr"), ptr_num)
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
  ptrTrckr[uintptr(p)] = C.size_t(arr_len)
  return p
}

//export dispose
func dispose(ptr unsafe.Pointer, ctx unsafe.Pointer) {
  ptr_num := uintptr(ptr)
  if _, ok := ptrTrckr[ptr_num]; ok {
    delete(ptrTrckr, ptr_num)
    defer C.free(ptr)
  } else {
    fmt.Sprintf("panic(Error: pointer `%d` not found in ptrTrckr"), ptr_num)
  }
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
  ptrTrckr[uintptr(p)] = C.size_t(arr_len)
  return p
}

//export dispose
func dispose(ptr unsafe.Pointer, ctx unsafe.Pointer) {
  ptr_num := uintptr(ptr)
  if _, ok := ptrTrckr[ptr_num]; ok {
    delete(ptrTrckr, ptr_num)
    defer C.free(ptr)
  } else {
    fmt.Sprintf("panic(Error: pointer `%d` not found in ptrTrckr"), ptr_num)
  }
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
  ptrTrckr[uintptr(p)] = C.size_t(arr_len)
  return p
}

//export dispose
func dispose(ptr unsafe.Pointer, ctx unsafe.Pointer) {
  ptr_num := uintptr(ptr)
  if _, ok := ptrTrckr[ptr_num]; ok {
    delete(ptrTrckr, ptr_num)
    defer C.free(ptr)
  } else {
    fmt.Sprintf("panic(Error: pointer `%d` not found in ptrTrckr"), ptr_num)
  }
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
  ptrTrckr[uintptr(p)] = C.size_t(arr_len)
  return p
}

//export dispose
func dispose(ptr unsafe.Pointer, ctx unsafe.Pointer) {
  ptr_num := uintptr(ptr)
  if _, ok := ptrTrckr[ptr_num]; ok {
    delete(ptrTrckr, ptr_num)
    defer C.free(ptr)
  } else {
    fmt.Sprintf("panic(Error: pointer `%d` not found in ptrTrckr"), ptr_num)
  }
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
  ptrTrckr[uintptr(p)] = C.size_t(arr_len)
  return p
}

//export dispose
func dispose(ptr unsafe.Pointer, ctx unsafe.Pointer) {
  ptr_num := uintptr(ptr)
  if _, ok := ptrTrckr[ptr_num]; ok {
    delete(ptrTrckr, ptr_num)
    defer C.free(ptr)
  } else {
    fmt.Sprintf("panic(Error: pointer `%d` not found in ptrTrckr"), ptr_num)
  }
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
  ptrTrckr[uintptr(p)] = C.size_t(arr_len)
  return p
}

//export dispose
func dispose(ptr unsafe.Pointer, ctx unsafe.Pointer) {
  ptr_num := uintptr(ptr)
  if _, ok := ptrTrckr[ptr_num]; ok {
    delete(ptrTrckr, ptr_num)
    defer C.free(ptr)
  } else {
    fmt.Sprintf("panic(Error: pointer `%d` not found in ptrTrckr"), ptr_num)
  }
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
 func _IntTest (foo *C.char) C.int {
  _foo := C.GoString(foo)
  _returned_value := C.int(abstract.IntTest(_foo))
  return _returned_value
}

//export _Float32ArrayTest
 func _Float32ArrayTest (foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CFloat32(abstract.Float32ArrayTest(_foo))
  return _returned_value
}

//export _Float64ArrayTest
 func _Float64ArrayTest (foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CFloat64(abstract.Float64ArrayTest(_foo))
  return _returned_value
}

//export _Int32ArrayTest
 func _Int32ArrayTest (foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CInt32(abstract.Int32ArrayTest(_foo))
  return _returned_value
}

//export _Int64ArrayTest
 func _Int64ArrayTest (foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CInt64(abstract.Int64ArrayTest(_foo))
  return _returned_value
}

//export _Uint32ArrayTest
 func _Uint32ArrayTest (foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CUint32(abstract.Uint32ArrayTest(_foo))
  return _returned_value
}

//export _Uint64ArrayTest
 func _Uint64ArrayTest (foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CUint64(abstract.Uint64ArrayTest(_foo))
  return _returned_value
}

//export _StringTest
 func _StringTest () *C.char {
  _returned_value := C.CString(abstract.StringTest())
  defer C.free(unsafe.Pointer(_returned_value))
  return _returned_value
}

//export _ArrayArgTest
 func _ArrayArgTest (foo unsafe.Pointer, _len int) unsafe.Pointer {
  _foo := unsafe.Slice((*float64)(foo), _len)
  _returned_value := CFloat64(abstract.ArrayArgTest(_foo))
  return _returned_value
}

//export _TestStruct
 func _TestStruct () *C.char {
  _temp_res_val := encodeJSON(abstract.TestStruct())
  _returned_value := C.CString(string(_temp_res_val))
  defer C.free(unsafe.Pointer(_returned_value))
  return _returned_value
}

//export _TestMap
 func _TestMap () *C.char {
  _temp_res_val := encodeJSON(abstract.TestMap())
  _returned_value := C.CString(string(_temp_res_val))
  defer C.free(unsafe.Pointer(_returned_value))
  return _returned_value
}

func main() {} // Required but ignored