package drawing

import (
	"github.com/fogleman/gg"
	"image"
	"math"
	"wicsie/agents"
)

type Board struct {
	ctx     *gg.Context
	mask    image.Image
	w       int
	h       int
	scaling float64
}

func CreateBoard(width, height int, mask image.Image, scaling float64) *Board {
	w := int(float64(width) * scaling)
	h := int(float64(height) * scaling)

	return &Board{
		ctx:     gg.NewContext(w, h),
		mask:    mask,
		w:       w,
		h:       h,
		scaling: scaling,
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
	board.ctx.DrawCircle(math.Mod(agent.X, float64(board.w))*board.scaling, math.Mod(agent.Y, float64(board.h))*board.scaling, 2)
	board.ctx.SetRGB(agent.Health.GetColor())
	board.ctx.Fill()
}

func (board *Board) drawGridMap(gridMap agents.GridMap) {
	for x := 0; x < board.w; x++ {
		for y := 0; y < board.h; y++ {
			pixel := gridMap.
		}
	}
}

func chunkFor(x, y, cw, ch int) (int, int) {
	return x / cw, y / ch
}
