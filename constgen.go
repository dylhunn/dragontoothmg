package main

import (
	"fmt"
)

func main() {
	fmt.Printf("Rook: %#x\n\n", rookOccupancyMasks())
	fmt.Printf("Bishop: %#x\n", bishopOccupancyMasks())
}

// Generate the occupancy masks for a rook at each index
// This represents the locations the rook can slide to that don't block it
// Thus, the edges of the board are not included
// Used for magic bitboards
func rookOccupancyMasks() [64]uint64 {
	masks := [64]uint64{}
	for i := uint64(0); i < 64; i++ { // for a rook at index i
		// files: compute an index to activate
		for j := uint64((i % 8) + 8); j <= 55; j += 8 { 
			masks[i] |= 1 << j
		}
		// ranks: compute an index to activate
		for k := uint64((i / 8) * 8 + 1); k < ((i / 8) * 8 + 7); k ++ { 
			masks[i] |= 1 << k
		}
		masks[i] &= (^(uint64(1) << i)) // unset the origin bit
	}
	return masks
}

// Generate the occupancy masks for a bishop at each index
// This represents the locations the bishop can slide to that don't block it
// Thus, the edges of the board are not included
// Used for magic bitboards
func bishopOccupancyMasks() [64]uint64 {
	masks := [64]uint64{}
	for i := uint64(0); i < 64; i++ { // for a bishop at index i
		distanceFromAFile := i % 8 // distance to the right of the A file
		distanceFromRank1 := i / 8 // distance above the first rank
		distanceFromHFile := 7 - distanceFromAFile
		distanceFromRank8 := 7 - distanceFromRank1
		NErayLength := Min(distanceFromHFile, distanceFromRank8) - 1
		NWrayLength := Min(distanceFromAFile, distanceFromRank8) - 1
		SErayLength := Min(distanceFromHFile, distanceFromRank1) - 1
		SWrayLength := Min(distanceFromAFile, distanceFromRank1) - 1
		maxVal_uint64 := uint64(18446744073709551615)
		if (NErayLength == maxVal_uint64) { // A -1 might have overflowed
			NErayLength = 0
		}
		if (NWrayLength == maxVal_uint64) {
			NWrayLength = 0
		}
		if (SErayLength == maxVal_uint64) {
			SErayLength = 0
		}
		if (SWrayLength == maxVal_uint64) {
			SWrayLength = 0
		}
		for j := uint64(1); j <= NErayLength; j++ {
			masks[i] |= 1 << (j * 9 + i)
		}
		for j := uint64(1); j <= NWrayLength; j++ {
			masks[i] |= 1 << (j * 7 + i)
		}
		for j := uint64(1); j <= SErayLength; j++ {
			masks[i] |= 1 << (i - (j * 7))
		}
		for j := uint64(1); j <= SWrayLength; j++ {
			masks[i] |= 1 << (i - (j * 9))
		}
	}
	return masks
}

func Min(i uint64, j uint64) uint64 {
	if i < j {
		return i
	}
	return j
}