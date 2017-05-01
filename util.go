package movegen

import (
	//"math/bits"
	"strconv"
	"strings"
)

// Accepts an algebraic notation chess square, and converts it to a square ID
// as used by Dragontooth (in both the board and move types).
func AlgebraicToIndex(alg string) uint8 {
	return (strings.ToLower(alg)[0] - 'a') + ((alg[1] - '1') * 8)
}

/*func generateRookMagicTable() {
	// For a rook at every board position
	for i := 0; i < 64; i++ {
		blockerMask := magicRookBlockerMasks[i]
		generateRookBlockerPermutations(Square(i), blockerMask, 0)
	}
}

// Recursively generate all permutations of active and inactive bits in the
// blocker mask. Origin is the piece's starting square. BlockerMaskProgress is
// the original blocker bitboard, from which we eliminate bits.
// CurrPerm is the bitboard permutation we are accumulating.
func generateRookBlockerPermutations(origin Square, blockerMaskProgress uint64, currPerm uint64) {
	if blockerMaskProgress == 0 {
		// the currPerm represents one possible occupancy pattern on the rook blocker bitboard
		dbindex := (currPerm * magicNumberRook[origin]) >> magicRookShifts[origin]

		magicMovesRook[origin][dbindex] = moves
		return
	}
	nextBit := bits.TrailingZeros64(blockerMaskProgress)
	blockerMaskProgress &= blockerMaskProgress - 1
	without := currPerm
	with := currPerm | (uint64(1) << uint8(nextBit))
	generateRookBlockerPermutations(origin, blockerMaskProgress, without)
	generateRookBlockerPermutations(origin, blockerMaskProgress, with)
}*/

// Parse a board from a FEN string.
func ParseFen(fen string) Board {
	// BUG(dylhunn): This FEN parsing implementation doesn't handle malformed inputs.
	tokens := strings.Fields(fen)
	var b Board
	// replace digits with the appropriate number of dashes
	for i := 1; i <= 8; i++ {
		var replacement string
		for j := 0; j < i; j++ {
			replacement += "-"
		}
		tokens[0] = strings.Replace(tokens[0], strconv.Itoa(i), replacement, -1)
	}
	// reverse the order of the ranks, removing slashes
	ranks := strings.Split(tokens[0], "/")
	for i := 0; i < len(ranks)/2; i++ {
		j := len(ranks) - i - 1
		ranks[i], ranks[j] = ranks[j], ranks[i]
	}
	tokens[0] = ranks[0]
	for i := 1; i < len(ranks); i++ {
		tokens[0] += ranks[i]
	}
	// add every piece to the board
	for i := uint8(0); i < 64; i++ {
		switch tokens[0][i] {
		case 'p':
			b.black.pawns |= 1 << i
		case 'n':
			b.black.knights |= 1 << i
		case 'b':
			b.black.bishops |= 1 << i
		case 'r':
			b.black.rooks |= 1 << i
		case 'q':
			b.black.queens |= 1 << i
		case 'k':
			b.black.kings |= 1 << i
		case 'P':
			b.white.pawns |= 1 << i
		case 'N':
			b.white.knights |= 1 << i
		case 'B':
			b.white.bishops |= 1 << i
		case 'R':
			b.white.rooks |= 1 << i
		case 'Q':
			b.white.queens |= 1 << i
		case 'K':
			b.white.kings |= 1 << i
		}
	}
	b.white.all = b.white.pawns | b.white.knights | b.white.bishops | b.white.rooks | b.white.queens | b.white.kings
	b.black.all = b.black.pawns | b.black.knights | b.black.bishops | b.black.rooks | b.black.queens | b.black.kings

	b.wtomove = tokens[1] == "w" || tokens[1] == "W"
	if strings.Contains(tokens[2], "K") {
		b.flipWhiteKingsideCastle()
	}
	if strings.Contains(tokens[2], "Q") {
		b.flipWhiteQueensideCastle()
	}
	if strings.Contains(tokens[2], "k") {
		b.flipBlackKingsideCastle()
	}
	if strings.Contains(tokens[2], "q") {
		b.flipBlackQueensideCastle()
	}
	if tokens[3] != "-" {
		b.enpassant = AlgebraicToIndex(tokens[3])
	}
	return b
}
