package errs

import (
	"path"
	"runtime"
	"strconv"
	"strings"
)

type stack []uintptr

func (s *stack) String() string {
	if len(*s) == 0 || !isTrace {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\nStack Trace:")
	for _, pc := range *s {
		fn := runtime.FuncForPC(pc - 1)
		if fn == nil {
			continue
		}
		name := path.Base(fn.Name())
		if strings.HasPrefix(name, "runtime.") {
			continue
		}
		if strings.HasPrefix(name, "testing.") {
			continue
		}
		file, line := fn.FileLine(pc)
		sb.WriteString("\n\t")
		sb.WriteString(name)
		sb.WriteString("() ")
		sb.WriteString(file)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(line))
	}
	return sb.String()
}

func callers(skip int) *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0:n]
	return &st
}
