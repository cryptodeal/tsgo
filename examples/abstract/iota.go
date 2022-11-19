package abstract

import (
	"fmt"
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
	fmt.Println("(logged from Golang) foo:", foo)
	return rand.Int()
}

func Float32ArrayTest(foo string) []float32 {
	out := []float32{}
	for i := 0; i < 100; i++ {
		out = append(out, rand.Float32())
	}
	fmt.Println("(logged from Golang) foo:", foo)
	return out
}

func Float64ArrayTest(foo string) []float64 {
	out := []float64{}
	for i := 0; i < 100; i++ {
		out = append(out, rand.Float64())
	}
	fmt.Println("(logged from Golang) foo:", foo)
	return out
}

func Int32ArrayTest(foo string) []int32 {
	out := []int32{}
	for i := 0; i < 10; i++ {
		out = append(out, rand.Int31())
	}
	fmt.Println("(logged from Golang) foo:", foo)
	return out
}

func Int64ArrayTest(foo string) []int64 {
	out := []int64{}
	for i := 0; i < 10; i++ {
		out = append(out, rand.Int63())
	}
	fmt.Println("(logged from Golang) foo:", foo)
	return out
}

func Uint32ArrayTest(foo string) []uint32 {
	out := []uint32{}
	for i := 0; i < 10; i++ {
		out = append(out, rand.Uint32())
	}
	fmt.Println("(logged from Golang) foo:", foo)
	return out
}

func Uint64ArrayTest(foo string) []uint64 {
	out := []uint64{}
	for i := 0; i < 10; i++ {
		out = append(out, rand.Uint64())
	}
	fmt.Println("(logged from Golang) foo:", foo)
	return out
}

func StringTest() string {
	return "Hello, World!"
}
