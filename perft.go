package dragontoothmg

import "fmt"

// Run perft to count the number of moves.
// Useful for testing and benchmarking.
func Perft(b *Board, n int) int {
	if n <= 0 {
		return 1
	}
	moves := b.GenerateLegalMoves()
	if n == 1 {
		return len(moves)
	}
	var count int = 0
	
	for _, move := range moves {
		unapply := b.Apply(move)
		count += Perft(b, n-1)
		unapply()

	}
	return count
}

// Performs the Perft move count division operation. Useful for debugging.
func Divide(b *Board, n int) {
	moves := b.GenerateLegalMoves()
	for i, move := range moves {
		unapply := b.Apply(move)
		result := Perft(b, n-1)
		unapply()
		fmt.Printf("Move   #%3d:   %-6s =%9d\n", i+1, &move, result)
	}
}
