package dragontoothmg

// Applies a move to the board, and returns a function that can be used to unapply it.
// This function assumes that the given move is valid (i.e., is in the set of moves found by GenerateLegalMoves()).
// If the move is not valid, this function has undefined behavior.
func (b *Board) Apply(m Move) func() {
	// Configure data about which pieces move
	var ourBitboardPtr, oppBitboardPtr *bitboards
	var epDelta int8                                // add this to the e.p. square to find the captured pawn
	var oppStartingRankBb, ourStartingRankBb uint64 // the starting rank of out opponent's major pieces
	// the constant that represents the index into pieceSquareZobristC for the pawn of our color
	var ourPiecesPawnZobristIndex int
	var oppPiecesPawnZobristIndex int
	if b.wtomove {
		ourBitboardPtr = &(b.white)
		oppBitboardPtr = &(b.black)
		epDelta = -8
		oppStartingRankBb = onlyRank[7]
		ourStartingRankBb = onlyRank[0]
		ourPiecesPawnZobristIndex = 0
		oppPiecesPawnZobristIndex = 6
	} else {
		ourBitboardPtr = &(b.black)
		oppBitboardPtr = &(b.white)
		epDelta = 8
		oppStartingRankBb = onlyRank[0]
		ourStartingRankBb = onlyRank[7]
		b.fullmoveno++ // increment after black's move
		ourPiecesPawnZobristIndex = 6
		oppPiecesPawnZobristIndex = 0
	}
	fromBitboard := (uint64(1) << m.From())
	toBitboard := (uint64(1) << m.To())
	pieceType, pieceTypeBitboard := determinePieceType(ourBitboardPtr, fromBitboard)
	castleStatus := 0
	var oldRookLoc, newRookLoc uint8
	var flippedKsCastle, flippedQsCastle, flippedOppKsCastle, flippedOppQsCastle bool

	// King moves strip castling rights
	if pieceType == King {
		// TODO(dylhunn): do this without a branch
		if m.To()-m.From() == 2 { // castle short
			castleStatus = 1
			oldRookLoc = m.To() + 1
			newRookLoc = m.To() - 1
		} else if int(m.To())-int(m.From()) == -2 { // castle long
			castleStatus = -1
			oldRookLoc = m.To() - 2
			newRookLoc = m.To() + 1
		}
		// King moves always strip castling rights
		if b.canCastleKingside() {
			b.flipKingsideCastle()
			flippedKsCastle = true
		}
		if b.canCastleQueenside() {
			b.flipQueensideCastle()
			flippedQsCastle = true
		}
	}

	// Rook moves strip castling rights
	if pieceType == Rook {
		if b.canCastleKingside() && (fromBitboard&onlyFile[7] != 0) &&
			fromBitboard&ourStartingRankBb != 0 { // king's rook
			flippedKsCastle = true
			b.flipKingsideCastle()
		} else if b.canCastleQueenside() && (fromBitboard&onlyFile[0] != 0) &&
			fromBitboard&ourStartingRankBb != 0 { // queen's rook
			flippedQsCastle = true
			b.flipQueensideCastle()
		}
	}

	// Apply the castling rook movement
	if castleStatus != 0 {
		ourBitboardPtr.rooks |= (uint64(1) << newRookLoc)
		ourBitboardPtr.all |= (uint64(1) << newRookLoc)
		ourBitboardPtr.rooks &= ^(uint64(1) << oldRookLoc)
		ourBitboardPtr.all &= ^(uint64(1) << oldRookLoc)
		// Update rook location in hash
		// (Rook - 1) assumes that "Nothing" precedes "Rook" in the Piece constants list
		b.hash ^= pieceSquareZobristC[ourPiecesPawnZobristIndex+(Rook-1)][oldRookLoc]
		b.hash ^= pieceSquareZobristC[ourPiecesPawnZobristIndex+(Rook-1)][newRookLoc]
	}

	// Is this an e.p. capture? Strip the opponent pawn and reset the e.p. square
	oldEpCaptureSquare := b.enpassant
	var actuallyPerformedEpCapture bool = false
	if pieceType == Pawn && m.To() == oldEpCaptureSquare && oldEpCaptureSquare != 0 {
		actuallyPerformedEpCapture = true
		epOpponentPawnLocation := uint8(int8(oldEpCaptureSquare) + epDelta)
		oppBitboardPtr.pawns &= ^(uint64(1) << epOpponentPawnLocation)
		oppBitboardPtr.all &= ^(uint64(1) << epOpponentPawnLocation)
		// Remove the opponent pawn from the board hash.
		b.hash ^= pieceSquareZobristC[oppPiecesPawnZobristIndex][epOpponentPawnLocation]
	}
	// Update the en passant square
	if pieceType == Pawn && (int8(m.To())+2*epDelta == int8(m.From())) { // pawn double push
		b.enpassant = uint8(int8(m.To()) + epDelta)
	} else {
		b.enpassant = 0
	}

	// Is this a promotion?
	var destTypeBitboard *uint64
	var promotedToPieceType Piece // if not promoted, same as pieceType
	switch m.Promote() {
	case Queen:
		destTypeBitboard = &(ourBitboardPtr.queens)
		promotedToPieceType = Queen
	case Knight:
		destTypeBitboard = &(ourBitboardPtr.knights)
		promotedToPieceType = Knight
	case Rook:
		destTypeBitboard = &(ourBitboardPtr.rooks)
		promotedToPieceType = Rook
	case Bishop:
		destTypeBitboard = &(ourBitboardPtr.bishops)
		promotedToPieceType = Bishop
	default:
		destTypeBitboard = pieceTypeBitboard
		promotedToPieceType = pieceType
	}

	// Apply the move
	capturedPieceType, capturedBitboard := determinePieceType(oppBitboardPtr, toBitboard)
	ourBitboardPtr.all &= ^fromBitboard // remove at "from"
	ourBitboardPtr.all |= toBitboard    // add at "to"
	*pieceTypeBitboard &= ^fromBitboard // remove at "from"
	*destTypeBitboard |= toBitboard     // add at "to"
	if capturedPieceType != Nothing {   // This does not account for e.p. captures
		*capturedBitboard &= ^toBitboard
		oppBitboardPtr.all &= ^toBitboard
		b.hash ^= pieceSquareZobristC[oppPiecesPawnZobristIndex+(int(capturedPieceType)-1)][m.To()] // remove the captured piece from the hash
	}
	b.hash ^= pieceSquareZobristC[(int(pieceType)-1)+ourPiecesPawnZobristIndex][m.From()]         // remove piece at "from"
	b.hash ^= pieceSquareZobristC[(int(promotedToPieceType)-1)+ourPiecesPawnZobristIndex][m.To()] // add piece at "to"

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

	// flip the side to move in the hash
	b.hash ^= whiteToMoveZobristC
	b.wtomove = !b.wtomove

	// remove the old en passant square from the hash, and add the new one
	b.hash ^= uint64(oldEpCaptureSquare)
	b.hash ^= uint64(b.enpassant)

	// Return the unapply function (closure)
	unapply := func() {
		// Flip the player to move
		b.hash ^= whiteToMoveZobristC
		b.wtomove = !b.wtomove


		// Unapply move
		ourBitboardPtr.all &= ^toBitboard                                                             // remove at "to"
		ourBitboardPtr.all |= fromBitboard                                                            // add at "from"
		*destTypeBitboard &= ^toBitboard                                                              // remove at "to"
		*pieceTypeBitboard |= fromBitboard                                                            // add at "from"
		b.hash ^= pieceSquareZobristC[(int(promotedToPieceType)-1)+ourPiecesPawnZobristIndex][m.To()] // remove the piece at "to"
		b.hash ^= pieceSquareZobristC[(int(pieceType)-1)+ourPiecesPawnZobristIndex][m.From()]         // add the piece at "from"
		
		// Restore captured piece (excluding e.p.)
		if capturedPieceType != Nothing {                                                             // doesn't consider e.p. captures
			*capturedBitboard |= toBitboard
			oppBitboardPtr.all |= toBitboard
			// restore the captured piece to the hash (excluding e.p.)
			b.hash ^= pieceSquareZobristC[oppPiecesPawnZobristIndex+(int(capturedPieceType)-1)][m.To()]
		}

		// Restore rooks from castling move
		if castleStatus != 0 {
			ourBitboardPtr.rooks &= ^(uint64(1) << newRookLoc)
			ourBitboardPtr.all &= ^(uint64(1) << newRookLoc)
			ourBitboardPtr.rooks |= (uint64(1) << oldRookLoc)
			ourBitboardPtr.all |= (uint64(1) << oldRookLoc)
			// Revert castling rook move
			b.hash ^= pieceSquareZobristC[ourPiecesPawnZobristIndex+(Rook-1)][oldRookLoc]
			b.hash ^= pieceSquareZobristC[ourPiecesPawnZobristIndex+(Rook-1)][newRookLoc]
		}

		// Unapply en-passant square change, and capture if necessary
		b.hash ^= uint64(b.enpassant)     // undo the new en passant square from the hash
		b.hash ^= uint64(oldEpCaptureSquare) // restore the old one to the hash
		b.enpassant = oldEpCaptureSquare
		if actuallyPerformedEpCapture {
			epOpponentPawnLocation := uint8(int8(oldEpCaptureSquare) + epDelta)
			oppBitboardPtr.pawns |= (uint64(1) << epOpponentPawnLocation)
			oppBitboardPtr.all |= (uint64(1) << epOpponentPawnLocation)
			// Add the opponent pawn to the board hash.
			b.hash ^= pieceSquareZobristC[oppPiecesPawnZobristIndex][epOpponentPawnLocation]
		}

		// Decrement move clock
		if !b.wtomove {
			b.fullmoveno-- // decrement after undoing black's move
		}
		
		// Restore castling flags
		// Must update castling flags AFTER turn swap
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
