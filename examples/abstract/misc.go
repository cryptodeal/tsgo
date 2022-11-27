// Package level
// Second line of package level comment.
package abstract

// DROPPED: Floating comment at the top

// Comment belonging to Foo
type Foo string
type FooInt64 int64

// Comment for the const group declaration
const (
	ConstNumberValue = 123 // Line comment behind field with value 123
	// Individual comment for field ConstStringValue
	ConstStringValue     = "abc"
	ConstFooValue    Foo = "foo_const_value"
) // DROPPED: Line comment after grouped const

const Alice = "Alice"

/*
 DROPPED: Floating multiline comment somewhere in the middle
 Line two
*/

/*
Multiline comment for StructBar
Some more text
*/
type StructBar struct {
	// Comment for field Field of type Foo
	Field                 Foo   `json:"field"` // Line Comment for field Field of type Foo
	FieldWithWeirdJSONTag int64 `json:"weird"`

	FieldThatShouldBeOptional    *string `json:"field_that_should_be_optional"`
	FieldThatShouldNotBeOptional *string `json:"field_that_should_not_be_optional" tstype:",required"`
	FieldThatShouldBeReadonly    string  `json:"field_that_should_be_readonly" tstype:",readonly"`
	ArrayField                   []float32
	StructField                  *DemoStruct
}

/*
Another example multiline comment
for DemoStruct
*/
type DemoStruct struct {
	ArrayField           *[]float32
	FieldToAnotherStruct *DemoStruct2
}

type DemoStruct2 struct {
	AnotherArray        *[]float64
	BacktoAnotherStruct *DemoStruct3
}

type DemoStruct3 struct {
	AnotherArray *[]float32
}

// DROPPED: Floating comment at the end

func TestStruct() StructBar {
	str := "bar"
	structBar := StructBar{
		Field:                        "foo",
		FieldWithWeirdJSONTag:        123,
		FieldThatShouldNotBeOptional: &str,
		FieldThatShouldBeReadonly:    "readonly",
		ArrayField:                   []float32{1.1, 2.2, 3.3},
		StructField: &DemoStruct{
			ArrayField: &[]float32{1, 2, 3},
			FieldToAnotherStruct: &DemoStruct2{
				AnotherArray: &[]float64{1.1, 2.2, 3.3},
				BacktoAnotherStruct: &DemoStruct3{
					AnotherArray: &[]float32{1, 2, 3},
				}},
		},
	}
	return structBar
}

func TestStruct2() *StructBar {
	str := "bar"
	structBar := &StructBar{
		Field:                        "foo",
		FieldWithWeirdJSONTag:        123,
		FieldThatShouldNotBeOptional: &str,
		FieldThatShouldBeReadonly:    "readonly",
	}
	return structBar
}

func TestMap() *map[int]string {
	m := map[int]string{
		1: "foo",
		2: "bar",
		3: "baz",
		4: "qux",
	}
	return &m
}
