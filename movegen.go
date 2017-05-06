// dragontoothmg is a fast chess legal move generator library based on magic bitboards.
package dragontoothmg

// The main Dragontooth move generator file.
// Functions are in this file if (and only if) they are performance-critical
// move generator components, called while actually generating moves in-game.
// (The exception is a few one-line helpers for Move and Board in types.go)

import (
	"math/bits"
	//"fmt"
)

// The main API entrypoint. Generates all legal moves for a given board.
func (b *Board) GenerateLegalMoves() []Move {
	moves := make([]Move, 0, 45)
	// First, see if we are currently in check. If we are, invoke a special check-
	// evasion move generator.

	// Then, calculate all the absolutely pinned pieces, and compute their moves.
	pinnedPieces := b.generatePinnedMoves(&moves)
	nonpinnedPieces := ^pinnedPieces

	// Finally, compute ordinary moves, ignoring absolutely pinned pieces.
	b.pawnPushes(&moves, nonpinnedPieces)
	b.pawnCaptures(&moves, nonpinnedPieces)
	b.knightMoves(&moves, nonpinnedPieces)
	b.kingMoves(&moves)
	b.rookMoves(&moves, nonpinnedPieces)
	b.bishopMoves(&moves, nonpinnedPieces)
	b.queenMoves(&moves, nonpinnedPieces)
	return moves
}

// Calculate the available moves for absolutely pinned pieces (pinned to the king).
// Return a bitboard of all pieces that are pinned.
func (b *Board) generatePinnedMoves(moveList *[]Move) uint64 {
	var ourKingIdx uint8
	var ourPieces, oppPieces *bitboards
	var allPinnedPieces uint64 = 0
	if b.wtomove { // Assumes only one king on the board
		ourKingIdx = uint8(bits.TrailingZeros64(b.white.kings))
		ourPieces = &(b.white)
		oppPieces = &(b.black)
	} else {
		ourKingIdx = uint8(bits.TrailingZeros64(b.black.kings))
		ourPieces = &(b.black)
		oppPieces = &(b.white)
	}
	allPieces := oppPieces.all | ourPieces.all

	// Calculate king moves as if it was a rook.
	// "king targets" includes our own friendly pieces, for the purpose of identifying pins.
	kingOrthoTargets := calculateRookMoveBitboard(ourKingIdx, allPieces)
	oppRooks := oppPieces.rooks | oppPieces.queens
	for oppRooks != 0 { // For each opponent ortho slider
		currRookIdx := uint8(bits.TrailingZeros64(oppRooks))
		oppRooks &= oppRooks - 1
		rookTargets := calculateRookMoveBitboard(currRookIdx, allPieces) & (^(oppPieces.all))
		pinnedPiece := rookTargets & kingOrthoTargets // A piece is pinned iff it falls along both attack rays.
		if pinnedPiece == 0 {                         // there is no pin
			continue
		}
		pinnedPieceIdx := uint8(bits.TrailingZeros64(pinnedPiece))
		sameRank := pinnedPieceIdx/8 == ourKingIdx/8 && pinnedPieceIdx/8 == currRookIdx/8
		sameFile := pinnedPieceIdx%8 == ourKingIdx%8 && pinnedPieceIdx%8 == currRookIdx%8
		if !sameRank && !sameFile {
			continue // it's just an intersection, not a pin
		}
		allPinnedPieces |= pinnedPiece        // store the pinned piece location
		if pinnedPiece&ourPieces.pawns != 0 { // it's a pawn; we might be able to push it
			if sameFile { // push the pawn
				pawnPushesSingle, pawnPushesDouble := b.pawnPushBitboards(everything)
				pawnTargets := (pawnPushesSingle | pawnPushesDouble) & onlyFile[pinnedPieceIdx%8]
				genMovesFromTargets(moveList, Square(pinnedPieceIdx), pawnTargets)
			}
			continue
		}
		// If it's not a rook or queen, it can't move
		if pinnedPiece&ourPieces.rooks == 0 && pinnedPiece&ourPieces.queens == 0 {
			continue
		}
		// all ortho moves, as if it was not pinned
		pinnedPieceAllMoves := calculateRookMoveBitboard(pinnedPieceIdx, allPieces) & (^(ourPieces.all))
		// actually available moves
		pinnedTargets := pinnedPieceAllMoves & (rookTargets | kingOrthoTargets | (uint64(1) << currRookIdx))
		genMovesFromTargets(moveList, Square(pinnedPieceIdx), pinnedTargets)
	}

	// Calculate king moves as if it was a bishop.
	// "king targets" includes our own friendly pieces, for the purpose of identifying pins.
	kingDiagTargets := calculateBishopMoveBitboard(ourKingIdx, allPieces)
	oppBishops := oppPieces.bishops | oppPieces.queens
	for oppBishops != 0 {
		currBishopIdx := uint8(bits.TrailingZeros64(oppBishops))
		oppBishops &= oppBishops - 1
		bishopTargets := calculateBishopMoveBitboard(currBishopIdx, allPieces) & (^(oppPieces.all))
		pinnedPiece := bishopTargets & kingDiagTargets // A piece is pinned iff it falls along both attack rays.
		if pinnedPiece == 0 {                          // there is no pin
			continue
		}
		pinnedPieceIdx := uint8(bits.TrailingZeros64(pinnedPiece))
		bishopToPinnedSlope := (float32(pinnedPieceIdx)/8 - float32(currBishopIdx)/8) /
			(float32(pinnedPieceIdx%8) - float32(currBishopIdx%8))
		bishopToKingSlope := (float32(ourKingIdx)/8 - float32(currBishopIdx)/8) /
			(float32(ourKingIdx%8) - float32(currBishopIdx%8))
		if bishopToPinnedSlope != bishopToKingSlope { // just an intersection, not a pin
			continue
		}

		//fmt.Println(pinnedPieceIdx/8, (currBishopIdx/8)+1)

		allPinnedPieces |= pinnedPiece        // store pinned piece
		if pinnedPiece&ourPieces.pawns != 0 { // it's a pawn; we might be able to capture with it
			if (b.wtomove && (pinnedPieceIdx/8)+1 == currBishopIdx/8) ||
				(!b.wtomove && pinnedPieceIdx/8 == (currBishopIdx/8)+1) {

				var move Move
				move.Setfrom(Square(pinnedPieceIdx)).Setto(Square(currBishopIdx))
				*moveList = append(*moveList, move)
			}
			continue
		}
		// If it's not a bishop or queen, it can't move
		if pinnedPiece&ourPieces.bishops == 0 && pinnedPiece&ourPieces.queens == 0 {
			continue
		}
		// all diag moves, as if it was not pinned
		pinnedPieceAllMoves := calculateBishopMoveBitboard(pinnedPieceIdx, allPieces) & (^(ourPieces.all))
		// actually available moves
		pinnedTargets := pinnedPieceAllMoves & (bishopTargets | kingDiagTargets | (uint64(1) << currBishopIdx))
		genMovesFromTargets(moveList, Square(pinnedPieceIdx), pinnedTargets)
	}
	return allPinnedPieces
}

