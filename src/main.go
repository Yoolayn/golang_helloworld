package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
)

type hello struct {
	hello string
	world string
}

func (h *hello) Hello() (ok bool, err error) {
	ok = true
	err = nil
	if h.hello != "" {
		ok = false
		err = errors.New("hello exists")
		return
	}
	(*h).hello = "hello"
	return
}

func (h *hello) World() (ok bool, err error) {
	ok = true
	err = nil
	if h.world != "" {
		ok = false
		err = errors.New("world exists")
		return
	}
	(*h).world = "world"
	return
}

func formatter (h1, h2, w1, w2 string, ch chan<- string, wg *sync.WaitGroup) {
	hello_str := fmt.Sprintf(
		"%s%s, %s%s",
		strings.ToUpper(h1),
		h2,
		strings.ToUpper(w1),
		w2,
	)
	ch <- hello_str
	wg.Done()
}

func (h *hello) get_hello() string {
	return h.hello
}

func (h *hello) get_world() string {
	return h.world
}

func printer(wg *sync.WaitGroup, ch <-chan string) {
	str := <- ch
	fmt.Println(str)
	wg.Done()
}

func exclamation(rec <- chan string, send chan <- string) {
	str := <- rec
	send <- fmt.Sprintf("%s!", str)
}

func main() {
	var wg sync.WaitGroup
	hello_world := hello{}
	ok, err := hello_world.Hello()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("not ok")
	}

	ok, err = hello_world.World()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("not ok")
	}

	wg.Add(2)
	noex := make(chan string)
	ex := make(chan string)
	go printer(&wg, ex)
	go formatter(
		hello_world.get_hello()[0:1],
		hello_world.get_hello()[1:],
		hello_world.get_world()[0:1],
		hello_world.get_world()[1:],
		noex,
		&wg,
	)
	go exclamation(noex, ex)
	wg.Wait()
	close(noex)
}
