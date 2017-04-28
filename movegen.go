package main

import (
	"fmt"
	"math/bits"
)


func (b *Board) whitePawnPushes(moveList *[]Move) {
	targets, doubleTargets := b.whitePawnPushBitboards()
	// push all pawns by one square
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1 // unset the lowest active bit
		var move Move
		move.setfrom(Square(target - 8)).setto(Square(target))
		if target >= 56 { // promotion
			for i := Piece(knight); i <= queen; i++ {
				move.setpromote(i)
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
		var move Move
		move.setfrom(Square(doubleTarget - 16)).setto(Square(doubleTarget))
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
		var move Move
		move.setfrom(Square(target - 9)).setto(Square(target))
		if target >= 56 { // promotion
			for i := Piece(knight); i <= queen; i++ {
				move.setpromote(i)
				*moveList = append(*moveList, move)
			}
		} else {
			*moveList = append(*moveList, move)
		}
	}
	for west != 0 {
		target := bits.TrailingZeros64(west)
		west &= west - 1
		var move Move
		move.setfrom(Square(target - 7)).setto(Square(target))
		if target >= 56 { // promotion
			for i := Piece(knight); i <= queen; i++ {
				move.setpromote(i)
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
