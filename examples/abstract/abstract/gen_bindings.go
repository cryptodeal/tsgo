package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

//export _testFunc
func _testFunc(foo *C.char) C.int {
	_foo := C.GoString(foo)
	defer C.free(_foo)
	return TestFunc(_foo)
}

func main() {} // Required but ignored
