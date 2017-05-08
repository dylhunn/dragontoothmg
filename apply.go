package dragontoothmg

// Applies a move to the board, and returns a function that can be used to unapply it.
func (b *Board) Apply(m Move) func() {
	var ourBitboardPtr, oppBitboardPtr *bitboards
	if b.wtomove {
		ourBitboardPtr = &(b.white)
		oppBitboardPtr = &(b.black)
	} else {
		ourBitboardPtr = &(b.black)
		oppBitboardPtr = &(b.white)
	}
	fromBitboard := (uint64(1) << m.From())
	toBitboard := (uint64(1) << m.To())
	_, pieceTypeBitboard := determinePieceType(ourBitboardPtr, fromBitboard)
	ourBitboardPtr.all &= ^fromBitboard // remove at "from"
	ourBitboardPtr.all |= toBitboard // add at "to"
	*pieceTypeBitboard &= ^fromBitboard // remove at "from"
	*pieceTypeBitboard |= toBitboard // add at "to"
	capturedPieceType, capturedBitboard := determinePieceType(oppBitboardPtr, toBitboard)
	if (capturedPieceType != Nothing) {
		*capturedBitboard &= ^toBitboard
		oppBitboardPtr.all  &= ^toBitboard
	}

	unapply := func() {
		ourBitboardPtr.all &= ^toBitboard
		ourBitboardPtr.all |= fromBitboard
		*pieceTypeBitboard &= ^toBitboard
		*pieceTypeBitboard |= fromBitboard
		if (capturedPieceType != Nothing) {
			*capturedBitboard |= toBitboard
			oppBitboardPtr.all  |= toBitboard
		}
	}
	return unapply
}

func determinePieceType(ourBitboardPtr *bitboards, squareMask uint64) (Piece, *uint64) {
	var pieceType Piece = Nothing
	pieceTypeBitboard := &(ourBitboardPtr.all)
	if squareMask & ourBitboardPtr.pawns != 0 {
		pieceType = Pawn
		pieceTypeBitboard = &(ourBitboardPtr.pawns)
	} else if squareMask & ourBitboardPtr.knights != 0 {
		pieceType = Knight
		pieceTypeBitboard = &(ourBitboardPtr.knights)
	} else if squareMask & ourBitboardPtr.bishops != 0 {
		pieceType = Bishop
		pieceTypeBitboard = &(ourBitboardPtr.bishops)
	} else if squareMask & ourBitboardPtr.rooks != 0 {
		pieceType = Rook
		pieceTypeBitboard = &(ourBitboardPtr.rooks)
	} else if squareMask & ourBitboardPtr.queens != 0 {
		pieceType = Queen
		pieceTypeBitboard = &(ourBitboardPtr.queens)
	} else if squareMask & ourBitboardPtr.kings != 0 {
		pieceType = King
		pieceTypeBitboard = &(ourBitboardPtr.kings)
	}
	return pieceType, pieceTypeBitboard
}