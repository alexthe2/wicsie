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
	gridMap.clearGridMap()
}

func (gridMap *GridMap) clearGridMap() {
	for i := 0; i < gridMap.CountX; i++ {
		for j := 0; j < gridMap.CountY; j++ {
			gridMap.Cells[i][j] = Cell{
				Agents: make([]*Agent, 0),
			}
		}
	}
}

func (gridMap *GridMap) recalculateGridMap(agents []*Agent) {
	for _, agent := range agents {
		corX := int(agent.X / float64(gridMap.chunkSize))
		corY := int(agent.Y / float64(gridMap.chunkSize))
		gridMap.Cells[corX][corY].Agents = append(gridMap.Cells[corX][corY].Agents, agent)
		switch agent.Health {
		case Incubated:
		case Healthy:
			gridMap.Cells[corX][corY].Healthy++
		case Infected:
		case UnknownInfected:
			gridMap.Cells[corX][corY].Infected++
		case Cured:
			gridMap.Cells[corX][corY].Cured++
		}
	}
}
