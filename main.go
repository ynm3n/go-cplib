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
	startScannerBufSize = 4096
)

// 解答欄
func solve(in Input, out Output) {
}

// 入出力のバッファリング設定とsolve関数の呼び出し
func Solve(r io.Reader, w io.Writer) {
	in := NewInput(r)
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
		flushLog := setLogBuffer()
		defer flushLog()
	}
	Solve(os.Stdin, os.Stdout)
}

type Input interface {
	ScanInt(a ...*int)
	ScanFloat(a ...*float64)
	ScanString(a ...*string)
	ScanIntSlice(s []int)
	ScanFloatSlice(s []float64)
	ScanStringSlice(s []string)
	Split(split bufio.SplitFunc)
	Discard()
}

// default splitfunc: bufio.ScanWords
func NewInput(r io.Reader) Input {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, startScannerBufSize), math.MaxInt)
	sc.Split(bufio.ScanWords)
	return &input{sc}
}

type Output interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
}

func NewOutput(w io.Writer, prefix string, flag int) Output {
	return log.New(w, prefix, flag)
}

type input struct{ *bufio.Scanner }

func (in *input) ScanInt(a ...*int) {
	if err := scan(in, strconv.Atoi, a...); err != nil {
		log.Panic(fmt.Errorf("ScanInt: %w", err))
	}
}

func (in *input) ScanFloat(a ...*float64) {
	if err := scan(in, func(s string) (float64, error) { return strconv.ParseFloat(s, 64) }, a...); err != nil {
		log.Panic(fmt.Errorf("ScanFloat: %w", err))
	}
}

func (in *input) ScanString(a ...*string) {
	if err := scan(in, func(s string) (string, error) { return s, nil }, a...); err != nil {
		log.Panic(fmt.Errorf("ScanString: %w", err))
	}
}

func (in *input) ScanIntSlice(s []int) {
	if err := scanSlice(in, strconv.Atoi, s); err != nil {
		log.Panic(fmt.Errorf("ScanIntSlice: %w", err))
	}
}

func (in *input) ScanFloatSlice(s []float64) {
	if err := scanSlice(in, func(s string) (float64, error) { return strconv.ParseFloat(s, 64) }, s); err != nil {
		log.Panic(fmt.Errorf("ScanFloatSlice: %w", err))
	}
}

func (in *input) ScanStringSlice(s []string) {
	if err := scanSlice(in, func(s string) (string, error) { return s, nil }, s); err != nil {
		log.Panic(fmt.Errorf("ScanStringSlice: %w", err))
	}
}

func (in *input) Split(split bufio.SplitFunc) { in.Scanner.Split(split) }

func (in *input) Discard() { in.Scan() }

func scan[T any](in *input, f func(s string) (T, error), a ...*T) error {
	var err error
	for _, p := range a {
		if !in.Scan() {
			if err = in.Err(); err != nil {
				return err
			}
			return io.EOF
		}
		if *p, err = f(in.Text()); err != nil {
			return err
		}
	}
	return nil
}

func scanSlice[T any](in *input, f func(s string) (T, error), s []T) error {
	var err error
	for i := range s {
		if !in.Scan() {
			if err = in.Err(); err != nil {
				return err
			}
			return io.EOF
		}
		if s[i], err = f(in.Text()); err != nil {
			return err
		}
	}
	return nil
}

func setLogBuffer() func() {
	const errHeader = "\n==========Stderr==========\n"
	buf := bytes.NewBufferString(errHeader)
	log.SetOutput(buf)
	log.SetFlags(log.Lshortfile)
	return func() {
		if buf.Len() > len(errHeader) {
			io.Copy(os.Stderr, io.LimitReader(buf, 1<<10))
			if buf.Len() > 0 {
				fmt.Fprintln(os.Stderr, " …and more")
			}
		}
	}
}