// Generate moves involving advancing pawns.
func (b *Board) pawnPushes(moveList *[]Move, nonpinned uint64) {
	targets, doubleTargets := b.pawnPushBitboards(nonpinned)
	oneRankBack := 8
	if b.wtomove {
		oneRankBack = -oneRankBack
	}
	// push all pawns by one square
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1 // unset the lowest active bit
		var canPromote bool
		if b.wtomove {
			canPromote = target >= 56
		} else {
			canPromote = target <= 7
		}
		var move Move
		move.Setfrom(Square(target + oneRankBack)).Setto(Square(target))
		if canPromote {
			for i := Piece(Knight); i <= Queen; i++ {
				move.Setpromote(i)
				*moveList = append(*moveList, move)
			}
		} else {
			*moveList = append(*moveList, move)
		}
	}
	// push some pawns by two squares
	for doubleTargets != 0 {
		doubleTarget := bits.TrailingZeros64(doubleTargets)
		doubleTargets &= doubleTargets - 1 // unset the lowest active bit
		var move Move
		move.Setfrom(Square(doubleTarget + 2*oneRankBack)).Setto(Square(doubleTarget))
		*moveList = append(*moveList, move)
	}
}

// A helper function that produces bitboards of valid pawn push locations.
func (b *Board) pawnPushBitboards(nonpinned uint64) (targets uint64, doubleTargets uint64) {
	free := (^b.white.all) & (^b.black.all)
	if b.wtomove {
		movableWhitePawns := b.white.pawns & nonpinned
		targets = movableWhitePawns << 8 & free
		fourthFile := uint64(0xFF000000)
		doubleTargets = targets << 8 & fourthFile & free
	} else {
		movableBlackPawns := b.black.pawns & nonpinned
		targets = movableBlackPawns >> 8 & free
		fifthFile := uint64(0xFF00000000)
		doubleTargets = targets >> 8 & fifthFile & free
	}
	return
}

