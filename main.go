package main

import "fmt"
import "time"

func main() {
	in := make(chan int)
	out := make(chan int)

	aStatements := []string{"ADD LEFT", "MOV ACC DOWN", "ADD LEFT", "ADD DOWN", "MOV ACC RIGHT"}
	a := Node{Left: in, Statements: aStatements}

	bStatements := []string{"MOV LEFT RIGHT"}
	b := Node{Right: out, Statements: bStatements}

	cStatements := []string{"MOV UP ACC", "MOV ACC UP"}
	c := Node{Statements: cStatements}

	LinkLR(&a, &b)
	LinkUD(&a, &c)

	go a.Run()
	go b.Run()
	go c.Run()

	go func() {
		for n := 1; ; n++ {
			in <- n
		}
	}()

	for {
		time.Sleep(time.Second)
		fmt.Println(<-out)
	}
}
