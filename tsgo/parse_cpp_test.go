package tsgo

import (
	"testing"
)

func TestHeaderParser(t *testing.T) {
	got := parseHeader("test_files/TensorBase.h")
	for k, v := range got {
		if k == "" {
			t.Errorf("got %q, parsed method requires `Identifier`", k)
		}
		t.Logf("Method: %s", k)
		for i, overload := range v.Overloads {
			t.Logf("Overload %d -", i+1)
			for j, arg := range *overload {
				if arg.Identifier == nil {
					t.Errorf("arg %v, missing required field `Identifier`", arg)
				}
				t.Logf("Arg %d: %v", j+1, arg.Identifier)
				if arg.IsPtr() && arg.RefDecl == nil {
					t.Errorf("%q, must be false if %q == nil", "arg.IsPtr()", "arg.RefDecl")
				}
			}
		}

	}
}
