// // Code generated by tsgo. DO NOT EDIT.
package main

/*
#include <stdlib.h>
*/
import "C"

import (
  "github.com/cryptodeal/tsgo/examples/abstract"
  "unsafe"
)


var ptrTrckr = make(map[uintptr]C.size_t)

func CFloat32(b []float32) unsafe.Pointer {
  p := C.malloc(C.size_t(len(b)))
  sliceHeader := struct {
    p   unsafe.Pointer
    len int
    cap int
  }{p, len(b), len(b)}
  s := *(*[]float32)(unsafe.Pointer(&sliceHeader))
  copy(s, b)
  return p
}

//export disposePtr
func disposePtr(ptr unsafe.Pointer, ctx unsafe.Pointer) {
  delete(ptrTrckr, uintptr(ptr))
  C.free(ptr)
}

//export ArraySize
func ArraySize(array unsafe.Pointer) C.size_t {
  return ptrTrckr[uintptr(array)]
}

//export _TestFunc
 func _TestFunc (foo *C.char) C.int {
  _foo := C.GoString(foo)
  _returned_value := C.int(abstract.TestFunc(_foo))
  return _returned_value
}

//export _TestFunc2
 func _TestFunc2 (foo *C.char) unsafe.Pointer {
  _foo := C.GoString(foo)
  _returned_value := CFloat32(abstract.TestFunc2(_foo))
  return _returned_value
}

func main() {} // Required but ignored