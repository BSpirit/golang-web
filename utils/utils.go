package utils

import (
	"fmt"
	"runtime"
)

func Trace(err error) string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	_, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s line %d:\n\t%s", f.Name(), line, err)
}
