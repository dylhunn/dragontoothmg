package dragontoothmg

// Applies a move to the board, and returns a function that can be used to unapply it.
func (b *Board) Apply(m Move) func() {
	// Configure data about which pieces move
	var ourBitboardPtr, oppBitboardPtr *bitboards
	var epDelta int8                                // add this to the e.p. square to find the captured pawn
	var oppStartingRankBb, ourStartingRankBb uint64 // the starting rank of out opponent's major pieces
	if b.wtomove {
		ourBitboardPtr = &(b.white)
		oppBitboardPtr = &(b.black)
		epDelta = -8
		oppStartingRankBb = onlyRank[7]
		ourStartingRankBb = onlyRank[0]
	} else {
		ourBitboardPtr = &(b.black)
		oppBitboardPtr = &(b.white)
		epDelta = 8
		oppStartingRankBb = onlyRank[0]
		ourStartingRankBb = onlyRank[7]
		b.fullmoveno++ // increment after black's move
	}
	fromBitboard := (uint64(1) << m.From())
	toBitboard := (uint64(1) << m.To())
	pieceType, pieceTypeBitboard := determinePieceType(ourBitboardPtr, fromBitboard)
	castleStatus := 0
	var oldRookLoc, newRookLoc uint8
	kingsideCastleRightsBefore := b.canCastleKingside()
	queensideCastleRightsBefore := b.canCastleQueenside()
	var flippedKsCastle, flippedQsCastle, flippedOppKsCastle, flippedOppQsCastle bool

	// Configure castling rights
	if pieceType == King {
		if m.To()-m.From() == 2 { // castle short
			castleStatus = 1
			oldRookLoc = m.To() + 1
			newRookLoc = m.To() - 1
			b.flipKingsideCastle()
			if queensideCastleRightsBefore {
				b.flipQueensideCastle()
				flippedQsCastle = true
			}
			flippedKsCastle = true
		} else if int(m.To())-int(m.From()) == -2 { // castle long
			castleStatus = -1
			oldRookLoc = m.To() - 2
			newRookLoc = m.To() + 1
			b.flipQueensideCastle()
			if kingsideCastleRightsBefore {
				b.flipKingsideCastle()
				flippedKsCastle = true
			}
			flippedQsCastle = true
		} else { // an ordinary non-castling king movement
			if kingsideCastleRightsBefore {
				b.flipKingsideCastle()
				flippedKsCastle = true
			}
			if queensideCastleRightsBefore {
				b.flipQueensideCastle()
				flippedQsCastle = true
			}
		}
	}
	// Apply the castling rook movement
	if castleStatus != 0 {
		ourBitboardPtr.rooks |= (uint64(1) << newRookLoc)
		ourBitboardPtr.all |= (uint64(1) << newRookLoc)
		ourBitboardPtr.rooks &= ^(uint64(1) << oldRookLoc)
		ourBitboardPtr.all &= ^(uint64(1) << oldRookLoc)
	}

	// Rook moves strip castling rights
	if pieceType == Rook {
		if b.canCastleKingside() && (fromBitboard&onlyFile[7] != 0) && fromBitboard&ourStartingRankBb != 0 { // king's rook
			flippedKsCastle = true
			b.flipKingsideCastle()
		} else if b.canCastleQueenside() && (fromBitboard&onlyFile[0] != 0) && fromBitboard&ourStartingRankBb != 0 { // queen's rook
			flippedQsCastle = true
			b.flipQueensideCastle()
		}
	}

	// remove the old en passant square from the hash
	b.hash ^= uint64(b.enpassant)
	// Is this an e.p. capture? Strip the opponent pawn and reset the e.p. square
	epCaptureSquare := b.enpassant
	if pieceType == Pawn && m.To() == epCaptureSquare && epCaptureSquare != 0 {
		oppBitboardPtr.pawns &= ^(uint64(1) << uint8(int8(epCaptureSquare)+epDelta))
		oppBitboardPtr.all &= ^(uint64(1) << uint8(int8(epCaptureSquare)+epDelta))
	}
	// Update the en passant square
	if pieceType == Pawn && (int8(m.To())+2*epDelta == int8(m.From())) { // pawn double push
		b.enpassant = uint8(int8(m.To()) + epDelta)
	} else {
		b.enpassant = 0
	}
	// add the new en passant square to the hash
	b.hash ^= uint64(b.enpassant)

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
	capturedPieceType, capturedBitboard := determinePieceType(oppBitboardPtr, toBitboard)
	ourBitboardPtr.all &= ^fromBitboard // remove at "from"
	ourBitboardPtr.all |= toBitboard    // add at "to"
	*pieceTypeBitboard &= ^fromBitboard // remove at "from"
	*destTypeBitboard |= toBitboard     // add at "to"
	if capturedPieceType != Nothing {
		*capturedBitboard &= ^toBitboard
		oppBitboardPtr.all &= ^toBitboard
	}

	// If a rook was captured, it strips castling rights
	if capturedPieceType == Rook {
		if b.oppCanCastleKingside() && m.To()%8 == 7 && toBitboard&oppStartingRankBb != 0 { // captured king rook
			b.flipOppKingsideCastle()
			flippedOppKsCastle = true
		} else if b.oppCanCastleQueenside() && m.To()%8 == 0 && toBitboard&oppStartingRankBb != 0 { // queen rooks
			b.flipOppQueensideCastle()
			flippedOppQsCastle = true
		}
	}

	b.hash ^= whiteToMoveZobristC
	b.wtomove = !b.wtomove

	// Return the unapply function (closure)
	unapply := func() {
		ourBitboardPtr.all &= ^toBitboard  // remove at "to"
		ourBitboardPtr.all |= fromBitboard // add at "from"
		*destTypeBitboard &= ^toBitboard   // remove at "to"
		*pieceTypeBitboard |= fromBitboard // add at "from"
		if capturedPieceType != Nothing {
			*capturedBitboard |= toBitboard
			oppBitboardPtr.all |= toBitboard
		}
		if castleStatus != 0 {
			ourBitboardPtr.rooks &= ^(uint64(1) << newRookLoc)
			ourBitboardPtr.all &= ^(uint64(1) << newRookLoc)
			ourBitboardPtr.rooks |= (uint64(1) << oldRookLoc)
			ourBitboardPtr.all |= (uint64(1) << oldRookLoc)
		}
		// undo the new en passant square from the hash
		b.hash ^= uint64(b.enpassant)
		b.enpassant = epCaptureSquare
		b.hash ^= uint64(b.enpassant) // restore the old one
		if epCaptureSquare != 0 {
			oppBitboardPtr.pawns |= (uint64(1) << uint8(int8(epCaptureSquare)+epDelta))
			oppBitboardPtr.all |= (uint64(1) << uint8(int8(epCaptureSquare)+epDelta))
		}
		if b.wtomove {
			b.fullmoveno-- // decrement after undoing black's move
		}

		b.hash ^= whiteToMoveZobristC
		b.wtomove = !b.wtomove

		// must update castling flags after turn swap
		if flippedKsCastle {
			b.flipKingsideCastle()
		}
		if flippedQsCastle {
			b.flipQueensideCastle()
		}
		if flippedOppKsCastle {
			b.flipOppKingsideCastle()
		}
		if flippedOppQsCastle {
			b.flipOppQueensideCastle()
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
