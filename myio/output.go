package myio

import (
	"io"
	"log"
)

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
