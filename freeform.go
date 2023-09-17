package hexz

import "fmt"

//
// The freeform single-player hexz game.
//

type GameEngineFreeform struct {
	board *Board
}

func (g *GameEngineFreeform) GameType() GameType { return gameTypeFreeform }
func (g *GameEngineFreeform) Board() *Board      { return g.board }

func (g *GameEngineFreeform) Init() {
	flatFields, fields := makeFields()
	b := &Board{
		Turn:       1, // Player 1 begins
		FlatFields: flatFields,
		Fields:     fields,
		State:      Initial,
	}
	numPlayers := g.NumPlayers()
	b.Score = make([]int, numPlayers)
	b.Resources = make([]ResourceInfo, numPlayers)
	for i := 0; i < numPlayers; i++ {
		b.Resources[i] = g.InitialResources()
	}
	g.board = b
	g.board.State = Running
}

func (g *GameEngineFreeform) NumPlayers() int {
	return 1
}
func (g *GameEngineFreeform) MoveHistory() []GameEngineMove {
	panic("Not implemented")
}

func (g *GameEngineFreeform) ValidCellTypes() []CellType {
	r := make([]CellType, 0, cellTypeLen)
	for i, v := range g.InitialResources().NumPieces {
		if v != 0 {
			r = append(r, CellType(i))
		}
	}
	return r
}

func (g *GameEngineFreeform) InitialResources() ResourceInfo {
	var ps [cellTypeLen]int
	ps[cellNormal] = -1 // unlimited
	ps[cellFire] = -1
	ps[cellFlag] = -1
	ps[cellPest] = -1
	ps[cellDeath] = -1
	return ResourceInfo{
		NumPieces: ps}
}

func (g *GameEngineFreeform) Reset()       { g.Init() }
func (g *GameEngineFreeform) IsDone() bool { return false }
func (g *GameEngineFreeform) Winner() (playerNum int) {
	return 0 // No one ever wins here.
}

func (g *GameEngineFreeform) MakeMove(m GameEngineMove) bool {
	board := g.board
	if !board.valid(idx{m.Row, m.Col}) {
		// Invalid move request.
		return false
	}
	board.Move++
	f := &board.Fields[m.Row][m.Col]
	f.Owner = board.Turn
	f.Type = m.CellType
	board.Turn++
	if board.Turn > 2 {
		board.Turn = 1
	}
	f.Value = 1
	return true
}

func (g *GameEngineFreeform) Encode() ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (g *GameEngineFreeform) Decode([]byte) error {
	return fmt.Errorf("not implemented")
}
