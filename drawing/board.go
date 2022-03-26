package drawing

import (
	"github.com/fogleman/gg"
	"image"
	"image/color"
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

func (board *Board) DrawGridMap(gridMap agents.GridMap) {
	for x := 0; x < gridMap.CountX; x++ {
		for y := 0; y < gridMap.CountY; y++ {
			r, g, b := colorForCell(gridMap.GetCell(x, y))
			for i := 0; i < gridMap.ChunkSize; i++ {
				for j := 0; j < gridMap.ChunkSize; j++ {
					board.ctx.SetPixel(x*gridMap.ChunkSize+i, y*gridMap.ChunkSize+j)
					board.ctx.SetRGBA(r, g, b, float64(board.mask.At(x*gridMap.ChunkSize+i, y*gridMap.ChunkSize+j).(color.NRGBA).A)/255.0)
					board.ctx.Fill()
				}
			}
		}
	}
}

func colorForCell(cell agents.Cell) (float64, float64, float64) {
	dominant := agents.Healthy
	dominantCount := cell.Healthy

	if cell.Infected > dominantCount {
		dominant = agents.Infected
		dominantCount = cell.Infected
	}

	if cell.Cured > dominantCount {
		dominant = agents.Cured
		dominantCount = cell.Cured
	}

	if dominantCount == 0 {
		return 0, 0, 0
	}

	return dominant.GetColor()
}
