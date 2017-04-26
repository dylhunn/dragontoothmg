package movegen

import "fmt"

type Board struct {
	wtomove  bool
	wpawns   uint64 // white pieces
	wbishops uint64
	wknights uint64
	wrooks   uint64
	wqueens  uint64
	wking    uint64
	bpawns   uint64 // black pieces
	bbishops uint64
	bknights uint64
	brooks   uint64
	bqueens  uint64
	bking    uint64
}

func main() {
	fmt.Printf("hello, world\n")
}
