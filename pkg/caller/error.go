// Package caller contains functions for getting file name and line number of function call
package caller

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
)

// CodeLine returns file name and line number of function call
func CodeLine() string {
	_, file, lineNo, ok := runtime.Caller(1)
	if !ok {
		return "runtime.CodeLine() failed"
	}

	fileName := path.Base(file)
	dir := filepath.Base(filepath.Dir(file))
	return fmt.Sprintf("%s/%s:%d", dir, fileName, lineNo)
}
