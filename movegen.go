package main

import "fmt"

type Board struct {
	wtomove bool
	wp      pieces
	bp      pieces
}

// Each bitboard shall use little-endian rank-file mapping:
// 56  57  58  59  60  61  62  63
// 48  49  50  51  52  53  54  55
// 40  41  42  43  44  45  46  47
// 32  33  34  35  36  37  38  39
// 24  25  26  27  28  29  30  31
// 16  17  18  19  20  21  22  23
// 8   9   10  11  12  13  14  15
// 0   1   2   3   4   5   6   7
// The binary bitboard uint64 thus uses this ordering:
// H8 G8 F8 E8 D8 C8 B8 A8 H7 ... A2 H1 G1 F1 E1 D1 C1 B1 A1
type pieces struct {
	p   uint64
	b   uint64
	n   uint64
	r   uint64
	q   uint64
	k   uint64
	all uint64
}

type Move struct {
}

func (b *Board) genPawnPushesW(moveList []Move) {
	targets := b.wp.p << 8 & (^b.wp.all) & (^b.bp.all)
	

}

func main() {
	var test uint16 = 1

	fmt.Printf("%v %v %v\n", test, test << 8, ^test)

}
