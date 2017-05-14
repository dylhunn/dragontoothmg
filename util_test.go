package dragontoothmg

import (
	"testing"
)

// Some example valid move strings:
// e1e2 b4d6 e7e8q a2a1n
func TestParseMove(t *testing.T) {
	move, _ := ParseMove("b4d6")
	if move.To() != algebraicToIndexFatal("d6") ||
		move.From() != algebraicToIndexFatal("b4") ||
		move.Promote() != Nothing {
		t.Error("Incorrectly parsed move.")
	}
	move2, _ := ParseMove("a2a1n")
	if move2.To() != algebraicToIndexFatal("a1") ||
		move2.From() != algebraicToIndexFatal("a2") ||
		move2.Promote() != Knight {
		t.Error("Incorrectly parsed move.")
	}
}

func TestAlgToIdx(t *testing.T) {
	if algebraicToIndexFatal("A8") != 56 {
		t.Error("Algebraic to index conversion failed.")
	}
	if algebraicToIndexFatal("A1") != 0 {
		t.Error("Algebraic to index conversion failed.")
	}
	if algebraicToIndexFatal("h3") != 23 {
		t.Error("Algebraic to index conversion failed.")
	}
	if algebraicToIndexFatal("a6") != 40 {
		t.Error("Algebraic to index conversion failed.")
	}
	if algebraicToIndexFatal("H4") != 31 {
		t.Error("Algebraic to index conversion failed.")
	}
	_, err := AlgebraicToIndex("H9")
	if err == nil {
		t.Error("Algebraic to index conversion failed.")
	}
	_, err2 := AlgebraicToIndex("qq")
	if err2 == nil {
		t.Error("Algebraic to index conversion failed.")
	}
}

func TestIdxToAlg(t *testing.T) {
	if IndexToAlgebraic(56) != "a8" {
		t.Error("Index to algebraic conversion failed:", IndexToAlgebraic(56), "instead of a8")
	}
	if IndexToAlgebraic(0) != "a1" {
		t.Error("Index to algebraic conversion failed:", IndexToAlgebraic(0), "instead of a1")
	}
	if IndexToAlgebraic(40) != "a6" {
		t.Error("Index to algebraic conversion failed:", IndexToAlgebraic(31), "instead of a6")
	}
	if IndexToAlgebraic(31) != "h4" {
		t.Error("Index to algebraic conversion failed:", IndexToAlgebraic(40), "instead of h4")
	}
}

func TestParseFen(t *testing.T) {
	b := ParseFen("1Q2rk2/2p2p2/1n4b1/N7/2B1Pp1q/2B4P/1QPP4/4K2R b K e3 4 30")
	if b.Wtomove {
		t.Error("Error parsing FEN")
	}
	if b.enpassant != 20 {
		t.Error("Error parsing FEN")
	}
	if !b.whiteCanCastleKingside() {
		t.Error("Error parsing FEN")
	}
	if b.whiteCanCastleQueenside() {
		t.Error("Error parsing FEN")
	}
	if b.blackCanCastleKingside() {
		t.Error("Error parsing FEN")
	}
	if b.blackCanCastleQueenside() {
		t.Error("Error parsing FEN")
	}
	if b.White.Kings != 1<<4 {
		t.Error("Error parsing FEN")
	}
	if b.Black.Kings != 1<<61 {
		t.Error("Error parsing FEN")
	}
	if b.White.Rooks != 1<<7 {
		t.Error("Error parsing FEN")
	}
	if b.White.Knights != 1<<32 {
		t.Error("Error parsing FEN")
	}
	if b.Halfmoveclock != 4 {
		t.Error("Error parsing FEN")
	}
	if b.Fullmoveno != 30 {
		t.Error("Error parsing FEN")
	}
}

func TestToFen(t *testing.T) {
	fenTests := []string{
		"1Q2rk2/2p2p2/1n4b1/N7/2B1Pp1q/2B4P/1QPP4/4K2R b K e3 4 30",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 10",
		"6nq/6p1/2B4n/1rB2r1R/5q2/2P5/1Q4n1/2B5 w - h8 6 12",
		"6nq/6p1/2B4n/1rB2r1R/5q2/2P5/1Q4n1/2B5 b - - 2 999"}
	for _, fen := range fenTests {
		b := ParseFen(fen)
		if b.ToFen() != fen {
			t.Error("Error serializing FEN.\nOutput:  ", b.ToFen(), "\nExpected:", fen)
		}
	}
}
