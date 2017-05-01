[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
![Build Status](http://img.shields.io/travis/dylhunn/dragontooth-movegen.svg)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-green.svg)](https://godoc.org/github.com/dylhunn/dragontooth-movegen)


Dragontooth Movegen | Dylan D. Hunn
==================================

Dragontooth Movegen is a fast, magic-bitboard chess move generator written entirely in Go. It provides a simple API for generating pseudo-legal moves.

Repo summary
============

Here is a summary of the important files in the repo:

| **File**         | **Description**                                                                                                                                         |
|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| movegen.go   | This is the "primary" source file. Functions are located here if, and only if, they are performance critical and executed to generate moves in-game. |
| types.go     | This file contains the Board and Moves types, along with some supporting helper functions and types.                                                 |
| constants.go | All constants for move generation are hard-coded here, along with functions to compute the magic bitboard lookup tables when the file loads.         |
| util.go      | This file contains supporting library functions, for FEN reading and conversions.                                                                    |

**This project is currently a prototype.** It works, although optimizations are underway, along with a legal move generation feature.

Installing and building the library
===================================

This project requires Go 1.9. As of the time of writing, 1.9 is still a pre-release version. You can get it by cloning the official [Go Repo](https://github.com/golang/go), and building it yourself.

To use this package in your own code, make sure your `GO_PATH` environment variable is correctly set, and install it using `go get`:

    go get github.com/dylhunn/dragontooth-movegen

Then, you can include it in your project:

	import "github.com/dylhunn/dragontooth-movegen"

Alternatively, you can clone it yourself:

    git clone https://github.com/dylhunn/dragontooth-movegen.git


Documentation and examples
==========================

You can find the documentation [here](https://godoc.org/github.com/dylhunn/dragontooth-movegen).

Here is a simple example invocation:

    board := movegen.ParseFen("1Q2rk2/2p2p2/1n4b1/N7/2B1Pp1q/2B4P/1QPP4/4K2R b K e3 4 30")
    moveList := board.GeneratePseudolegalMoves()
    for _, curr := range moveList {
    	fmt.Println("Moved to: %v", dragontoothmg.IndexToAlgebraic(curr.To()))
    }