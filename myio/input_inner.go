package myio

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
)

// Inputインターフェースの中身
type input struct{ *bufio.Scanner }

func (in *input) Split(split bufio.SplitFunc)                  { in.Scanner.Split(split) }
func (in *input) Discard()                                     { in.Scan() }
func (in *input) Int() int                                     { return in.i() }
func (in *input) Int2() (int, int)                             { return in.i(), in.i() }
func (in *input) Int3() (int, int, int)                        { return in.i(), in.i(), in.i() }
func (in *input) Int4() (int, int, int, int)                   { return in.i(), in.i(), in.i(), in.i() }
func (in *input) Float() float64                               { return in.f() }
func (in *input) Float2() (float64, float64)                   { return in.f(), in.f() }
func (in *input) Float3() (float64, float64, float64)          { return in.f(), in.f(), in.f() }
func (in *input) Float4() (float64, float64, float64, float64) { return in.f(), in.f(), in.f(), in.f() }
func (in *input) String() string                               { return in.s() }
func (in *input) String2() (string, string)                    { return in.s(), in.s() }
func (in *input) String3() (string, string, string)            { return in.s(), in.s(), in.s() }
func (in *input) String4() (string, string, string, string)    { return in.s(), in.s(), in.s(), in.s() }
func (in *input) Runes() []rune                                { return []rune(in.s()) }
func (in *input) Ints(n int) []int {
	res := make([]int, n)
	for i := range res {
		res[i] = in.i()
	}
	return res
}
func (in *input) Floats(n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = in.f()
	}
	return res
}
func (in *input) Strings(n int) []string {
	res := make([]string, n)
	for i := range res {
		res[i] = in.s()
	}
	return res
}

func (in *input) i() int {
	if err := in.checkScan(); err != nil {
		log.Panicln(fmt.Errorf("input int: %w", err))
	}
	res, err := strconv.Atoi(in.Text())
	if err != nil {
		log.Panicln(fmt.Errorf("input int: %w", err))
	}
	return res
}
func (in *input) f() float64 {
	if err := in.checkScan(); err != nil {
		log.Panicln(fmt.Errorf("input float: %w", err))
	}
	res, err := strconv.ParseFloat(in.Text(), 64)
	if err != nil {
		log.Panicln(fmt.Errorf("input float: %w", err))
	}
	return res
}
func (in *input) s() string {
	if err := in.checkScan(); err != nil {
		log.Panicln(fmt.Errorf("input string: %w", err))
	}
	return in.Text()
}
func (in *input) checkScan() error {
	if !in.Scan() {
		if err := in.Err(); err != nil {
			return err
		}
		return io.EOF
	}
	return nil
}
