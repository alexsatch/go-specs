package specs

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func nameOfFunc[T any](fn func(t T) bool) string {
	if fn == nil {
		panic("nil func")
	}

	value := reflect.ValueOf(fn)
	valuePtr := value.Pointer()
	funcForPC := runtime.FuncForPC(valuePtr)
	fqnName := funcForPC.Name()

	lastDot := strings.LastIndex(fqnName, ".")
	if lastDot == -1 {
		return "<no-name>"
	}

	if firstChar := fqnName[lastDot+1]; firstChar >= '0' && firstChar <= '9' {
		file, line := funcForPC.FileLine(valuePtr)
		return fmt.Sprintf("<anonymous: %s:%d>", file, line)
	}

	if li := strings.LastIndex(fqnName[lastDot+1:], "-fm"); li > 0 {
		fqnName = fqnName[:lastDot+1+li]
	}

	fullPathToType := fullPathOfType[T]()
	if strings.HasPrefix(fqnName, fullPathToType) {
		return "." + fqnName[len(fullPathToType):]
	}

	return fqnName
}

func fullPathOfType[T any]() string {
	var dummy T
	tType := reflect.TypeOf(dummy)
	return tType.PkgPath() + "." + tType.Name() + "."
}
