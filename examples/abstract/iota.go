package abstract

import (
	"math/rand"
)

type MyIotaType int

const (
	Zero MyIotaType = iota
	One
	Two
	_
	Four
	FourString string = "four"
	_
	AlsoFourString
	Five = 5
	FiveAgain

	Sixteen = iota + 6
	Seventeen
)

func IntTest(foo string) int {
	return rand.Int()
}

func Float32ArrayTest(foo string) []float32 {
	out := []float32{}
	for i := 0; i < 100; i++ {
		out = append(out, rand.Float32())
	}
	return out
}

func Float64ArrayTest(foo string) []float64 {
	out := []float64{}
	for i := 0; i < 100; i++ {
		out = append(out, rand.Float64())
	}
	return out
}

func Int32ArrayTest(foo string) []int32 {
	out := []int32{}
	for i := 0; i < 10; i++ {
		out = append(out, rand.Int31())
	}
	return out
}

func Int64ArrayTest(foo string) []int64 {
	out := []int64{}
	for i := 0; i < 10; i++ {
		out = append(out, rand.Int63())
	}
	return out
}

func Uint32ArrayTest(foo string) []uint32 {
	out := []uint32{}
	for i := 0; i < 10; i++ {
		out = append(out, rand.Uint32())
	}
	return out
}

func Uint64ArrayTest(foo string) []uint64 {
	out := []uint64{}
	for i := 0; i < 10; i++ {
		out = append(out, rand.Uint64())
	}
	return out
}

func StringTest() string {
	return "Hello, World!"
}

func Float32ArgTest(foo []float32) []float32 {
	for i := 0; i < len(foo); i++ {
		foo[i] = (foo[i] * float32(2))
	}
	return foo
}

func Float64ArgTest(foo []float64) []float64 {
	for i := 0; i < len(foo); i++ {
		foo[i] = (foo[i] * float64(2))
	}
	return foo
}

func Int32ArgTest(foo []int32) []int32 {
	for i := 0; i < len(foo); i++ {
		foo[i] = (foo[i] * int32(2))
	}
	return foo
}

func Int64ArgTest(foo []int64) []int64 {
	for i := 0; i < len(foo); i++ {
		foo[i] = (foo[i] * int64(2))
	}
	return foo
}

func Uint32ArgTest(foo []uint32) []uint32 {
	for i := 0; i < len(foo); i++ {
		foo[i] = (foo[i] * uint32(2))
	}
	return foo
}

func Uint64ArgTest(foo []uint64) []uint64 {
	for i := 0; i < len(foo); i++ {
		foo[i] = (foo[i] * uint64(2))
	}
	return foo
}
