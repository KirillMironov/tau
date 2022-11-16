package cmdutil

import (
	"fmt"
	"os"
)

func Exit(code int, a ...any) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(code)
}

func Exitf(code int, format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(code)
}
