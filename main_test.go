package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
)

var (
	testCount = 10000
)

func genTestCase(tb testing.TB, w io.Writer) {
	tb.Helper()
	// テストケースを作り、wに書き込むプログラムを書く
}

func genCorrect(tb testing.TB, r io.Reader, w io.Writer) {
	tb.Helper()
	// rからテストケースを取得し、正答をwに書き込むプログラムを書く
}

// 出力が正しいかどうか確認するためのテスト
func TestSolve_Correct(t *testing.T) {
	for range testCount {
		sTestcase := new(bytes.Buffer)
		cTestcase := new(bytes.Buffer)
		mw := io.MultiWriter(sTestcase, cTestcase)
		genTestCase(t, mw)

		sAns := new(bytes.Buffer)
		Solve(sTestcase, sAns)
		cAns := new(bytes.Buffer)
		genCorrect(t, cTestcase, cAns)

		if d := gocmp.Diff(sAns.String(), cAns.String()); len(d) > 0 {
			t.Fatal(d)
		}
	}
}

// panicするケースを探すためのテスト
func TestSolve_Panic(t *testing.T) {
	for range testCount {
		buf := new(bytes.Buffer)
		genTestCase(t, buf)
		testcase := buf.String()

		ok := t.Run(testcase, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil {
					t.Fatal(p)
				}
				out(t, testcase)
			}()
			Solve(buf, io.Discard)
		})
		if !ok {
			return
		}
	}
}

func BenchmarkSolve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		genTestCase(b, buf)
		Solve(buf, io.Discard)
	}
}

func out(tb testing.TB, testcase string) {
	tb.Helper()
	out, err := os.Create("out")
	if err != nil {
		tb.Fatalf("out(helper): %v", err)
	}
	defer out.Close()
	out.WriteString(testcase)
}
