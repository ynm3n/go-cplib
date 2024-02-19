package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	bufferedOutput      = true
	initialInputBufSize = 1 << 15
)

// 解答欄
func solve(in Input, out Output) {

}

// 入出力のバッファリング設定とsolve関数の呼び出し
func Solve(r io.Reader, w io.Writer) {
	in := NewInput(r, initialInputBufSize)
	var out Output
	if bufferedOutput {
		bw := bufio.NewWriter(w)
		defer bw.Flush()
		out = NewOutput(bw, "", 0)
	} else {
		out = NewOutput(w, "", 0)
	}
	for i := 0; i < 1; i++ {
		solve(in, out)
	}
}

// logパッケージの設定とSolve関数の呼び出し
func main() {
	if os.Getenv("ATCODER") == "1" {
		log.SetOutput(io.Discard)
	} else {
		const errHeader = "\n==========Stderr==========\n"
		buf := bytes.NewBufferString(errHeader)
		log.SetOutput(buf)
		log.SetFlags(log.Lshortfile)
		defer func() {
			if buf.Len() > len(errHeader) {
				io.Copy(os.Stderr, io.LimitReader(buf, 1<<10))
				if buf.Len() > 0 {
					fmt.Fprintln(os.Stderr, " …and more")
				}
			}
		}()
	}
	Solve(os.Stdin, os.Stdout)
}

type Output interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
}

// default: not buffered
// バッファリングしたい場合は *bufio.Writer を渡す
func NewOutput(w io.Writer, prefix string, flag int) Output {
	return log.New(w, prefix, flag)
}

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
