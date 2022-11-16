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
	fmt.Println("foo:", foo)
	return rand.Int()
}
