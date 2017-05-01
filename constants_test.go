package movegen

import (
	"testing"
)

func TestRookMovesFromBlockers(t *testing.T) {
	// Blockers:
	// 00010000
	// 00000000
	// 00010000
	// 000R0010
	// 00000000
	// 00010000
	// 00010000
	// 00000000
	// Bitstring: 0000100000000000000010000100000000000000000010000000100000000000
	// Bitstring: 0x800084000080800
	// Rook at D5 = 35
	// Moves:
	// 00000000
	// 00000000
	// 00010000
	// 11101110
	// 00010000
	// 00010000
	// 00000000
	// 00000000
	// Bitstring: 0000000000000000000010000111011100001000000010000000000000000000
	// Bitstring: 0x87708080000
	moves := rookMovesFromBlockers(35, 0x800084000080800)
	if moves != 0x87708080000 {
		t.Error("Failed to generate rook moves from blocker board. Output:", moves)
	}
}
