package dragontoothmg

// Applies a move to the board, and returns a function that can be used to unapply it.
func (b *Board) Apply(m Move) func() {
	/*var ourBitboardPtr, oppBitboardPtr *bitboards
	if b.wtomove {
		ourBitboardPtr = &(b.white)
		oppBitboardPtr = &(b.black)
	} else {
		ourBitboardPtr = &(b.black)
		oppBitboardPtr = &(b.white)
	}
	ourBitboardPtr*/

	unapply := func() {

	}
	return unapply
}
