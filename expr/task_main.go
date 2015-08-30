package main

import (
	"log"
)

type Task func()

type HostWindow struct {
}

func (h *HostWindow) CreateWindow() {
	log.Printf("%p", h)
}

func PostTask(task Task) {
	log.Printf("%p", task)
	task()
}

func main() {
	t0 := func() {
		log.Printf("t0")
	}
	PostTask(t0)
	PostTask(func() {
		log.Printf("t1")
	})
	var h HostWindow
	h.CreateWindow()
}
