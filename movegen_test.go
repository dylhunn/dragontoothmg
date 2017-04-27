package main

import (
	"math/bits"
	"testing"
)

func TestWhitePawnPush(t *testing.T) {
	var whitePawnsBefore uint64 = 0xFF00 // white on second rank
	var whitePawnsAfter uint64 = 0xFCFD0000
	var blackPawns uint64 = 0x1020000 // black on 24 and 17
	whitepieces := pieces{pawns: whitePawnsBefore, all: whitePawnsBefore}
	blackpieces := pieces{pawns: blackPawns, all: blackPawns}
	testboard := Board{white: whitepieces, black: blackpieces, wtomove: true}
	moves := make([]Move, 0)
	testboard.whitePawnPushes(&moves)
	for _, v := range moves {
		if ((1 << v.to) & whitePawnsAfter) == 0 {
			t.Error("Generated move was not expected:", v)
		}
		whitePawnsAfter -= 1 << v.to
	}
	if whitePawnsAfter != 0 {
		t.Error("An expected move was not found to square", bits.TrailingZeros64(whitePawnsAfter))
	}
}

func TestWhitePawnCapture(t *testing.T) {
	// Board setup:
	// 56  57  58  59  60  BN  62  63
	// 48  49  50  51  52  53  WW  55
	// 40  41  42  43  44  45  46  47
	// 32  33  BB  WW  36  37  38  39
	// 24  BB  BB  27  28  29  30  31
	// 16  WW  18  19  20  21  22  23
	// 8   9   WW  11  WW  13  14  15
	// 0   1   2   3   4   5   6   7
	// white: 0000000001000000000000000000100000000000000000100001010000000000
	// black: 0000000000000000000000000000010000000110000000000000000000000000
	var whitePawns uint64 = 0x40000800021400 // white on 10, 12, 17, 35, 54
	var blackPawns uint64 = 0x406000000      // black on 25, 26, 34
	var blackKnight uint64 = 1 << 61         // black on 61 (for capture promotion)
	// en passant target is 42
	whitepieces := pieces{pawns: whitePawns, all: whitePawns}
	blackpieces := pieces{pawns: blackPawns, knights: blackKnight, all: blackPawns | blackKnight}
	testboard := Board{white: whitepieces, black: blackpieces, wtomove: true, enpassant: 42}
	moves := make([]Move, 0)
	testboard.whitePawnCaptures(&moves)
	if len(moves) != 6 {
		t.Error("Pawn capture moves: wrong length. Expected 6, got", len(moves))
	}
}
