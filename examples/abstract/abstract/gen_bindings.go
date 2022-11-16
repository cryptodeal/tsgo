package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
  "github.com/cryptodeal/tsgo/examples/abstract"
)

//export _TestFunc
 func _TestFunc (foo *C.char) C.int {
  _foo := C.GoString(foo)
    _returned_value := C.int(abstract.TestFunc(_foo))
  return _returned_value
}

func main() {} // Required but ignored