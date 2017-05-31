package main

import (
	"strconv"
	"strings"
)

type Node struct {
	Left       chan int
	Up         chan int
	Right      chan int
	Down       chan int
	Acc        int
	Bak        int
	last       chan int
	Statements []string
	pc         int
}

type Port string //todo don't know if this is needed. Might just change funcs to accept chan int passed in
const (
	Left  = "LEFT"
	Up    = "UP"
	Right = "RIGHT"
	Down  = "DOWN"
	Any   = "ANY"
	Last  = "LAST"
)

func (n *Node) MovToPort(from Port, to Port) {
	f := n.portToChan(from)
	t := n.portToChan(to)
	t <- <-f
}

func (n *Node) MovToAcc(from Port) {
	c := n.portToChan(from)
	n.Acc = <-c
}

func (n *Node) MovFromAcc(to Port) {
	c := n.portToChan(to)
	c <- n.Acc
}

func (n *Node) portToChan(p Port) chan int {
	switch p {
	case Left:
		return n.Left
	case Up:
		return n.Up
	case Right:
		return n.Right
	case Down:
		return n.Down
	}
	//todo
	panic(nil)
}

func (n *Node) ReadAny() int {
	select {
	case i := <-n.Left:
		n.last = n.Left
		return i
	case i := <-n.Up:
		n.last = n.Up
		return i
	case i := <-n.Right:
		n.last = n.Right
		return i
	case i := <-n.Down:
		n.last = n.Down
		return i
	}
}

func (n *Node) Add(i int) {
	n.Acc += i
}

func (n *Node) AddPort(p Port) {
	i := <-n.portToChan(p)
	n.Add(i)
}

func (n *Node) Sub(i int) {
	n.Acc -= i
}

func (n *Node) ReadLast() int {
	return <-n.last
}

func (n *Node) WriteLast(i int) {
	n.last <- i
}

func (n *Node) Sav() {
	n.Bak = n.Acc
}

func (n *Node) Swp() {
	n.Acc = n.Acc ^ n.Bak
	n.Bak = n.Acc ^ n.Bak
	n.Acc = n.Acc ^ n.Bak
}

func LinkLR(l, r *Node) {
	c := make(chan int)
	l.Right = c
	r.Left = c
}

func LinkUD(u, d *Node) {
	c := make(chan int)
	u.Down = c
	d.Up = c
}

func (n *Node) nextStatement() string {
	s := n.Statements[n.pc]
	n.pc++
	if n.pc == len(n.Statements) {
		n.pc = 0
	}
	return s
}

func (n *Node) mov(l, r string) {
	lIsAcc := l == "ACC"
	rIsAcc := r == "ACC"
	if lIsAcc && rIsAcc {
		panic("Two Acc in a MOV")
	}
	if lIsAcc {
		n.MovFromAcc(Port(r))
		return
	}
	if rIsAcc {
		n.MovToAcc(Port(l))
		return
	}
	n.MovToPort(Port(l), Port(r))
}

func (n *Node) Run() {
	for {
		s := n.nextStatement()
		statement := strings.Split(s, " ")
		if statement[0] == "MOV" {
			n.mov(statement[1], statement[2])
		}
		if statement[0] == "ADD" {
			i, err := strconv.Atoi(statement[1])
			if err != nil {
				n.AddPort(Port(statement[1]))
			} else {
				n.Add(i)
			}
		}
		if statement[0] == "SWP" {
			n.Swp()
		}
		if statement[0] == "SAV" {
			n.Sav()
		}
		if statement[0] == "SUB" {
			i, err := strconv.Atoi(statement[1])
			if err == nil {
				panic("SUB int can't be parsed")
			}
			n.Sub(i)
		}
	}
}
