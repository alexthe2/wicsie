package drawing

import (
	"github.com/fogleman/gg"
	"math"
	"wicsie/agents"
)

type Board struct {
	ctx *gg.Context
	w   int
	h   int
}

func CreateBoard(width, height int) *Board {
	return &Board{
		ctx: gg.NewContext(width, height),
		w:   width * 2,
		h:   height * 2,
	}
}

func (board *Board) DrawAgents(agents []*agents.Agent) {
	board.ctx.Clear()
	for _, agent := range agents {
		board.drawAgent(agent)
	}
	board.ctx.SetRGB(0, 0, 0)
	board.ctx.Fill()
}

func (board *Board) SaveBoard(filename string) {
	err := board.ctx.SavePNG(filename)
	if err != nil {
		return
	}
}

func (board *Board) drawAgent(agent *agents.Agent) {
	board.ctx.DrawCircle(math.Mod(agent.X, float64(board.w))*3, math.Mod(agent.Y, float64(board.h))*3, 2)
	board.ctx.SetRGB(agent.Health.GetColor())
	board.ctx.Fill()
}
