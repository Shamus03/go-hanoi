# Hanoi

Simple game of Hanoi in the terminal.

1/2/3 to pick up/put down a disk, R to reset, and ctrl-c to quit.

Hanoi knows how to solve itself. Hitting S will execute the next move in the pre-generated solution. Right now it breaks if you've made any manual moves, and you will have to reset the game to properly execute a solution.

## Installation

```bash
$ go get github.com/shamus03/go-hanoi/cmd/hanoi
$ hanoi
```