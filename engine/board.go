package engine

// array of all pieces on a given board
type Board struct {
	Board    []*Piece // all of the pieces on the board
	Lastmove Move
	Turn     int // 1 : white , -1 : black
}

// Checks if a king is in check.
// Pass the color of the king that you want to check.
// Returns true if king in check / false if not.
func (b *Board) IsCheck(color int) bool {
	var kingsquare Square
	if color == 1 {
		kingsquare = b.Board[0].Position
	} else if color == -1 {
		kingsquare = b.Board[1].Position
	}
	for _, piece := range b.Board {
		if piece.Color == color*-1 {
			for _, move := range piece.legalMoves(b, false) {
				if move.End == kingsquare {
					return true
				}
			}
		}
	}
	return false
}

// Returns all legal moves available to the player whose turn it is.
func (b *Board) AllLegalMoves() []*Move {
	legals := make([]*Move, 0)
	for _, p := range b.Board {
		if p.Color == b.Turn {
			for _, m := range p.legalMoves(b, true) {
				legals = append(legals, m)
			}
		}
	}
	return legals
}

// Returns a deep copy of a given board
func (b *Board) CopyBoard() *Board {
	newboard := new(Board)
	s := make([]*Piece, len(b.Board))
	newboard.Lastmove, newboard.Turn = b.Lastmove, b.Turn
	for i, _ := range b.Board {
		p := new(Piece)
		*p = *b.Board[i]
		s[i] = p
	}
	newboard.Board = s
	return newboard
}

// Returns a slice of pointers to all possible boards
func (b *Board) NewGen() []*Board {
	legals := b.AllLegalMoves()
	s := make([]*Board, len(legals))
	for i, m := range legals {
		newboard := b.CopyBoard()
		newboard.Move(m)
		s[i] = newboard
	}
	return s
}

// Checks if the game has ended.
// Returns 2 if white wins, -2 if black wins, 1 if it's stalemate, 0 if the game is still going.
func (b *Board) IsOver() int {
	var kingindex int
	if b.Turn == 1 {
		kingindex = 0
	} else if b.Turn == -1 {
		kingindex = 1
	}
	if len(b.Board[kingindex].legalMoves(b, true)) == 0 {
		if b.IsCheck(b.Turn) {
			return -2 * b.Turn
		}
		if len(b.AllLegalMoves()) == 0 {
			return 1
		}
	}
	return 0
}

// Given a name, color, and coordinates, place the appropriate piece on the board.
// Does not add flags such as Can_Castle, must be done manually.
func (b *Board) PlacePiece(name byte, color, x, y int) {
	p := &Piece{
		Name:  name,
		Color: color,
		Position: Square{
			X: x,
			Y: y,
		},
	}
	if name == 'b' || name == 'r' || name == 'q' {
		p.Infinite_direction = true
	}
	if name == 'p' {
		p.Directions = [][2]int{
			{0, 1 * color},
		}
	} else if name == 'b' {
		p.Directions = [][2]int{
			{1, 1},
			{1, -1},
			{-1, 1},
			{-1, -1},
		}
	} else if name == 'n' {
		p.Directions = [][2]int{
			{1, 2},
			{-1, 2},
			{1, -2},
			{-1, -2},
			{2, 1},
			{-2, 1},
			{2, -1},
			{-2, -1},
		}

	} else if name == 'r' {
		p.Directions = [][2]int{
			{1, 0},
			{-1, 0},
			{0, 1},
			{0, -1},
		}
	} else if name == 'q' {
		p.Directions = [][2]int{
			{1, 1},
			{1, 0},
			{1, -1},
			{0, 1},
			{0, -1},
			{-1, 1},
			{-1, 0},
			{-1, -1},
		}
	} else if name == 'k' {
		p.Directions = [][2]int{
			{1, 1},
			{1, 0},
			{1, -1},
			{0, 1},
			{0, -1},
			{-1, 1},
			{-1, 0},
			{-1, -1},
		}
	}
	b.Board = append(b.Board, p)
}

// Resets a given board to its starting position.
func (b *Board) SetUpPieces() {
	b.Board = make([]*Piece, 0)
	pawnrows := [2]int{2, 7}
	piecerows := [2]int{1, 8}
	rookfiles := [2]int{1, 8}
	knightfiles := [2]int{2, 7}
	bishopfiles := [2]int{3, 6}
	queenfile := 4
	kingfile := 5
	for _, rank := range piecerows {
		// put the kings in first
		var color int
		if rank == 1 {
			color = 1
		} else {
			color = -1
		}
		b.PlacePiece('k', color, kingfile, rank)
		b.Board[len(b.Board)-1].Can_castle = true
	}
	for _, rank := range piecerows {
		var color int
		if rank == 1 {
			color = 1
		} else {
			color = -1
		}
		for _, file := range rookfiles {
			b.PlacePiece('r', color, file, rank)
			b.Board[len(b.Board)-1].Can_castle = true
		}
		for _, file := range knightfiles {
			b.PlacePiece('n', color, file, rank)
		}
		for _, file := range bishopfiles {
			b.PlacePiece('b', color, file, rank)
		}
		b.PlacePiece('q', color, queenfile, rank)
	}
	for _, rank := range pawnrows {
		var color int
		if rank == 2 {
			color = 1
		} else {
			color = -1
		}
		for file := 1; file <= 8; file++ {
			b.PlacePiece('p', color, file, rank)
			b.Board[len(b.Board)-1].Can_double_move = true
		}
	}
}
