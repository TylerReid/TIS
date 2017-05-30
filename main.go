package main

import "fmt"
import "time"

func main() {
	in := make(chan int)
	out := make(chan int)

	aStatements := make([]string, 0)
	aStatements = append(aStatements, "MOV LEFT RIGHT")
	a := Node{Acc: 1, Left: in, Statements: aStatements}

	bStatements := make([]string, 0)
	bStatements = append(bStatements, "MOV LEFT RIGHT")
	b := Node{Acc: 2, Right: out, Statements: bStatements}

	LinkLR(&a, &b)

	go a.Run()
	go b.Run()

	go func() {
		for n := 0; ; n++ {
			in <- n
		}
	}()

	for {
		time.Sleep(time.Second)
		fmt.Println(<-out)
	}
}
