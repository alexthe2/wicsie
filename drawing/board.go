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
			r, g, b, a := colorForCell(gridMap.GetCell(x, y))

			for i := 0; i < gridMap.ChunkSize; i++ {
				for j := 0; j < gridMap.ChunkSize; j++ {
					board.ctx.SetPixel(x*gridMap.ChunkSize+i, y*gridMap.ChunkSize+j)

					if board.mask.At(x*gridMap.ChunkSize+i, y*gridMap.ChunkSize+j).(color.NRGBA).A == 0 {
						board.ctx.SetRGBA(.6784, .847, .901, .8)
					} else {
						board.ctx.SetRGBA(r, g, b, a)
					}

					board.ctx.Fill()
				}
			}
		}
	}
}

func colorForCell(cell agents.Cell) (float64, float64, float64, float64) {
	/*all := cell.Healthy + cell.Infected + cell.Cured
	infectedPart := float64(cell.Infected) / float64(all)
	healthyPart := float64(cell.Healthy) / float64(all)
	curedPart := float64(cell.Cured) / float64(all)

	r := math.Max(infectedPart*255-(curedPart*0.8*255), 20)
	g := math.Max(healthyPart*255-(curedPart*0.8*255), 20)
	b := curedPart * 0.2 * 255

	return r / 255.0, g / 255.0, b / 255.0, dim(all)*/

	dominant := agents.Healthy
	dominantCount := cell.Healthy

	if cell.Cured > dominantCount {
		dominant = agents.Cured
		dominantCount = cell.Cured
	}

	if cell.Infected > dominantCount {
		dominant = agents.Infected
		dominantCount = cell.Infected
	}

	if dominantCount == 0 {
		return 0, 0, 0, 0
	}

	r, g, b := dominant.GetColor()
	return r, g, b, dim(cell.Healthy + cell.Infected + cell.Cured)

}

func dim(count int) float64 { return math.Min(255, float64(count)/800) }
