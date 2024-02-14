package main

import (
	"io"
	"testing"
)

func genRandomCase(t *testing.T, w io.WriteCloser) {
	t.Helper()
	defer w.Close()
	// ここでテストケースを作る
}

func TestSolve(t *testing.T) {
	for range 10000 {
		pr, pw := io.Pipe()
		go genRandomCase(t, pw)
		w := io.Discard // solveから結果を取得したい場合は変更する
		Solve(pr, w)
	}
}
