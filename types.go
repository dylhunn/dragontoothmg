package main

import (
	"fmt"
)

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
// MSB---------------------------------------------------LSB
// H8 G8 F8 E8 D8 C8 B8 A8 H7 ... A2 H1 G1 F1 E1 D1 C1 B1 A1

type Board struct {
	wtomove   bool
	enpassant uint8 // square id (16-23 or 40-47) where en passant capture is possible
	white     Bitboards
	black     Bitboards
}

type Bitboards struct {
	pawns   uint64
	bishops uint64
	knights uint64
	rooks   uint64
	queens  uint64
	kings   uint64
	all     uint64
}

// Move bitwise structure
// Data stored inside, from LSB
// 6 bits: destination square
// 6 bits: source square
// 3 bits: promotion
type Move uint32

func (m *Move) To() Square {
	return Square(*m & 0x3F)
}
func (m *Move) From() Square {
	return Square((*m & 0xFC0) >> 6)
}
func (m *Move) Promote() Piece {
	return Piece((*m & 0x7000) >> 12)
}
func (m *Move) Setto(s Square) *Move {
	*m = *m & ^(Move(0x3F)) | Move(s)
	return m
}
func (m *Move) Setfrom(s Square) *Move {
	*m = *m & ^(Move(0xFC0)) | (Move(s) << 6)
	return m
}
func (m *Move) Setpromote(p Piece) *Move {
	*m = *m & ^(Move(0x7000)) | (Move(p) << 12)
	return m
}
func (m *Move) String() string {
	return fmt.Sprintf("[from: %v, to: %v, promote: %v]", m.From(), m.To(), m.Promote())
}

// Square index values from 0-63
type Square uint8

// Piece types; valid in range 0-6
type Piece uint8

const (
	nothing = iota
	pawn    = iota
	knight  = iota // list before bishop for promotion loops
	bishop  = iota
	rook    = iota
	queen   = iota
	king    = iota
)

// Masks for knight attacks
// In order: knight on A1, B1, C1, ... F8, G8, H8
var knightMasks = [64]uint64{0x20400, 0x50800, 0xa1100, 0x142200,
	0x284400, 0x508800, 0xa01000, 0x402000, 0x2040004, 0x5080008, 0xa110011,
	0x14220022, 0x28440044, 0x50880088, 0xa0100010, 0x40200020, 0x204000402,
	0x508000805, 0xa1100110a, 0x1422002214, 0x2844004428, 0x5088008850,
	0xa0100010a0, 0x4020002040, 0x20400040200, 0x50800080500, 0xa1100110a00,
	0x142200221400, 0x284400442800, 0x508800885000, 0xa0100010a000,
	0x402000204000, 0x2040004020000, 0x5080008050000, 0xa1100110a0000,
	0x14220022140000, 0x28440044280000, 0x50880088500000, 0xa0100010a00000,
	0x40200020400000, 0x204000402000000, 0x508000805000000, 0xa1100110a000000,
	0x1422002214000000, 0x2844004428000000, 0x5088008850000000,
	0xa0100010a0000000, 0x4020002040000000, 0x400040200000000,
	0x800080500000000, 0x1100110a00000000, 0x2200221400000000,
	0x4400442800000000, 0x8800885000000000, 0x100010a000000000,
	0x2000204000000000, 0x4020000000000, 0x8050000000000, 0x110a0000000000,
	0x22140000000000, 0x44280000000000, 0x88500000000000, 0x10a00000000000,
	0x20400000000000}
