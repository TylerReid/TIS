package main

import (
	"bytes"
	"strconv"
	"strings"
)

type Board struct {
	Nodes [][]*Node
}

func (b *Board) LinkNodes() {
	for row := 0; row < len(b.Nodes); row++ {
		for column := 0; column < len(b.Nodes[row]); column++ {
			if column+1 != len(b.Nodes[row]) {
				LinkLR(b.Nodes[row][column], b.Nodes[row][column+1])
			}
			if row+1 != len(b.Nodes) {
				LinkUD(b.Nodes[row][column], b.Nodes[row+1][column])
			}
		}
	}
}

func (b *Board) Run() {
	for _, row := range b.Nodes {
		for _, n := range row {
			go n.run()
		}
	}
}

func (b *Board) Stop() {
	for _, row := range b.Nodes {
		for _, n := range row {
			go n.stop()
		}
	}
}

func (b *Board) Print() string {
	var sb bytes.Buffer
	for _, row := range b.Nodes {
		for _, n := range row {
			sb.WriteString("Acc:")
			sb.WriteString(strconv.Itoa(n.Acc))
			sb.WriteString("\t")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

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
	shouldStop chan bool
}

//todo check for a better and idomatic way of doing this
var EmptyNode = &Node{}

type Port string //todo don't know if this is needed. Might just change funcs to accept chan int passed in
const (
	Left  = "LEFT"
	Up    = "UP"
	Right = "RIGHT"
	Down  = "DOWN"
	Any   = "ANY"
	Last  = "LAST"
)

func (n *Node) movToPort(from Port, to Port) {
	f := n.portToChan(from)
	t := n.portToChan(to)
	t <- <-f
}

func (n *Node) movToAcc(from Port) {
	c := n.portToChan(from)
	n.Acc = <-c
}

func (n *Node) movFromAcc(to Port) {
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

func (n *Node) readAny() int {
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

func (n *Node) add(s string) {
	i, err := strconv.Atoi(s)
	if err != nil {
		n.addPort(Port(s))
	} else {
		n.addConst(i)
	}
}

func (n *Node) addConst(i int) {
	n.Acc += i
}

func (n *Node) addPort(p Port) {
	i := <-n.portToChan(p)
	n.addConst(i)
}

func (n *Node) sub(s string) {
	i, err := strconv.Atoi(s)
	if err != nil {
		n.subPort(Port(s))
	} else {
		n.subConst(i)
	}
}

func (n *Node) subConst(i int) {
	n.Acc -= i
}

func (n *Node) subPort(p Port) {
	i := <-n.portToChan(p)
	n.subConst(i)
}

func (n *Node) readLast() int {
	return <-n.last
}

func (n *Node) writeLast(i int) {
	n.last <- i
}

func (n *Node) sav() {
	n.Bak = n.Acc
}

func (n *Node) swp() {
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

func (n *Node) CurrentStatement() string {
	if len(n.Statements) > 0 {
		return n.Statements[n.pc]
	}
	return ""
}

func (n *Node) mov(l, r string) {
	lIsAcc := l == "ACC"
	rIsAcc := r == "ACC"
	if lIsAcc && rIsAcc {
		panic("Two Acc in a MOV")
	}
	if lIsAcc {
		n.movFromAcc(Port(r))
		return
	}
	if rIsAcc {
		n.movToAcc(Port(l))
		return
	}
	n.movToPort(Port(l), Port(r))
}

func (n *Node) run() {
	if n == EmptyNode {
		return
	}
	if n.shouldStop != nil {
		n.shouldStop <- true
	} else {
		n.shouldStop = make(chan bool)
	}
	for {
		select {
		case stop := <-n.shouldStop:
			if stop {
				break
			}
		default:
			s := n.nextStatement()
			statement := strings.Split(s, " ")
			instruction := statement[0]

			if instruction == "MOV" {
				n.mov(statement[1], statement[2])
			}
			if instruction == "ADD" {
				n.add(statement[1])
			}
			if instruction == "SWP" {
				n.swp()
			}
			if instruction == "SAV" {
				n.sav()
			}
			if instruction == "SUB" {
				n.sub(statement[1])
			}
		}
	}
}

func (n *Node) stop() {
	if n.shouldStop != nil {
		n.shouldStop <- true
	}
}
