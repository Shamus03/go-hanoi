package main

import (
	"flag"

	termbox "github.com/nsf/termbox-go"
)

//go:generate stacker -type disk

type disk int

func (d disk) Size() int {
	return int(d)
}

func (d disk) FitsOnTopOf(o disk) bool {
	return d < o
}

func main() {
	size := flag.Int("size", 4, "size of the tower")
	drawSolver := flag.Bool("drawsolver", false, "also show solution")
	flag.Parse()

	h := hanoi{
		numDisks: *size,
	}
	var solver hanoiSolver

	reset := func() {
		h.Reset()
		solver.GenerateMoves(*size)
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
				if h.MoveA() {
					solver.Move('A')
				}
			case ev.Ch == '2':
				if h.MoveB() {
					solver.Move('B')
				}
			case ev.Ch == '3':
				if h.MoveC() {
					solver.Move('C')
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		if *drawSolver {
			draw(h, solver)
		} else {
			draw(h, hanoiSolver{})
		}
	}
}

func drawDisk(center, row int, d disk) {
	for i := 0; i < d.Size(); i++ {
		termbox.SetCell(center+i, row, '=', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(center-i, row, '=', termbox.ColorWhite, termbox.ColorDefault)
	}
}

func draw(h hanoi, s hanoiSolver) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	_, height := termbox.Size()

	mod := h.numDisks*2 + 1
	drawRow := func(s diskStack, col int) {
		row := height - 2
		s.Walk(func(d disk) {
			for i := 0; i < d.Size(); i++ {
				drawDisk(col, row, d)
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

	for i := 0; i < len(s.moves); i++ {
		termbox.SetCell(i, 2, rune(s.moves[i]), termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.Flush()
}

type hanoi struct {
	numDisks int
	hand     disk
	a        diskStack
	b        diskStack
	c        diskStack
}

func (h *hanoi) Reset() {
	h.hand = 0
	h.a = diskStack{}
	h.b = diskStack{}
	h.c = diskStack{}
	for i := h.numDisks; i > 0; i-- {
		h.a.Push(disk(i))
	}
}

func (h *hanoi) MoveA() bool {
	return h.move(&h.a)
}

func (h *hanoi) MoveB() bool {
	return h.move(&h.b)
}

func (h *hanoi) MoveC() bool {
	return h.move(&h.c)
}

func (h *hanoi) move(s *diskStack) bool {
	if h.hand == 0 {
		d, ok := s.Pop()
		if ok {
			h.hand = d
			return true
		}
	} else {
		top, ok := s.Peek()
		if ok && !h.hand.FitsOnTopOf(top) {
			return false
		}
		s.Push(h.hand)
		h.hand = 0
		return true
	}
	return false
}

type hanoiSolver struct {
	moves string
}

func (s *hanoiSolver) Next(h *hanoi) {
	if len(s.moves) == 0 {
		return
	}
	m := rune(s.moves[0])
	switch m {
	case 'A':
		h.MoveA()
	case 'B':
		h.MoveB()
	case 'C':
		h.MoveC()
	}
	s.Move(m)
}

func (s *hanoiSolver) Move(m rune) {
	if len(s.moves) == 0 {
		s.moves = string(m)
		return
	}
	if rune(s.moves[0]) == m {
		s.moves = s.moves[1:]
	} else {
		s.moves = string(m) + s.moves
	}
}

func (s *hanoiSolver) GenerateMoves(size int) {
	s.moves = generateMoves(size, 'A', 'B', 'C')
}

func generateMoves(size int, from, via, to rune) string {
	if size == 0 {
		return ""
	}
	return generateMoves(size-1, from, to, via) +
		string(from) + string(to) +
		generateMoves(size-1, via, from, to)
}
