package main

import "fmt"
import "time"

func main() {
	in := make(chan int)
	out := make(chan int)

	aStatements := make([]string, 0)
	aStatements = append(aStatements, "ADD LEFT")
	aStatements = append(aStatements, "MOV ACC DOWN")
	aStatements = append(aStatements, "ADD LEFT")
	aStatements = append(aStatements, "ADD DOWN")
	aStatements = append(aStatements, "MOV ACC RIGHT")
	a := Node{Left: in, Statements: aStatements}

	bStatements := make([]string, 0)
	bStatements = append(bStatements, "MOV LEFT RIGHT")
	b := Node{Right: out, Statements: bStatements}

	cStatements := make([]string, 0)
	cStatements = append(cStatements, "MOV UP ACC")
	cStatements = append(cStatements, "MOV ACC UP")
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
