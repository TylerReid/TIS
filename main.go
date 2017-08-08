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

	nodes := []*Node{&a, &b, &c}
	board := Board{Nodes: nodes}

	LinkLR(&a, &b)
	LinkUD(&a, &c)

	board.Run()

	go func() {
		for n := 1; n < 20; n++ {
			in <- n
		}
		board.Stop()
	}()

	for {
		time.Sleep(time.Second)
		select {
		case o := <-out:
			fmt.Println(o)
		default:
		}
	}
}