// A function that computes available pawn captures.
func (b *Board) pawnCaptures(moveList *[]Move, nonpinned uint64) {
	east, west := b.pawnCaptureBitboards(nonpinned)
	bitboards := [2]uint64{east, west}
	if !b.wtomove {
		bitboards[0], bitboards[1] = bitboards[1], bitboards[0]
	}
	for dir, board := range bitboards { // for east and west
		for board != 0 {
			target := bits.TrailingZeros64(board)
			board &= board - 1
			var move Move
			move.Setto(Square(target))
			canPromote := false
			if b.wtomove {
				move.Setfrom(Square(target - (9 - (dir * 2))))
				canPromote = target >= 56
			} else {
				move.Setfrom(Square(target + (9 - (dir * 2))))
				canPromote = target <= 7
			}
			if canPromote {
				for i := Piece(Knight); i <= Queen; i++ {
					move.Setpromote(i)
					*moveList = append(*moveList, move)
				}
				continue
			}
			*moveList = append(*moveList, move)
		}
	}
}

// A helper than generates bitboards for available pawn captures.
func (b *Board) pawnCaptureBitboards(nonpinned uint64) (east uint64, west uint64) {
	notHFile := uint64(0x7F7F7F7F7F7F7F7F)
	notAFile := uint64(0xFEFEFEFEFEFEFEFE)
	var targets uint64
	if b.enpassant > 0 { // an en-passant target is active
		targets = (1 << b.enpassant)
	}
	if b.wtomove {
		targets |= b.black.all
		ourpawns := b.white.pawns & nonpinned
		east = ourpawns << 9 & notAFile & targets
		west = ourpawns << 7 & notHFile & targets
	} else {
		targets |= b.white.all
		ourpawns := b.black.pawns & nonpinned
		east = ourpawns >> 7 & notAFile & targets
		west = ourpawns >> 9 & notHFile & targets
	}
	return
}

// Generate all knight moves.
func (b *Board) knightMoves(moveList *[]Move, nonpinned uint64) {
	var ourKnights, noFriendlyPieces uint64
	if b.wtomove {
		ourKnights = b.white.knights & nonpinned
		noFriendlyPieces = (^b.white.all)
	} else {
		ourKnights = b.black.knights & nonpinned
		noFriendlyPieces = (^b.black.all)
	}
	for ourKnights != 0 {
		currentKnight := bits.TrailingZeros64(ourKnights)
		ourKnights &= ourKnights - 1
		targets := knightMasks[currentKnight] & noFriendlyPieces
		genMovesFromTargets(moveList, Square(currentKnight), targets)
	}
}

