// go build -ldflags="-linkmode=external"

package pin

// #cgo pkg-config: gtk+-3.0
// #include "application_linux.h"
import "C"

type Runnable func()

//export CallVoidFunc
func CallVoidFunc(ptr interface{}) {
	if runnable, ok := ptr.(Runnable); ok {
		runnable()
	}
}

func ShowGTKWindow() C.int {
	return C.GtkShowWindow()
}

func main() {
	ShowGTKWindow()
}
