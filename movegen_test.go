package dragontoothmg

import (
	"math/bits"
	"testing"
)

func TestWhitePawnPush(t *testing.T) {
	var whitePawnsBefore uint64 = 0xFF00 // white on second rank
	var whitePawnsAfter uint64 = 0xFCFD0000
	var blackPawns uint64 = 0x1020000 // black on 24 and 17
	whitepieces := bitboards{pawns: whitePawnsBefore, all: whitePawnsBefore}
	blackpieces := bitboards{pawns: blackPawns, all: blackPawns}
	testboard := Board{white: whitepieces, black: blackpieces, wtomove: true}
	moves := make([]Move, 0, 45)
	testboard.pawnPushes(&moves)
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

func TestPawnPosition0(t *testing.T) {
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
	whitepieces := bitboards{pawns: whitePawns, all: whitePawns}
	blackpieces := bitboards{pawns: blackPawns, knights: blackKnight, all: blackPawns | blackKnight}
	testboard := Board{white: whitepieces, black: blackpieces, wtomove: true, enpassant: 42}

	moves := make([]Move, 0, 45)
	testboard.pawnCaptures(&moves)
	if len(moves) != 6 {
		t.Error("Pawn capture moves: wrong length. Expected 6, got", len(moves))
	}

	movesc := make([]Move, 0, 45)
	testboard.pawnPushes(&movesc)
	if len(movesc) != 8 {
		t.Error("Pawn push moves: wrong length. Expected 8, got", len(movesc))
	}

	testboard.wtomove = false
	testboard.enpassant = 0
	moves2 := make([]Move, 0, 45)
	testboard.pawnCaptures(&moves2)
	if len(moves2) != 1 {
		t.Error("Pawn capture moves: wrong length. Expected 1, got", len(moves2))
	}

	movesc2 := make([]Move, 0, 45)
	testboard.pawnPushes(&movesc2)
	if len(movesc2) != 1 {
		t.Error("Pawn push moves: wrong length. Expected 1, got", len(movesc2))
	}
}

func TestPawnPosition1(t *testing.T) {
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
	whitepieces := bitboards{pawns: whitePawns, knights: whiteKnight, all: whitePawns | whiteKnight}
	blackpieces := bitboards{pawns: blackPawns, all: blackPawns}
	testboard := Board{white: whitepieces, black: blackpieces, wtomove: false, enpassant: 19}

	moves := make([]Move, 0, 45)
	testboard.pawnCaptures(&moves)
	if len(moves) != 7 {
		t.Error("Pawn capture moves: wrong length. Expected 7, got", len(moves))
	}

	movesc := make([]Move, 0, 45)
	testboard.pawnPushes(&movesc)
	if len(movesc) != 9 {
		t.Error("Pawn push moves: wrong length. Expected 9, got", len(movesc))
	}

	testboard.wtomove = true
	testboard.enpassant = 0
	moves2 := make([]Move, 0, 45)
	testboard.pawnCaptures(&moves2)
	if len(moves2) != 2 {
		t.Error("Pawn capture moves: wrong length. Expected 2, got", len(moves2))
	}

	movesc2 := make([]Move, 0, 45)
	testboard.pawnPushes(&movesc2)
	if len(movesc2) != 9 {
		t.Error("Pawn push moves: wrong length. Expected 9, got", len(movesc2))
	}
}

func TestKnightPosition0(t *testing.T) {
	// Board setup:
	// WN  57  WN  59  60  61  WN  63	W: 2, 4, 3
	// 48  49  50  51  52  53  WN  55	W: 4
	// 40  BN  42  BP  44  45  46  47	B: 5
	// 32  33  WN  35  36  BN  38  39	W: 7	B: 7
	// BN  25  26  27  28  29  30  31	B: 3
	// 16  WP  18  BN  20  21  22  23	B: 8
	// 8   9   10  11  12  13  BN  15	B: 4
	// 0   1   2   3   4   5   6   7

	var whitePawns uint64 = 1 << 17
	var blackPawns uint64 = 1 << 43

	// 0100010101000000000000000000010000000000000000000000000000000000
	var whiteKnights uint64 = 0x4540000400000000

	// 0000000000000000000000100010000000000001000010000100000000000000
	var blackKnights uint64 = 0x22001084000

	whitepieces := bitboards{pawns: whitePawns, knights: whiteKnights, all: whitePawns | whiteKnights}
	blackpieces := bitboards{pawns: blackPawns, knights: blackKnights, all: blackPawns | blackKnights}
	testboard := Board{white: whitepieces, black: blackpieces, wtomove: true}

	moves := make([]Move, 0, 45)
	testboard.knightMoves(&moves)
	if len(moves) != 20 {
		t.Error("Knight moves: wrong length. Expected 20, got", len(moves))
	}

	testboard.wtomove = false
	moves2 := make([]Move, 0, 45)
	testboard.knightMoves(&moves2)
	if len(moves2) != 27 {
		t.Error("Knight moves: wrong length. Expected 27, got", len(moves2))
	}
}

func TestKingPositions(t *testing.T) {
	testboard := ParseFen("1Q2rk2/2p2p2/1n4b1/N7/2B1Pp1q/2B4P/1QPP1P2/4K2R b K e3 4 30")
	moves := make([]Move, 0, 45)
	testboard.kingMoves(&moves)
	if len(moves) != 3 {
		t.Error("King moves: wrong length. Expected 3, got", len(moves))
	}
	testboard.wtomove = true
	moves2 := make([]Move, 0, 45)
	testboard.kingMoves(&moves2)
	if len(moves2) != 4 {
		t.Error("King moves: wrong length. Expected 4, got", len(moves2))
	}
}

func TestRookPositions(t *testing.T) {
	positions := map[string]int{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -":  0,
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq -":  0,
		"rnbqkbnr/ppppppp1/8/8/7R/8/1PPPPPPP/RNBQKBNR w KQkq -": 18,
		"rnbqkbnr/ppppppp1/8/8/7R/8/1PPPPPPP/RNBQKBNR b KQkq -": 4,
		"r1N2bnN/3pp1p1/8/2rR4/7R/8/1PP1P1P1/RN5R w KQkq -":     37,
		"r1N2bnN/3pp1p1/8/2rR4/7R/8/1PP1P1P1/RN5R b KQkq -":     18,
		"8/8/8/3r4/8/8/8/8 w KQkq -":                            0,
		"8/8/8/3r4/8/8/8/8 b KQkq -":                            14,
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.rookMoves(&moves)
		if len(moves) != v {
			t.Error("Rook moves: wrong length. Expected", v, "but got", len(moves))
		}
	}
}

func TestBishopPositions(t *testing.T) {
	positions := map[string]int{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -":    0,
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq -":    0,
		"rnbqkb1r/pp2pppp/8/4P3/5bN1/8/PPP2PPP/RNBQKBNR w KQkq -": 8,
		"rnbqkb1r/pp2pppp/8/4P3/5bN1/8/PPP2PPP/RNBQKBNR b KQkq -": 12,
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.bishopMoves(&moves)
		if len(moves) != v {
			t.Error("Bishop moves: wrong length. Expected", v, "but got", len(moves))
		}
	}
}

func TestQueenPositions(t *testing.T) {
	positions := map[string]int{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -":    0,
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq -":    0,
		"rnbqkb1r/pp2pppp/8/4P3/5bN1/8/PPP2PPP/RNBQKBNR w KQkq -": 8,
		"rnbqkb1r/pp2pppp/8/4P3/5bN1/8/PPP2PPP/RNBQKBNR b KQkq -": 12,
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.bishopMoves(&moves)
		if len(moves) != v {
			t.Error("Bishop moves: wrong length. Expected", v, "but got", len(moves))
		}
	}
}
