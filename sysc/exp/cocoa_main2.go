package main

/*
extern int c_main(int argc, char *argv[]);
*/
import "C"

func main() {
	C.c_main(0, nil)
}
