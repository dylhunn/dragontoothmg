package dragontoothmg

// Applies a move to the board, and returns a function that can be used to unapply it.
func (b *Board) Apply(m Move) func() {
	// Configure data about which pieces move
	var ourBitboardPtr, oppBitboardPtr *bitboards
	var epDelta int8 // add this to the e.p. square to find the captured pawn
	if b.wtomove {
		ourBitboardPtr = &(b.white)
		oppBitboardPtr = &(b.black)
		epDelta = -8
	} else {
		ourBitboardPtr = &(b.black)
		oppBitboardPtr = &(b.white)
		epDelta = 8
		b.fullmoveno++ // increment after black's move
	}
	fromBitboard := (uint64(1) << m.From())
	toBitboard := (uint64(1) << m.To())
	pieceType, pieceTypeBitboard := determinePieceType(ourBitboardPtr, fromBitboard)
	castleStatus := 0
	var oldRookLoc, newRookLoc uint8
	kingsideCastleRightsBefore := b.canCastleKingside()
	queensideCastleRightsBefore := b.canCastleQueenside()
	var flippedKsCastle, flippedQsCastle bool

	// Configure castling rights
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
	if pieceType == Rook {
		originBitboard := uint64(1) << m.From()
		if b.canCastleKingside() && (originBitboard & onlyFile[7] != 0) { // king's rook
			flippedKsCastle = true
			b.flipKingsideCastle()
		} else if b.canCastleQueenside() && (originBitboard & onlyFile[0] != 0) { // queen's rook
			flippedQsCastle = true
			b.flipQueensideCastle()
		}
	}

	// Is this an e.p. capture? Strip the opponent pawn and reset the e.p. square
	epCaptureSquare := b.enpassant
	if (epCaptureSquare != 0) {
		oppBitboardPtr.pawns &= ^(uint64(1) << uint8(int8(epCaptureSquare) + epDelta))
		oppBitboardPtr.all &= ^(uint64(1) << uint8(int8(epCaptureSquare) + epDelta))
		b.enpassant = 0
	}

	// Is this a promotion?
	var destTypeBitboard *uint64
	switch m.Promote() {
	case Queen:
		destTypeBitboard = &(ourBitboardPtr.queens)
	case Knight:
		destTypeBitboard = &(ourBitboardPtr.knights)
	case Rook:
		destTypeBitboard = &(ourBitboardPtr.rooks)
	case Bishop:
		destTypeBitboard = &(ourBitboardPtr.bishops)
	default:
		destTypeBitboard = pieceTypeBitboard
	}

	// Apply the move
	ourBitboardPtr.all &= ^fromBitboard // remove at "from"
	ourBitboardPtr.all |= toBitboard    // add at "to"
	*pieceTypeBitboard &= ^fromBitboard // remove at "from"
	*destTypeBitboard |= toBitboard    // add at "to"
	capturedPieceType, capturedBitboard := determinePieceType(oppBitboardPtr, toBitboard)
	if capturedPieceType != Nothing {
		*capturedBitboard &= ^toBitboard
		oppBitboardPtr.all &= ^toBitboard
	}
	b.wtomove = !b.wtomove

	// Return the unapply function (closure)
	unapply := func() {
		ourBitboardPtr.all &= ^toBitboard // remove at "to"
		ourBitboardPtr.all |= fromBitboard // add at "from"
		*destTypeBitboard &= ^toBitboard // remove at "to"
		*pieceTypeBitboard |= fromBitboard // add at "from"
		if capturedPieceType != Nothing {
			*capturedBitboard |= toBitboard
			oppBitboardPtr.all |= toBitboard
		}
		if castleStatus != 0 {
			ourBitboardPtr.rooks &= ^(uint64(1) << newRookLoc)
			ourBitboardPtr.rooks |= (uint64(1) << oldRookLoc)
		}
		if (epCaptureSquare != 0) {
			b.enpassant = epCaptureSquare
			oppBitboardPtr.pawns |= (uint64(1) << uint8(int8(epCaptureSquare) + epDelta))
			oppBitboardPtr.all |= (uint64(1) << uint8(int8(epCaptureSquare) + epDelta))
		}
		if b.wtomove {
			b.fullmoveno-- // decrement after undoing black's move
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
