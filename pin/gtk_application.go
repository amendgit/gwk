// go build -ldflags="-linkmode=external"

package pin

// #cgo pkg-config: gtk+-3.0
// #include "gtk_application.h"
// #include "events.h"
import "C"

const (
	WindowEvent_Restore = C.WindowEvent_Restore
)

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
