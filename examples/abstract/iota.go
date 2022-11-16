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

func TestFunc(foo string) int {
	fmt.Println("(logged from Golang) foo:", foo)
	return rand.Int()
}

func TestFunc2(foo string) []float32 {
	out := []float32{}
	for i := 0; i < 10; i++ {
		out = append(out, rand.Float32())
	}
	fmt.Println("(logged from Golang) foo:", foo)
	return out
}
