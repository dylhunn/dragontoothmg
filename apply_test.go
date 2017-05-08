package dragontoothmg

import (
	"testing"
)

func TestApplyUnapply(t *testing.T) {
	movesMap := map[string]Move{
		// ordinary move
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0": parseMove("e2e4"),
		// castle 1: white short
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 0": parseMove("e1g1"),
		// castle 2: black long, without kingside rights
		"r3kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R b KQq - 0 0": parseMove("e8c8"),
		// rook move strips castling rights
		// en passant capture
		// promotion 1
		// promotion 2
		// promotion-capture
		// capture 1
		// capture
	}
	results := map[string]string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 0",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 0":  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1RK1 b kq - 0 0",
		"r3kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R b KQq - 0 0":     "2kr1bnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQ - 0 0",
	}
	for k, v := range movesMap {
		b := ParseFen(k)
		fenBefore := b.ToFen()
		movesList := b.GenerateLegalMoves()
		fenAfter := b.ToFen()
		if fenBefore != k {
			t.Error("Fen changed during parsing for board", k)
		}
		if fenBefore != fenAfter {
			t.Error("Fen changed during generation for board", k)
		}
		unapply := b.Apply(v)
		if b.ToFen() != results[k] {
			t.Error("Move application didn't produce expected result for\n", k, "->\n",
				results[k], "\nInstead, we got:\n", b.ToFen())
		}
		unapply()
		if k != b.ToFen() {
			t.Error("Board changed during unapply for\n", k, "\nResult was\n", b.ToFen())
		}
		for _, mv := range movesList {
			unapply := b.Apply(mv)
			unapply()
			if b.ToFen() != k {
				t.Error("Move apply/unapply changed board\n", &mv, "\n", k)
			}
		}
	}
}
