package main

import (
	"math/bits"
	"testing"
)

func TestWhitePawnPush(t *testing.T) {
	var whitePawnsBefore uint64 = 0xFF00 // white on second rank
	var whitePawnsAfter uint64 = 0xFCFD0000
	var blackPawns uint64 = 0x1020000 // black on 24 and 17
	whitepieces := Bitboards{pawns: whitePawnsBefore, all: whitePawnsBefore}
	blackpieces := Bitboards{pawns: blackPawns, all: blackPawns}
	testboard := Board{white: whitepieces, black: blackpieces, wtomove: true}
	moves := make([]Move, 0, 45)
	testboard.whitePawnPushes(&moves)
	for _, v := range moves {
		if ((1 << v.To()) & whitePawnsAfter) == 0 {
			t.Error("Generated move was not expected:", v)
		}
		whitePawnsAfter -= 1 << v.To()
	}
	if whitePawnsAfter != 0 {
		t.Error("An expected move was not found to square", bits.TrailingZeros64(whitePawnsAfter))
	}
	if len(moves) != 13 {
		t.Error("Unexpected number of moves")
	}
}

func TestPawnCapturePosition0(t *testing.T) {
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
	// black pawns: 0000000000000000000000000000010000000110000000000000000000000000
	var whitePawns uint64 = 0x40000800021400 // white on 10, 12, 17, 35, 54
	var blackPawns uint64 = 0x406000000      // black on 25, 26, 34
	var blackKnight uint64 = 1 << 61         // black on 61 (for capture promotion)
	// en passant target is 42
	whitepieces := Bitboards{pawns: whitePawns, all: whitePawns}
	blackpieces := Bitboards{pawns: blackPawns, knights: blackKnight, all: blackPawns | blackKnight}
	testboard := Board{white: whitepieces, black: blackpieces, wtomove: true, enpassant: 42}
	moves := make([]Move, 0, 45)
	testboard.pawnCaptures(&moves)
	if len(moves) != 6 {
		t.Error("Pawn capture moves: wrong length. Expected 6, got", len(moves))
	}

	testboard.wtomove = false
	testboard.enpassant = 0
	moves2 := make([]Move, 0, 45)
	testboard.pawnCaptures(&moves2)
	if len(moves2) != 1 {
		t.Error("Pawn capture moves: wrong length. Expected 1, got", len(moves2))
	}
}

func TestPawnCapturePosition1(t *testing.T) {
	// Board setup:
	// 56  57  58  59  60  61  62  63
	// 48  49  50  51  52  53  BB  55
	// 40  41  42  43  44  45  46  47
	// 32  33  34  35  BB  37  38  39
	// 24  25  BB  WW  28  BB  BB  31
	// 16  17  18  19  20  WW  22  23
	// BB  WW  WW  11  WW  13  WW  WW
	// 0   WN  2   3   4   5   6   7
	// white pawns: 0000000000000000000000000000000000001000001000001101011000000000
	// black: 0000000001000000000000000001000001100100000000000000000100000000
	var whitePawns uint64 = 0x820D600
	var blackPawns uint64 = 0x40001064000100
	var whiteKnight uint64 = 1 << 1 // white on 1 (for capture promotion)
	// en passant target is 19
	whitepieces := Bitboards{pawns: whitePawns, knights: whiteKnight, all: whitePawns | whiteKnight}
	blackpieces := Bitboards{pawns: blackPawns, all: blackPawns}
	testboard := Board{white: whitepieces, black: blackpieces, wtomove: false, enpassant: 19}
	moves := make([]Move, 0, 45)
	testboard.pawnCaptures(&moves)
	if len(moves) != 7 {
		t.Error("Pawn capture moves: wrong length. Expected 7, got", len(moves))
	}

	testboard.wtomove = true
	testboard.enpassant = 0
	moves2 := make([]Move, 0, 45)
	testboard.pawnCaptures(&moves2)
	if len(moves2) != 2 {
		t.Error("Pawn capture moves: wrong length. Expected 2, got", len(moves2))
	}
}
