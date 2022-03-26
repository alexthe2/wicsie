package agents

import "math"

type Cell struct {
	Agents   []*Agent
	Healthy  int
	Infected int
	Cured    int
}

type GridMap struct {
	CountX, CountY int
	width, height  int
	chunkSize      int
	Cells          [][]Cell
}

func CreateGridMap(width, height, chunkSize int) *GridMap {
	countX := int(math.Ceil(float64(width) / float64(chunkSize)))
	countY := int(math.Ceil(float64(height) / float64(chunkSize)))

	arr := make([][]Cell, countY)
	for i := 0; i < countY; i++ {
		arr[i] = make([]Cell, countX)
	}

	return &GridMap{
		width:     width,
		height:    height,
		chunkSize: chunkSize,
		CountX:    countX,
		CountY:    countY,
		Cells:     arr,
	}
}

func (gridMap *GridMap) UpdateGridMap(agents []*Agent) {

}
