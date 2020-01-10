package main

import "fmt"

type Card interface {
	display()
}

type Memory interface {
	storage()
}

type CPU interface {
	calculate()
}

type Inter struct {
}

func (i Inter) display() {
	fmt.Println("inter display")
}

func (i Inter) storage() {
	fmt.Println("inter storage")
}

func (i Inter) calculate() {
	fmt.Println("inter calculate")
}

type Kingston struct {
	Memory
}

func (k Kingston) storage() {
	fmt.Println("kingston storage")
}

type NVIDIA struct {
	Card
}

func (n NVIDIA) display() {
	fmt.Println("nvidia diaplay")
}

type Computer struct {
	Card
	Memory
	CPU
}

func (c Computer) Run(card Card, memory Memory, cpu CPU) {
	card.display()
	memory.storage()
	cpu.calculate()
}

func main() {
	c := new(Computer)
	interc := new(Inter)
	c.Run(NVIDIA{}, Kingston{}, interc)
}
