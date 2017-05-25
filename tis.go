package main

type Node struct {
	Left  chan int
	Up    chan int
	Right chan int
	Down  chan int
	Acc   int
	Bak   int
	last  chan int
}

type Port int //todo don't know if this is needed. Might just change funcs to accept chan int passed in

const (
	Left Port = iota
	Up
	Right
	Down
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
