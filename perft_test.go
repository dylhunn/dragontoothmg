package dragontoothmg

import (
	"testing"
)

func testDivide(t *testing.T) {
	//b := ParseFen("rnbqkbnr/1ppppppp/p7/P7/8/8/1PPPPPPP/RNBQKBNR b KQkq - 0 1") // my b7b5 finds 21 instead of 22 moves
	//b := ParseFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0")
	b := ParseFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	//b := ParseFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/4P3/1pN2Q1p/PPPBBPPP/R4RK1 w kq - 0 1")
	Divide(&b, 4)
}

func TestStartingPosition(t *testing.T) {
	perftSolutions := map[int]int64{
		1: 20,
		2: 400,
		3: 8902,
		4: 197281,
		5: 4865609,
		//6: 119060324,
	}
	checkPerftResults(startpos, perftSolutions, t)
}

func TestKiwipetePosition(t *testing.T) {
	perftSolutions := map[int]int64{
		1: 48,
		2: 2039,
		3: 97862,
		4: 4085603,
		//5: 193690690,
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
		//7: 178633661,
	}
	pos := "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 0"
	checkPerftResults(pos, perftSolutions, t)
}

func testMidgameDense(t *testing.T) {
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

func checkPerftResults(fen string, perftSolutions map[int]int64, t *testing.T) {
	b := ParseFen(fen)
	for i := 1; i <= len(perftSolutions); i++ {
		result := Perft(&b, i)
		if result != perftSolutions[i] {
			t.Error("Perft error in position\n", b.ToFen(), "\nExpected",
				perftSolutions[i], "but got", result, "for depth", i)
		}
	}
}

func BenchmarkFoo(b *testing.B) {

}
