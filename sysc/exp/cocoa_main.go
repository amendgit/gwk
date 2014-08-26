package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

extern int NSApplicationMain(int, const char *[]);
extern int c_main(int argc, char *argv[]);
*/
import "C"

func main() {
	// C.NSApplicationMain(0, nil)
	C.c_main(0, nil)
}
