package abstract_gen_tsgo

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
  "encoding/json"
 = C.GoString(_foo)
defer C.free(_foo)
)

//export _testFunc
 func _testFunc (foo *C.char) C.int {
  return testFunc(_foo);
}