// Generate all available king moves.
// First, if castling is possible, verifies the checking prohibitions on castling.
// Then, outputs castling moves (if any), and king moves.
// Not thread-safe, since the king is removed from the board to compute
// king-danger squares.
func (b *Board) kingMoves(moveList *[]Move) {
	var ourKingLocation uint8
	var noFriendlyPieces uint64
	var canCastleQueenside, canCastleKingside bool
	var ptrToOurBitboards *bitboards
	allPieces := b.white.all | b.black.all
	if b.wtomove {
		ourKingLocation = uint8(bits.TrailingZeros64(b.white.kings))
		ptrToOurBitboards = &(b.white)
		noFriendlyPieces = ^(b.white.all)
		// To castle, we must have rights and a clear path
		kingsideClear := allPieces&((1<<5)|(1<<6)) == 0
		queensideClear := allPieces&((1<<3)|(1<<2)|(1<<1)) == 0
		// skip the king square, since this won't be called while in check
		canCastleQueenside = b.whiteCanCastleQueenside() &&
			queensideClear && !b.anyUnderDirectAttack(true, 0, 1, 2, 3)
		canCastleKingside = b.whiteCanCastleKingside() &&
			kingsideClear && !b.anyUnderDirectAttack(true, 5, 6, 7)
	} else {
		ourKingLocation = uint8(bits.TrailingZeros64(b.black.kings))
		ptrToOurBitboards = &(b.black)
		noFriendlyPieces = ^(b.black.all)
		kingsideClear := allPieces&((1<<61)|(1<<62)) == 0
		queensideClear := allPieces&((1<<57)|(1<<58)|(1<<59)) == 0
		// skip the king square, since this won't be called while in check
		canCastleQueenside = b.blackCanCastleQueenside() &&
			queensideClear && !b.anyUnderDirectAttack(false, 56, 57, 58, 59)
		canCastleKingside = b.blackCanCastleKingside() &&
			kingsideClear && !b.anyUnderDirectAttack(false, 61, 62, 63)
	}
	if canCastleKingside {
		var move Move
		move.Setfrom(Square(ourKingLocation)).Setto(Square(ourKingLocation + 2))
		*moveList = append(*moveList, move)
	}
	if canCastleQueenside {
		var move Move
		move.Setfrom(Square(ourKingLocation)).Setto(Square(ourKingLocation - 2))
		*moveList = append(*moveList, move)
	}

	// TODO(dylhunn): Modifying the board is NOT thread-safe.
	// We only do this to avoid the king danger problem, aka moving away from a
	// checking slider.
	oldKings := ptrToOurBitboards.kings
	ptrToOurBitboards.kings = 0
	ptrToOurBitboards.all &= ^(1 << ourKingLocation)

	targets := kingMasks[ourKingLocation] & noFriendlyPieces
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1
		if b.underDirectAttack(b.wtomove, uint8(target)) {
			continue
		}
		var move Move
		move.Setfrom(Square(ourKingLocation)).Setto(Square(target))
		*moveList = append(*moveList, move)
	}

	ptrToOurBitboards.kings = oldKings
	ptrToOurBitboards.all |= (1 << ourKingLocation)
}

// Generate all rook moves using magic bitboards.
func (b *Board) rookMoves(moveList *[]Move, nonpinned uint64) {
	var ourRooks, friendlyPieces uint64
	if b.wtomove {
		ourRooks = b.white.rooks & nonpinned
		friendlyPieces = b.white.all
	} else {
		ourRooks = b.black.rooks & nonpinned
		friendlyPieces = b.black.all
	}
	allPieces := b.white.all | b.black.all
	for ourRooks != 0 {
		currRook := uint8(bits.TrailingZeros64(ourRooks))
		ourRooks &= ourRooks - 1
		targets := calculateRookMoveBitboard(currRook, allPieces) & (^friendlyPieces)
		genMovesFromTargets(moveList, Square(currRook), targets)
	}
}

// Generate all bishop moves using magic bitboards.
func (b *Board) bishopMoves(moveList *[]Move, nonpinned uint64) {
	var ourBishops, friendlyPieces uint64
	if b.wtomove {
		ourBishops = b.white.bishops & nonpinned
		friendlyPieces = b.white.all
	} else {
		ourBishops = b.black.bishops & nonpinned
		friendlyPieces = b.black.all
	}
	allPieces := b.white.all | b.black.all
	for ourBishops != 0 {
		currBishop := uint8(bits.TrailingZeros64(ourBishops))
		ourBishops &= ourBishops - 1
		targets := calculateBishopMoveBitboard(currBishop, allPieces) & (^friendlyPieces)
		genMovesFromTargets(moveList, Square(currBishop), targets)
	}
}

// Generate all queen moves using magic bitboards.
func (b *Board) queenMoves(moveList *[]Move, nonpinned uint64) {
	var ourQueens, friendlyPieces uint64
	if b.wtomove {
		ourQueens = b.white.queens & nonpinned
		friendlyPieces = b.white.all
	} else {
		ourQueens = b.black.queens & nonpinned
		friendlyPieces = b.black.all
	}
	allPieces := b.white.all | b.black.all
	for ourQueens != 0 {
		currQueen := uint8(bits.TrailingZeros64(ourQueens))
		ourQueens &= ourQueens - 1
		// bishop motion
		diag_targets := calculateBishopMoveBitboard(currQueen, allPieces) & (^friendlyPieces)
		genMovesFromTargets(moveList, Square(currQueen), diag_targets)
		// rook motion
		ortho_targets := calculateRookMoveBitboard(currQueen, allPieces) & (^friendlyPieces)
		genMovesFromTargets(moveList, Square(currQueen), ortho_targets)
	}
}

