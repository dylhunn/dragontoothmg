package dragontoothmg

import (
	"testing"
)

func testDivide(t *testing.T) {
	b := ParseFen("nqn5/P1Pk4/8/8/8/6K1/7p/5N2 w - - 0 1")
	Divide(&b, 1)
}

func TestStartingPosition(t *testing.T) {
	perftSolutions := map[int]int64{
		1: 20,
		2: 400,
		3: 8902,
		4: 197281,
		5: 4865609,
		6: 119060324,
	}
	checkPerftResults(startpos, perftSolutions, t)
}

func TestKiwipetePosition(t *testing.T) {
	perftSolutions := map[int]int64{
		1: 48,
		2: 2039,
		3: 97862,
		4: 4085603,
		5: 193690690,
	}
	pos := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0"
	checkPerftResults(pos, perftSolutions, t)
}

func TestEndgameRP(t *testing.T) {
	perftSolutions := map[int]int64{
		1: 14,
		2: 191,
		3: 2812,
		4: 43238,
		5: 674624,
		6: 11030083,
		7: 178633661,
	}
	pos := "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 0"
	checkPerftResults(pos, perftSolutions, t)
}

func TestMidgameDense(t *testing.T) {
	perftSolutions := map[int]int64{
		1: 6,
		2: 264,
		3: 9467,
		4: 422333,
		5: 15833292,
		6: 706045033,
	}
	pos := "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
	checkPerftResults(pos, perftSolutions, t)
}

func TestPromotions(t *testing.T) {
	perftSolutions := map[int]int64{
		1: 24,
		2: 496,
		3: 9483,
		4: 182838,
		5: 3605103,
		6: 71179139,
	}
	pos := "n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1"
	checkPerftResults(pos, perftSolutions, t)
}

func checkPerftResults(fen string, perftSolutions map[int]int64, t *testing.T) {
	b := ParseFen(fen)
	for i := 1; i <= len(perftSolutions); i++ {
		beforeFen := b.ToFen()
		result := Perft(&b, i)
		afterFen := b.ToFen()
		if beforeFen != afterFen {
			t.Error("Perft corrupted board state.")
		}
		if result != perftSolutions[i] {
			t.Error("Perft error in position\n", b.ToFen(), "\nExpected",
				perftSolutions[i], "but got", result, "for depth", i)
		}
	}
}

func BenchmarkFoo(b *testing.B) {

}
