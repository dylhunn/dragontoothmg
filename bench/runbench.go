package main

import (
	"fmt"
	"github.com/dylhunn/dragontoothmg"
	"testing"
	"time"
)

func main() {
	res := testing.Benchmark(benchmarkStartpos)
	fmt.Println(res)
	fmt.Println("Ns per op:", res.NsPerOp())
	fmt.Println("Time per op:", time.Duration(res.NsPerOp()))
}

// -----------------
// BENCHMARK HELPERS
// -----------------

func benchmarkStartpos(b *testing.B) {
	board := dragontoothmg.ParseFen(dragontoothmg.Startpos)
	for i := 0; i < b.N; i++ {
		dragontoothmg.Perft(&board, 5)
	}
}

func benchmarkKiwipete(b *testing.B) {
	pos := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0"
	board := dragontoothmg.ParseFen(pos)
	for i := 0; i < b.N; i++ {
		dragontoothmg.Perft(&board, 4)
	}
}

func benchmarkDense(b *testing.B) {
	pos := "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
	board := dragontoothmg.ParseFen(pos)
	for i := 0; i < b.N; i++ {
		dragontoothmg.Perft(&board, 5)
	}
}

func benchmarkEndgameRP(b *testing.B) {
	pos := "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 0"
	board := dragontoothmg.ParseFen(pos)
	for i := 0; i < b.N; i++ {
		dragontoothmg.Perft(&board, 6)
	}
}
