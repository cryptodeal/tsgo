package abstract_gen_tsgo

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

//export _testFunc
func _testFunc(foo *C.char) C.int {
	__foo = C.GoString(_foo)
	defer C.free(_foo)
	return testFunc(_foo)
}
