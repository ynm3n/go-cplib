package myio

import (
	"bufio"
	"io"
	"math"
)

type Input interface {
	Split(split bufio.SplitFunc)
	Discard()

	Int() int
	Int2() (int, int)
	Int3() (int, int, int)
	Int4() (int, int, int, int)
	Float() float64
	Float2() (float64, float64)
	Float3() (float64, float64, float64)
	Float4() (float64, float64, float64, float64)
	String() string
	String2() (string, string)
	String3() (string, string, string)
	String4() (string, string, string, string)
	Runes() []rune

	Ints(n int) []int
	Floats(n int) []float64
	Strings(n int) []string
}

// default splitfunc: bufio.ScanWords
func NewInput(r io.Reader, bufSize int) Input {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, bufSize), math.MaxInt)
	sc.Split(bufio.ScanWords)
	return &input{sc}
}
