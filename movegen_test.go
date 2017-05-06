package dragontoothmg

import (
	"math/bits"
	"testing"
)

func TestPawnPushes(t *testing.T) {
	positions := map[string]int{
		"rnbqkbnr/ppp2pp1/3p4/4p3/3N1P2/P1n5/2PPP3/R1BQKBNR w KQkq - 0 0": 5,
		"rnbqkbnr/ppp2pp1/3p4/4p3/3N1P2/P1n5/2PPP3/R1BQKBNR b KQkq - 0 0": 12,
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.pawnPushes(&moves, everything)
		if len(moves) != v {
			t.Error("Pawn pushes: wrong length. Expected", v, "but got", 
				len(moves), "for FEN", b.ToFen())
		}
	}
}

func TestPawnCaptures(t *testing.T) {
	positions := map[string]int{
		"rnbqkbnr/ppp2pp1/3p4/4p3/3N1P2/P1n5/2PPP3/R1BQKBNR w KQkq - 0 0": 2,
		"rnbqkbnr/ppp2pp1/3p4/4p3/3N1P2/P1n5/2PPP3/R1BQKBNR b KQkq - 0 0": 2,
		"rnbqkbnr/ppp2pp1/3p4/4pP2/3N4/P1n5/2PPP3/R1BQKBNR w KQkq e6 0 0": 2,
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.pawnCaptures(&moves, everything)
		if len(moves) != v {
			t.Error("Pawn captures: wrong length. Expected", v, "but got", 
				len(moves), "for FEN", b.ToFen())
		}
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
	testboard.knightMoves(&moves, everything)
	if len(moves) != 20 {
		t.Error("Knight moves: wrong length. Expected 20, got", len(moves))
	}

	testboard.wtomove = false
	moves2 := make([]Move, 0, 45)
	testboard.knightMoves(&moves2, everything)
	if len(moves2) != 27 {
		t.Error("Knight moves: wrong length. Expected 27, got", len(moves2))
	}
}

func TestKingPositions(t *testing.T) {
	positions := map[string]int{
		"1Q2rk2/2p2p2/1n4b1/N7/2B1Pp1q/2B4P/1QPP1P2/4K2R b K e3 4 30": 2,
		"1Q2rk2/2p2p2/1n4b1/N7/2B1Pp1q/2B4P/1QPP1P2/4K2R w K e3 4 30": 4,
		"r3k2r/7B/8/8/3q4/8/P6P/R3K2R w KQkq - 0 0":                   2,
		"r3k2r/7B/8/8/3q4/8/P6P/R3K2R b KQkq - 0 0":                   6,
		"8/1pk5/8/8/8/2R4b/8/4K2R w K -":                              4,
		"8/1pk5/8/8/7b/2R5/8/4K2R b K - 0 0":                          5,
		"4k3/8/8/8/8/8/8/4K2R w K - 0 0":                              6, // white short castle
		"4k3/8/8/8/8/8/8/4K1NR w K - 0 0":                             5, // short castle blocked
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.kingMoves(&moves)
		if len(moves) != v {
			t.Error("King moves: wrong length. Expected", v, "but got",
				len(moves), "\nFor position:", k)
		}
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
		b.rookMoves(&moves, everything)
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
		b.bishopMoves(&moves, everything)
		if len(moves) != v {
			t.Error("Bishop moves: wrong length. Expected", v, "but got", len(moves))
		}
	}
}

func TestQueenPositions(t *testing.T) {
	positions := map[string]int{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -": 0,
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq -": 0,
		"6nq/6p1/2B4n/1rB2r1R/5q2/2P5/1Q4n1/2B5 w - -":         12,
		"6nq/6p1/2B4n/1rB2r1R/5q2/2P5/1Q4n1/2B5 b - -":         21,
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.queenMoves(&moves, everything)
		if len(moves) != v {
			t.Error("Queen moves: wrong length. Expected", v, "but got", len(moves))
		}
	}
}

func TestUnderDirectAttack(t *testing.T) {
	b1 := ParseFen("r1N1kbnN/3pp1p1/1p2q3/2rR1b2/2QP1nBR/6B1/1PP1P1P1/RNK4R b - - 0 0")
	solutionsByBlack := map[uint8]bool{
		AlgebraicToIndex("a5"): true,
		AlgebraicToIndex("a7"): true,
		AlgebraicToIndex("d8"): true,
		AlgebraicToIndex("f7"): true,
		AlgebraicToIndex("h7"): true,
		AlgebraicToIndex("h6"): true,
		AlgebraicToIndex("d8"): true,
		AlgebraicToIndex("e2"): true,
		AlgebraicToIndex("f5"): true,
		AlgebraicToIndex("b5"): true,
		AlgebraicToIndex("d1"): false,
		AlgebraicToIndex("g5"): false,
		AlgebraicToIndex("b7"): false,
	}
	for k, v := range solutionsByBlack {
		attacked := b1.underDirectAttack(true, k)
		if attacked != v {
			t.Error("Under attack failed for position", b1.ToFen(), "\nat coord ", IndexToAlgebraic(Square(k)))
		}
	}

	b2 := ParseFen("r1N1kbnN/3pp3/1p2q3/2rR1bpP/2QP1nBR/6B1/1PP1P1P1/RNK4R b - g6 0 0")
	solutionsByWhite := map[uint8]bool{
		AlgebraicToIndex("c2"): true, // TODO(dylhunn): this case is dubious
		AlgebraicToIndex("b3"): true,
		AlgebraicToIndex("b5"): true,
		AlgebraicToIndex("b6"): true,
		AlgebraicToIndex("g6"): true,
		AlgebraicToIndex("g8"): false,
		AlgebraicToIndex("e6"): false,
		AlgebraicToIndex("e8"): false,
	}
	for k, v := range solutionsByWhite {
		attacked := b2.underDirectAttack(false, k)
		if attacked != v {
			t.Error("Under attack failed for position", b2.ToFen(), "\nat coord ", IndexToAlgebraic(Square(k)))
		}
	}
}

// Test that the only legal moves are those that break check, through:
// - moving the king
// - capture the checking piece
// - breaking the check
func testBreakCheck(t *testing.T) {
	positions := map[string]int{
		"k1N5/3RrQ2/8/2B4R/8/2N5/8/4K3 w - - 0 0": 13, // Non-pawn check-breaks and captures
		"8/8/1p2p3/R6k/8/8/8/K7 b - - 0 0":        3,  // breaks and captures with a pawn
		"3k4/2P4r/1P6/8/8/8/8/K7 b - - 0 0":       5,  // break the check of a pawn
		"3k4/2P1P3/1P6/8/8/8/8/K7 b - - 0 0":      4,  // double check with pawns: must move king
		"3k4/7r/1P6/8/7B/8/3R4/K7 b - - 0 0":      4,  // double check: must move king
		"8/8/8/1k6/2Pp4/8/8/4K3 b - c3 0 0":       9,  // en passant check evasion
		"8/8/8/1k6/3Pp3/8/8/K4Q2 b - d3 0 0":      6,  //en passant check evasion
	}
	for k, v := range positions {
		b := ParseFen(k)
		moves := b.GenerateLegalMoves()
		if len(moves) != v {
			t.Error("Legal moves breaking check: wrong length. Expected", v, "but got", len(moves), "for position", b.ToFen())
		}
	}
}

// Test that pinned pieces can only move along the pin ray

func TestPinnedBishop(t *testing.T) {
	positions := map[string]int{
		"4k3/3b4/8/8/Q7/8/8/4K3 b - - 0 0":      3,  // pinned bishop
		"4k3/3b4/2b5/8/Q7/8/8/4K3 b - - 0 0":    0, // a "double" pin is not actually a pin
		"4k3/3b1b2/2Q3Q1/8/8/8/8/4K3 b - - 0 0": 2,  // two close pins
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.generatePinnedMoves(&moves)
		if len(moves) != v {
			t.Error("Legal moves for pinned bishops: wrong length. Expected", v, "but got", len(moves), "for position", b.ToFen())
		}
	}
}

func TestPinnedKnight(t *testing.T) {
	positions := map[string]int{
		"4k3/3n1n2/2Q3Q1/8/8/8/8/4K3 b - - 0 0": 0, // two close pins
		"4k3/8/8/8/1q6/2N5/8/4K3 w - - 0 0":     0, // normal pin
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.generatePinnedMoves(&moves)
		if len(moves) != v {
			t.Error("Legal moves for pinned bishops: wrong length. Expected", v, "but got", len(moves), "for position", b.ToFen())
		}
	}
}

func TestPinnedQueen(t *testing.T) {
	positions := map[string]int{
		"4k3/8/8/8/1q6/2Q5/8/4K3 w - - 0 0":     2, // normal pin
		"4k3/8/4r3/4Q3/1q6/2Q5/8/4K3 w - - 0 0": 6,
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		b.generatePinnedMoves(&moves)
		if len(moves) != v {
			t.Error("Legal moves for pinned bishops: wrong length. Expected", v, "but got", len(moves), "for position", b.ToFen())
		}
	}
}

func TestDiagPins(t *testing.T) {
	positions := map[string]int{
		"4k3/3p4/2B1p3/8/1q6/4R3/3P4/4K3 w - - 0 0": 0, // diagonal pawns
		"4k3/3p4/2B1p3/8/1q6/4R3/3P4/4K3 b - - 0 0": 2,
		"4k3/8/8/8/1q6/2Q5/8/4K3 w - - 0 0":         2, // normal queen pin
		"4k3/8/4r3/4Q3/1q6/2Q5/8/4K3 w - - 0 0":     6,
		"4k3/8/2p5/8/B7/6q1/5N2/4K3 w - - 0 0":      0,
		"4k3/8/2p5/8/B7/6q1/5N2/4K3 b - - 0 0":      0,
		"4k3/8/6p1/3b3Q/2P5/1K6/8/8 w - - 0 0":      1,
		"4k3/8/6p1/3b3Q/2P5/1K6/8/8 b - - 0 0":      1,
	}
	pinLocs := map[string]uint8{
		"4k3/3p4/2B1p3/8/1q6/4R3/3P4/4K3 w - - 0 0": AlgebraicToIndex("d2"), // diagonal pawns
		"4k3/3p4/2B1p3/8/1q6/4R3/3P4/4K3 b - - 0 0": AlgebraicToIndex("e6"),
		"4k3/8/8/8/1q6/2Q5/8/4K3 w - - 0 0":         AlgebraicToIndex("c3"),
		"4k3/8/4r3/4Q3/1q6/2Q5/8/4K3 w - - 0 0":     AlgebraicToIndex("c3"), // TODO: only checks one of two pins
		"4k3/8/2p5/8/B7/6q1/5N2/4K3 w - - 0 0":      AlgebraicToIndex("f2"),
		"4k3/8/2p5/8/B7/6q1/5N2/4K3 b - - 0 0":      AlgebraicToIndex("c6"),
		"4k3/8/6p1/3b3Q/2P5/1K6/8/8 w - - 0 0":      AlgebraicToIndex("c4"),
		"4k3/8/6p1/3b3Q/2P5/1K6/8/8 b - - 0 0":      AlgebraicToIndex("g6"),
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		result := b.generatePinnedMoves(&moves)
		if len(moves) != v {
			t.Error("Legal moves for diagonal pins: wrong length. Expected", v, "but got", len(moves), "for position", b.ToFen())
		}
		if pinLocs[k] == 64 {
			if result != 0 {
				t.Error("Found a false pin")
			}
		} else if pinLocs[k] != uint8(bits.TrailingZeros64(result)) {
			t.Error("Wrong pinned location for ", b.ToFen(), ":",
				IndexToAlgebraic(Square(bits.TrailingZeros64(result))),
				"not", IndexToAlgebraic(Square(pinLocs[k])))
		}
	}
}

func TestOrthoPins(t *testing.T) {
	positions := map[string]int{
		"4k3/8/4r3/4Q3/1q6/2Q5/8/4K3 b - - 0 0": 2,
		"7k/8/8/8/1r2R3/8/8/4K3 w - - 0 0":      0, // "false pin"
		"7k/8/8/8/1r2R3/8/8/4K3 b - - 0 0":      0, // no pin at all
		"3k4/8/3n4/8/8/8/3Q4/7K b - - 0 0":      0, // knight pin
		"8/8/1r3QK1/3QQ3/8/kr6/8/8 w - - 0 0":   4, // queen pin
		"4k3/4p3/8/8/8/4R3/q2PK3/8 w - - 0 0":   0, // horizontal pawn*/
		"4k3/4p3/8/8/8/4R3/q2PK3/8 b - - 0 0":   2,
		"8/4k3/8/4p3/8/4R3/q2PK3/8 b - - 0 0":   1,
	}
	pinLocs := map[string]uint8{
		"4k3/8/4r3/4Q3/1q6/2Q5/8/4K3 b - - 0 0": AlgebraicToIndex("e6"),
		"7k/8/8/8/1r2R3/8/8/4K3 w - - 0 0":      64, // "false pin"
		"7k/8/8/8/1r2R3/8/8/4K3 b - - 0 0":      64, // no pin at all
		"3k4/8/3n4/8/8/8/3Q4/7K b - - 0 0":      AlgebraicToIndex("d6"),
		"8/8/1r3QK1/3QQ3/8/kr6/8/8 w - - 0 0":   AlgebraicToIndex("f6"),
		"4k3/4p3/8/8/8/4R3/q2PK3/8 w - - 0 0":   AlgebraicToIndex("d2"), // horizontal
		"4k3/4p3/8/8/8/4R3/q2PK3/8 b - - 0 0":   AlgebraicToIndex("e7"),
		"8/4k3/8/4p3/8/4R3/q2PK3/8 b - - 0 0":   AlgebraicToIndex("e5"),
	}
	for k, v := range positions {
		moves := make([]Move, 0, 45)
		b := ParseFen(k)
		result := b.generatePinnedMoves(&moves)
		if len(moves) != v {
			t.Error("Legal moves for orthogonal pins: wrong length. Expected", v, "but got", len(moves), "for position", b.ToFen())
			printMoves(moves)
		}
		if pinLocs[k] == 64 {
			if result != 0 {
				t.Error("Found a false pin")
			}
		} else if pinLocs[k] != uint8(bits.TrailingZeros64(result)) {
			t.Error("Wrong pinned location")
		}
	}
}

// An incomplete, yet giant, test suite of positions. Tests legal move generation.
func testLegalMoves(t *testing.T) {
	positions := map[string]int{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1":             20,
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1": 48,
		"4k3/8/8/8/8/8/8/4K2R w K - 0 1":                                       15,
		"4k3/8/8/8/8/8/8/R3K3 w Q - 0 1":                                       16,
		"4k2r/8/8/8/8/8/8/4K3 w k - 0 1":                                       5,
		"r3k3/8/8/8/8/8/8/4K3 w q - 0 1":                                       5,
		"4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1":                                     26,
		"r3k2r/8/8/8/8/8/8/4K3 w kq - 0 1":                                     5,
		"8/8/8/8/8/8/6k1/4K2R w K - 0 1":                                       12,
		"8/8/8/8/8/8/1k6/R3K3 w Q - 0 1":                                       15,
		"4k2r/6K1/8/8/8/8/8/8 w k - 0 1":                                       3,
		"r3k3/1K6/8/8/8/8/8/8 w q - 0 1":                                       4,
		"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1":                                 26,
		"r3k2r/8/8/8/8/8/8/1R2K2R w Kkq - 0 1":                                 25,
		"r3k2r/8/8/8/8/8/8/2R1K2R w Kkq - 0 1":                                 25,
		"r3k2r/8/8/8/8/8/8/R3K1R1 w Qkq - 0 1":                                 25,
		"1r2k2r/8/8/8/8/8/8/R3K2R w KQk - 0 1":                                 26,
		"2r1k2r/8/8/8/8/8/8/R3K2R w KQk - 0 1":                                 25,
		"r3k1r1/8/8/8/8/8/8/R3K2R w KQq - 0 1":                                 25,
		"4k3/8/8/8/8/8/8/4K2R b K - 0 1":                                       5,
		"4k3/8/8/8/8/8/8/R3K3 b Q - 0 1":                                       5,
		"4k2r/8/8/8/8/8/8/4K3 b k - 0 1":                                       15,
		"r3k3/8/8/8/8/8/8/4K3 b q - 0 1":                                       16,
		"4k3/8/8/8/8/8/8/R3K2R b KQ - 0 1":                                     5,
		"r3k2r/8/8/8/8/8/8/4K3 b kq - 0 1":                                     26,
		"8/8/8/8/8/8/6k1/4K2R b K - 0 1":                                       3,
		"8/8/8/8/8/8/1k6/R3K3 b Q - 0 1":                                       4,
		"4k2r/6K1/8/8/8/8/8/8 b k - 0 1":                                       12,
		"r3k3/1K6/8/8/8/8/8/8 b q - 0 1":                                       15,
		"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1":                                 26,
		"r3k2r/8/8/8/8/8/8/1R2K2R b Kkq - 0 1":                                 26,
		"r3k2r/8/8/8/8/8/8/2R1K2R b Kkq - 0 1":                                 25,
		"r3k2r/8/8/8/8/8/8/R3K1R1 b Qkq - 0 1":                                 25,
		"1r2k2r/8/8/8/8/8/8/R3K2R b KQk - 0 1":                                 25,
		"2r1k2r/8/8/8/8/8/8/R3K2R b KQk - 0 1":                                 25,
		"r3k1r1/8/8/8/8/8/8/R3K2R b KQq - 0 1":                                 25,
		"8/1n4N1/2k5/8/8/5K2/1N4n1/8 w - - 0 1":                                14,
		"8/1k6/8/5N2/8/4n3/8/2K5 w - - 0 1":                                    11,
		"8/8/4k3/3Nn3/3nN3/4K3/8/8 w - - 0 1":                                  19,
		"K7/8/2n5/1n6/8/8/8/k6N w - - 0 1":                                     3,
		"k7/8/2N5/1N6/8/8/8/K6n w - - 0 1":                                     17,
		"8/1n4N1/2k5/8/8/5K2/1N4n1/8 b - - 0 1":                                15,
		"8/1k6/8/5N2/8/4n3/8/2K5 b - - 0 1":                                    16,
		"8/8/3K4/3Nn3/3nN3/4k3/8/8 b - - 0 1":                                  4,
		"K7/8/2n5/1n6/8/8/8/k6N b - - 0 1":                                     17,
		"k7/8/2N5/1N6/8/8/8/K6n b - - 0 1":                                     3,
		"B6b/8/8/8/2K5/4k3/8/b6B w - - 0 1":                                    17,
		"8/8/1B6/7b/7k/8/2B1b3/7K w - - 0 1":                                   21,
		"k7/B7/1B6/1B6/8/8/8/K6b w - - 0 1":                                    21,
		"K7/b7/1b6/1b6/8/8/8/k6B w - - 0 1":                                    7,
		"B6b/8/8/8/2K5/5k2/8/b6B b - - 0 1":                                    6,
		"8/8/1B6/7b/7k/8/2B1b3/7K b - - 0 1":                                   17,
		"k7/B7/1B6/1B6/8/8/8/K6b b - - 0 1":                                    7,
		"K7/b7/1b6/1b6/8/8/8/k6B b - - 0 1":                                    21,
		"7k/RR6/8/8/8/8/rr6/7K w - - 0 1":                                      19,
		"R6r/8/8/2K5/5k2/8/8/r6R w - - 0 1":                                    36,
		"7k/RR6/8/8/8/8/rr6/7K b - - 0 1":                                      19,
		"R6r/8/8/2K5/5k2/8/8/r6R b - - 0 1":                                    36,
		"6kq/8/8/8/8/8/8/7K w - - 0 1":                                         2,
		"K7/8/8/3Q4/4q3/8/8/7k w - - 0 1":                                      6,
		"6qk/8/8/8/8/8/8/7K b - - 0 1":                                         22,
		"6KQ/8/8/8/8/8/8/7k b - - 0 1":                                         2,
		"K7/8/8/3Q4/4q3/8/8/7k b - - 0 1":                                      6,
		"8/8/8/8/8/K7/P7/k7 w - - 0 1":                                         3,
		"8/8/8/8/8/7K/7P/7k w - - 0 1":                                         3,
		"K7/p7/k7/8/8/8/8/8 w - - 0 1":                                         1,
		"7K/7p/7k/8/8/8/8/8 w - - 0 1":                                         1,
		"8/2k1p3/3pP3/3P2K1/8/8/8/8 w - - 0 1":                                 7,
		"8/8/8/8/8/K7/P7/k7 b - - 0 1":                                         1,
		"8/8/8/8/8/7K/7P/7k b - - 0 1":                                         1,
		"K7/p7/k7/8/8/8/8/8 b - - 0 1":                                         3,
		"7K/7p/7k/8/8/8/8/8 b - - 0 1":                                         3,
		"8/2k1p3/3pP3/3P2K1/8/8/8/8 b - - 0 1":                                 5,
		"8/8/8/8/8/4k3/4P3/4K3 w - - 0 1":                                      2,
		"4k3/4p3/4K3/8/8/8/8/8 b - - 0 1":                                      2,
		"8/8/7k/7p/7P/7K/8/8 w - - 0 1":                                        3,
		"8/8/k7/p7/P7/K7/8/8 w - - 0 1":                                        3,
		"8/8/3k4/3p4/3P4/3K4/8/8 w - - 0 1":                                    5,
		"8/3k4/3p4/8/3P4/3K4/8/8 w - - 0 1":                                    8,
		"8/8/3k4/3p4/8/3P4/3K4/8 w - - 0 1":                                    8,
		"k7/8/3p4/8/3P4/8/8/7K w - - 0 1":                                      4,
		"8/8/7k/7p/7P/7K/8/8 b - - 0 1":                                        3,
		"8/8/k7/p7/P7/K7/8/8 b - - 0 1":                                        3,
		"8/8/3k4/3p4/3P4/3K4/8/8 b - - 0 1":                                    5,
		"8/3k4/3p4/8/3P4/3K4/8/8 b - - 0 1":                                    8,
		"8/8/3k4/3p4/8/3P4/3K4/8 b - - 0 1":                                    8,
		"k7/8/3p4/8/3P4/8/8/7K b - - 0 1":                                      4,
		"7k/3p4/8/8/3P4/8/8/K7 w - - 0 1":                                      4,
		"7k/8/8/3p4/8/8/3P4/K7 w - - 0 1":                                      5,
		"k7/8/8/7p/6P1/8/8/K7 w - - 0 1":                                       5,
		"k7/8/7p/8/8/6P1/8/K7 w - - 0 1":                                       4,
		"k7/8/8/6p1/7P/8/8/K7 w - - 0 1":                                       5,
		"k7/8/6p1/8/8/7P/8/K7 w - - 0 1":                                       4,
		"k7/8/8/3p4/4p3/8/8/7K w - - 0 1":                                      3,
		"k7/8/3p4/8/8/4P3/8/7K w - - 0 1":                                      4,
		"7k/3p4/8/8/3P4/8/8/K7 b - - 0 1":                                      5,
		"7k/8/8/3p4/8/8/3P4/K7 b - - 0 1":                                      4,
		"k7/8/8/7p/6P1/8/8/K7 b - - 0 1":                                       5,
		"k7/8/7p/8/8/6P1/8/K7 b - - 0 1":                                       4,
		"k7/8/8/6p1/7P/8/8/K7 b - - 0 1":                                       5,
		"k7/8/6p1/8/8/7P/8/K7 b - - 0 1":                                       4,
		"k7/8/8/3p4/4p3/8/8/7K b - - 0 1":                                      5,
		"k7/8/3p4/8/8/4P3/8/7K b - - 0 1":                                      4,
		"7k/8/8/p7/1P6/8/8/7K w - - 0 1":                                       5,
		"7k/8/p7/8/8/1P6/8/7K w - - 0 1":                                       4,
		"7k/8/8/1p6/P7/8/8/7K w - - 0 1":                                       5,
		"7k/8/1p6/8/8/P7/8/7K w - - 0 1":                                       4,
		"k7/7p/8/8/8/8/6P1/K7 w - - 0 1":                                       5,
		"k7/6p1/8/8/8/8/7P/K7 w - - 0 1":                                       5,
		"3k4/3pp3/8/8/8/8/3PP3/3K4 w - - 0 1":                                  7,
		"7k/8/8/p7/1P6/8/8/7K b - - 0 1":                                       5,
		"7k/8/p7/8/8/1P6/8/7K b - - 0 1":                                       4,
		"7k/8/8/1p6/P7/8/8/7K b - - 0 1":                                       5,
		"7k/8/1p6/8/8/P7/8/7K b - - 0 1":                                       4,
		"k7/7p/8/8/8/8/6P1/K7 b - - 0 1":                                       5,
		"k7/6p1/8/8/8/8/7P/K7 b - - 0 1":                                       5,
		"3k4/3pp3/8/8/8/8/3PP3/3K4 b - - 0 1":                                  7,
		"8/Pk6/8/8/8/8/6Kp/8 w - - 0 1":                                        11,
		"n1n5/1Pk5/8/8/8/8/5Kp1/5N1N w - - 0 1":                                24,
		"8/PPPk4/8/8/8/8/4Kppp/8 w - - 0 1":                                    18,
		"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N w - - 0 1":                              24,
		"8/Pk6/8/8/8/8/6Kp/8 b - - 0 1":                                        11,
		"n1n5/1Pk5/8/8/8/8/5Kp1/5N1N b - - 0 1":                                24,
		"8/PPPk4/8/8/8/8/4Kppp/8 b - - 0 1":                                    18,
		"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1":                              24}
	for k, v := range positions {
		b := ParseFen(k)
		moves := b.GenerateLegalMoves()
		if len(moves) != v {
			t.Error("Legal moves: wrong length. Expected", v, "but got", len(moves), "for position", b.ToFen())
		}
	}
}