// Helper: converts a targets bitboard into moves, and adds them to the moves list.
func genMovesFromTargets(moveList *[]Move, origin Square, targets uint64) {
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1
		var move Move
		move.Setfrom(origin).Setto(Square(target))
		*moveList = append(*moveList, move)
	}
}

// Variadic function that returns whether any of the specified squares is being attacked
// by the opponent. Potentially expensive.
func (b *Board) anyUnderDirectAttack(byBlack bool, squares ...uint8) bool {
	for _, v := range squares {
		if b.underDirectAttack(byBlack, v) {
			return true
		}
	}
	return false
}

// Compute whether an individual square is under direct attack. Potentially expensive.
func (b *Board) underDirectAttack(byBlack bool, origin uint8) bool {
	allPieces := b.white.all | b.black.all
	var opponentPieces *bitboards
	if byBlack {
		opponentPieces = &(b.black)
	} else {
		opponentPieces = &(b.white)
	}
	// find attacking knights
	knight_attackers := knightMasks[origin] & opponentPieces.knights
	if knight_attackers != 0 {
		return true
	}
	// find attacking bishops and queens
	diag_candidates := magicBishopBlockerMasks[origin] & allPieces
	diag_dbindex := (diag_candidates * magicNumberBishop[origin]) >> magicBishopShifts[origin]
	diag_potential_attackers := magicMovesBishop[origin][diag_dbindex] & opponentPieces.all
	diag_attackers := diag_potential_attackers & (opponentPieces.bishops | opponentPieces.queens)
	if diag_attackers != 0 {
		return true
	}
	// find attacking rooks and queens
	ortho_candidates := magicRookBlockerMasks[origin] & allPieces
	ortho_dbindex := (ortho_candidates * magicNumberRook[origin]) >> magicRookShifts[origin]
	ortho_potential_attackers := magicMovesRook[origin][ortho_dbindex] & opponentPieces.all
	ortho_attackers := ortho_potential_attackers & (opponentPieces.rooks | opponentPieces.queens)
	if ortho_attackers != 0 {
		return true
	}
	// find attacking kings
	// TODO(dylhunn): What if the opponent king can't actually move to the origin square?
	king_attackers := kingMasks[origin] & opponentPieces.kings
	if king_attackers != 0 {
		return true
	}
	// find attacking pawns
	var pawn_attackers uint64 = 0
	if byBlack {
		pawn_attackers = 1 << (origin + 7)
		pawn_attackers |= 1 << (origin + 9)
	} else {
		if origin-7 >= 0 {
			pawn_attackers = 1 << (origin - 7)
		}
		if origin-9 >= 0 {
			pawn_attackers |= 1 << (origin - 9)
		}
	}
	pawn_attackers &= opponentPieces.pawns
	if pawn_attackers != 0 {
		return true
	}
	return false
}

// Calculates the attack bitboard for a rook. This might include targeted squares
// that are actually friendly pieces, so the proper usage is:
// rookTargets := calculateRookMoveBitboard(myRookLoc, allPieces) & (^myPieces)
func calculateRookMoveBitboard(currRook uint8, allPieces uint64) uint64 {
	blockers := magicRookBlockerMasks[currRook] & allPieces
	dbindex := (blockers * magicNumberRook[currRook]) >> magicRookShifts[currRook]
	targets := magicMovesRook[currRook][dbindex]
	return targets
}

// Calculates the attack bitboard for a bishop. This might include targeted squares
// that are actually friendly pieces, so the proper usage is:
// bishopTargets := calculateBishopMoveBitboard(myBishopLoc, allPieces) & (^myPieces)
func calculateBishopMoveBitboard(currBishop uint8, allPieces uint64) uint64 {
	blockers := magicBishopBlockerMasks[currBishop] & allPieces
	dbindex := (blockers * magicNumberBishop[currBishop]) >> magicBishopShifts[currBishop]
	targets := magicMovesBishop[currBishop][dbindex]
	return targets
}
