package main

import (
	"fmt"
	"math/bits"
)

type Board struct {
	// is it white's turn to move?
	wtomove bool

	// the square id (16-23 or 40-47) on which en passant capture is possible
	// a value of 0 is no en passant
	enpassant uint8

	// the pieces for each player
	white pieces
	black pieces
}

// the piece bitboards for a player
type pieces struct {
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
	pawns   uint64
	bishops uint64
	knights uint64
	rooks   uint64
	queens  uint64
	kings   uint64
	all     uint64
}

type Move struct {
	from    uint8
	to      uint8
	promote uint8 // piece type
}

const ( // piece types
	nothing = iota
	pawn    = iota
	knight  = iota // list before bishop for promotion loops
	bishop  = iota
	rook    = iota
	queen   = iota
	king    = iota
)

func (b *Board) whitePawnPushes(moveList *[]Move) {
	targets, doubleTargets := b.whitePawnPushBitboards()
	// push all pawns by one square
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1 // unset the lowest active bit
		move := Move{from: uint8(target - 8), to: uint8(target)}
		if target >= 56 { // promotion
			for i := uint8(knight); i <= queen; i++ {
				move.promote = i
				*moveList = append(*moveList, move)
			}
		} else {
			*moveList = append(*moveList, move)
		}
	}
	// push some pawns by two squares
	for doubleTargets != 0 {
		doubleTarget := bits.TrailingZeros64(doubleTargets)
		doubleTargets &= doubleTargets - 1 // unset the lowest active bit
		move := Move{from: uint8(doubleTarget - 16), to: uint8(doubleTarget)}
		*moveList = append(*moveList, move)
	}
}

func (b *Board) whitePawnPushBitboards() (targets uint64, doubleTargets uint64) {
	free := (^b.white.all) & (^b.black.all)
	targets = b.white.pawns << 8 & free
	fourthFile := uint64(0xFF000000)
	doubleTargets = targets << 8 & fourthFile & free
	return
}

func (b *Board) whitePawnCaptures(moveList *[]Move) {
	east, west := b.whitePawnCaptureBitboards()
	for east != 0 {
		target := bits.TrailingZeros64(east)
		east &= east - 1
		move := Move{from: uint8(target - 9), to: uint8(target)}
		if target >= 56 { // promotion
			for i := uint8(knight); i <= queen; i++ {
				move.promote = i
				*moveList = append(*moveList, move)
			}
		} else {
			*moveList = append(*moveList, move)
		}
	}
	for west != 0 {
		target := bits.TrailingZeros64(west)
		west &= west - 1
		move := Move{from: uint8(target - 7), to: uint8(target)}
		if target >= 56 { // promotion
			for i := uint8(knight); i <= queen; i++ {
				move.promote = i
				*moveList = append(*moveList, move)
			}
		} else {
			*moveList = append(*moveList, move)
		}
	}

}

func (b *Board) whitePawnCaptureBitboards() (east uint64, west uint64) {
	notAFile := uint64(0x7F7F7F7F7F7F7F7F)
	notHFile := uint64(0xFEFEFEFEFEFEFEFE)
	blacktargets := b.black.all
	if b.enpassant >= 40 { // a black en-passant target exists
		blacktargets |= (1 << b.enpassant)
	}
	east = b.white.pawns << 9 & notAFile & blacktargets
	west = b.white.pawns << 7 & notHFile & blacktargets
	return
}

func main() {
	var test uint16 = 1
	var test2 uint64 = 0

	fmt.Printf("%v d%v %v %v\n", test, test<<8, ^test, bits.TrailingZeros64(test2))

}
