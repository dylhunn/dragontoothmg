package main

import (
	"fmt"
	"github.com/dylhunn/dragontoothmg"
	"testing"
	//"time"
)

const nsPerMs = 1000000
const nsPerS = 1000000000

func main() {
	fmt.Println("\nBENCHMARKS")
	printResultLine(testing.Benchmark(benchmarkStartpos5), "Start position", startposResult5, 5)
	printResultLine(testing.Benchmark(benchmarkStartpos6), "Start position", startposResult6, 6)
	printResultLine(testing.Benchmark(benchmarkKiwipete), "Kiwipete position", kpResult, 5)
	printResultLine(testing.Benchmark(benchmarkDense), "Dense position", denseResult, 6)
	printResultLine(testing.Benchmark(benchmarkEndgameRP), "Endgame R/P position", endgameResult, 6)
	fmt.Println()
}

func printResultLine(res testing.BenchmarkResult, name string, perftValue int64, depth int) {
	fmt.Printf("%-22s depth %-3d %8dms %12d nodes  %11.0fnps\n", name + ":", depth, res.NsPerOp() / nsPerMs,
		perftValue, float64(perftValue) / (float64(res.NsPerOp()) / nsPerS))
}

// -----------------
// BENCHMARK HELPERS
// -----------------

var startposResult5 int64 = 0
func benchmarkStartpos5(b *testing.B) {
	pos := dragontoothmg.Startpos
	board := dragontoothmg.ParseFen(pos)
	for i := 0; i < b.N; i++ {
		startposResult5 = dragontoothmg.Perft(&board, 5)
	}
}

var startposResult6 int64 = 0
func benchmarkStartpos6(b *testing.B) {
	pos := dragontoothmg.Startpos
	board := dragontoothmg.ParseFen(pos)
	for i := 0; i < b.N; i++ {
		startposResult6 = dragontoothmg.Perft(&board, 6)
	}
}

var kpResult int64 = 0
func benchmarkKiwipete(b *testing.B) {
	pos := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0"
	board := dragontoothmg.ParseFen(pos)
	for i := 0; i < b.N; i++ {
		kpResult = dragontoothmg.Perft(&board, 5)
	}
}

var denseResult int64 = 0
func benchmarkDense(b *testing.B) {
	pos := "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
	board := dragontoothmg.ParseFen(pos)
	for i := 0; i < b.N; i++ {
		denseResult = dragontoothmg.Perft(&board, 6)
	}
}

var endgameResult int64 = 0
func benchmarkEndgameRP(b *testing.B) {
	pos := "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 0"
	board := dragontoothmg.ParseFen(pos)
	for i := 0; i < b.N; i++ {
		endgameResult = dragontoothmg.Perft(&board, 6)
	}
}
