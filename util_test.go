package dragontoothmg

import (
	"testing"
)

func TestAlgToIdx(t *testing.T) {
	if AlgebraicToIndex("A8") != 56 {
		t.Error("Algebraic to index conversion failed.")
	}
	if AlgebraicToIndex("A1") != 0 {
		t.Error("Algebraic to index conversion failed.")
	}
	if AlgebraicToIndex("H3") != 23 {
		t.Error("Algebraic to index conversion failed.")
	}
}

func TestIdxToAlg(t *testing.T) {
	if IndexToAlgebraic(56) != "A8" {
		t.Error("Index to algebraic conversion failed:", IndexToAlgebraic(56))
	}
	if IndexToAlgebraic(0) != "A1" {
		t.Error("Index to algebraic conversion failed:", IndexToAlgebraic(0))
	}
	if IndexToAlgebraic(23) != "H3" {
		t.Error("Index to algebraic conversion failed:", IndexToAlgebraic(23))
	}
}

func TestParseFen(t *testing.T) {
	b := ParseFen("1Q2rk2/2p2p2/1n4b1/N7/2B1Pp1q/2B4P/1QPP4/4K2R b K e3 4 30")
	if b.wtomove {
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
	if b.white.kings != 1<<4 {
		t.Error("Error parsing FEN")
	}
	if b.black.kings != 1<<61 {
		t.Error("Error parsing FEN")
	}
	if b.white.rooks != 1<<7 {
		t.Error("Error parsing FEN")
	}
	if b.white.knights != 1<<32 {
		t.Error("Error parsing FEN")
	}

}
