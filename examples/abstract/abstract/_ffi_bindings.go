package abstract_Gen_TSGo

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
  "encoding/json"
)
//export _testFunc
 func _testFunc (foo *C.char) C.int {
  return _testFunc(foo);
}

