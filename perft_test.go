package dragontoothmg

import (
	"testing"
)

func TestDivide(t *testing.T) {
	//b := ParseFen("rnbqkbnr/1ppppppp/p7/P7/8/8/1PPPPPPP/RNBQKBNR b KQkq - 0 1") // my b7b5 finds 21 instead of 22 moves
	b := ParseFen(startpos)
	Divide(&b, 1)
}

func testStartingPosition(t *testing.T) {
	b := ParseFen(startpos)
	perftSolutions := map[int]int64{
		1: 20,
		2: 400,
		3: 8902,
		4: 197281,
		5: 4865609,
		6: 119060324,
	}
	for i := 1; i <= len(perftSolutions); i++ {
		result := Perft(&b, i)
		if result != perftSolutions[i] {
			t.Error("Starting position perft error.\nExpected",
				perftSolutions[i], "but got", result, "for depth", i)
		}
	}
}

func TestForPromotionBugs(t *testing.T) {

}

func BenchmarkFoo(b *testing.B) {

}
