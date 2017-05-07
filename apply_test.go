package dragontoothmg

import (
	"testing"
)

func testApplyUnapply(t *testing.T) {
	movesMap := map[string]Move{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0": parseMove("e2e4"),
	}
	results := map[string]string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 0",
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
			t.Error("Move application didn't produce expected result for", k)
		}
		unapply()
		if k != b.ToFen() {
			t.Error("Board changed during unapply for", k)
		}
		for _, mv := range movesList {
			unapply := b.Apply(mv)
			unapply()
			if b.ToFen() != k {
				t.Error("Move apply/unapply changed board", mv, k)
			}
		}
	}
}
