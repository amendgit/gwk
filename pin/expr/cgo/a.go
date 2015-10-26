package main

// #include "a.h"
import "C"

type A struct {
}

func main() {
	var a = new(A)
	C.SomeFunc(a)
}
