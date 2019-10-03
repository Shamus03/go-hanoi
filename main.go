package main

import (
	"flag"

	termbox "github.com/nsf/termbox-go"
)

//go:generate stacker -type int

func main() {
	var size int
	flag.IntVar(&size, "size", 4, "size of the tower")
	flag.Parse()

	h := hanoi{
		numDisks: size,
	}
	var solver hanoiSolver

	reset := func() {
		h.Reset()
		solver.GenerateMoves(size)
	}
	reset()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch {
			case ev.Key == termbox.KeyCtrlC:
				break loop
			case ev.Ch == 'r':
				reset()
			case ev.Ch == 's':
				solver.Next(&h)
			case ev.Ch == '1':
				h.MoveA()
			case ev.Ch == '2':
				h.MoveB()
			case ev.Ch == '3':
				h.MoveC()
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		draw(h)
	}
}

func drawDisk(center, row, size int) {
	for i := 0; i < size; i++ {
		termbox.SetCell(center+i, row, '=', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(center-i, row, '=', termbox.ColorWhite, termbox.ColorDefault)
	}
}

func draw(h hanoi) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	_, height := termbox.Size()

	mod := h.numDisks*2 + 1
	drawRow := func(s intStack, col int) {
		row := height - 2
		s.Walk(func(v int) {
			for i := 0; i < v; i++ {
				drawDisk(col, row, v)
			}
			row--
		})
	}
	drawRow(h.a, 1*mod)
	drawRow(h.b, 2*mod)
	drawRow(h.c, 3*mod)

	for i := 1; i < mod*4; i++ {
		ch := '_'
		if i%mod == 0 {
			ch = '|'
		}
		termbox.SetCell(i, height-1, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	drawDisk(2*mod, 0, h.hand)

	termbox.Flush()
}

type hanoi struct {
	numDisks int
	hand     int
	a        intStack
	b        intStack
	c        intStack
}

func (h *hanoi) Reset() {
	h.hand = 0
	h.a = intStack{}
	h.b = intStack{}
	h.c = intStack{}
	for i := h.numDisks; i > 0; i-- {
		h.a.Push(i)
	}
}

func (h *hanoi) MoveA() {
	h.move(&h.a)
}

func (h *hanoi) MoveB() {
	h.move(&h.b)
}

func (h *hanoi) MoveC() {
	h.move(&h.c)
}

func (h *hanoi) move(s *intStack) {
	if h.hand == 0 {
		v, ok := s.Pop()
		if ok {
			h.hand = v
		}
	} else {
		top, ok := s.Peek()
		if ok && h.hand > top {
			return
		}
		s.Push(h.hand)
		h.hand = 0
	}
}

type hanoiSolver struct {
	i     int
	moves []rune
}

func (s *hanoiSolver) Next(h *hanoi) {
	if s.i >= len(s.moves) {
		return
	}
	switch s.moves[s.i] {
	case 'A':
		h.MoveA()
	case 'B':
		h.MoveB()
	case 'C':
		h.MoveC()
	}
	s.i++
}

func (s *hanoiSolver) GenerateMoves(size int) {
	s.i = 0
	s.moves = generateMoves(size, 'A', 'B', 'C')
}

func generateMoves(size int, from, via, to rune) []rune {
	if size == 0 {
		return nil
	}
	var moves []rune
	moves = append(moves, generateMoves(size-1, from, to, via)...)
	moves = append(moves, from, to)
	moves = append(moves, generateMoves(size-1, via, from, to)...)
	return moves
}
