package dragontoothmg

// Applies a move to the board, and returns a function that can be used to unapply it.
func (b *Board) Apply(m Move) func() {
	// Configure data about which pieces move
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
	pieceType, pieceTypeBitboard := determinePieceType(ourBitboardPtr, fromBitboard)
	castleStatus := 0
	var oldRookLoc, newRookLoc uint8
	kingsideCastleRightsBefore := b.canCastleKingside()
	queensideCastleRightsBefore := b.canCastleQueenside()
	var flippedKsCastle, flippedQsCastle bool

	// Configure handling castling rights
	if pieceType == King && m.To()-m.From() == 2 { // castle short
		castleStatus = 1
		oldRookLoc = m.To() + 1
		newRookLoc = m.To() - 1
		b.flipKingsideCastle()
		if queensideCastleRightsBefore {
			b.flipQueensideCastle()
			flippedQsCastle = true
		}
		flippedKsCastle = true
	} else if pieceType == King && int(m.To())-int(m.From()) == -2 { // castle long
		castleStatus = -1
		oldRookLoc = m.To() - 2
		newRookLoc = m.To() + 1
		b.flipQueensideCastle()
		if kingsideCastleRightsBefore {
			b.flipKingsideCastle()
			flippedKsCastle = true
		}
		flippedQsCastle = true
	}
	// Apply the castling rook movement
	if castleStatus != 0 {
		ourBitboardPtr.rooks |= (uint64(1) << newRookLoc)
		ourBitboardPtr.rooks &= ^(uint64(1) << oldRookLoc)
	}

	// Rook moves strip castling rights
	/*if (pieceType == Rook) {
		if (pieceTypeBitboard & onlyFile[7] != 0) { 

		}
	}*/

	// Apply the move
	ourBitboardPtr.all &= ^fromBitboard // remove at "from"
	ourBitboardPtr.all |= toBitboard    // add at "to"
	*pieceTypeBitboard &= ^fromBitboard // remove at "from"
	*pieceTypeBitboard |= toBitboard    // add at "to"
	capturedPieceType, capturedBitboard := determinePieceType(oppBitboardPtr, toBitboard)
	if capturedPieceType != Nothing {
		*capturedBitboard &= ^toBitboard
		oppBitboardPtr.all &= ^toBitboard
	}
	b.wtomove = !b.wtomove

	// Return the unapply function (closure)
	unapply := func() {
		ourBitboardPtr.all &= ^toBitboard
		ourBitboardPtr.all |= fromBitboard
		*pieceTypeBitboard &= ^toBitboard
		*pieceTypeBitboard |= fromBitboard
		if capturedPieceType != Nothing {
			*capturedBitboard |= toBitboard
			oppBitboardPtr.all |= toBitboard
		}
		if castleStatus != 0 {
			ourBitboardPtr.rooks &= ^(uint64(1) << newRookLoc)
			ourBitboardPtr.rooks |= (uint64(1) << oldRookLoc)
		}
		b.wtomove = !b.wtomove
		// must update castling flags after turn swap
		if flippedKsCastle {
			b.flipKingsideCastle()
		}
		if flippedQsCastle {
			b.flipQueensideCastle()
		}

	}
	return unapply
}

func determinePieceType(ourBitboardPtr *bitboards, squareMask uint64) (Piece, *uint64) {
	var pieceType Piece = Nothing
	pieceTypeBitboard := &(ourBitboardPtr.all)
	if squareMask&ourBitboardPtr.pawns != 0 {
		pieceType = Pawn
		pieceTypeBitboard = &(ourBitboardPtr.pawns)
	} else if squareMask&ourBitboardPtr.knights != 0 {
		pieceType = Knight
		pieceTypeBitboard = &(ourBitboardPtr.knights)
	} else if squareMask&ourBitboardPtr.bishops != 0 {
		pieceType = Bishop
		pieceTypeBitboard = &(ourBitboardPtr.bishops)
	} else if squareMask&ourBitboardPtr.rooks != 0 {
		pieceType = Rook
		pieceTypeBitboard = &(ourBitboardPtr.rooks)
	} else if squareMask&ourBitboardPtr.queens != 0 {
		pieceType = Queen
		pieceTypeBitboard = &(ourBitboardPtr.queens)
	} else if squareMask&ourBitboardPtr.kings != 0 {
		pieceType = King
		pieceTypeBitboard = &(ourBitboardPtr.kings)
	}
	return pieceType, pieceTypeBitboard
}
