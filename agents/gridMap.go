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
	ChunkSize      int
	Cells          [][]Cell
}

func CreateGridMap(width, height, chunkSize int) *GridMap {
	countX := int(math.Ceil(float64(width) / float64(chunkSize)))
	countY := int(math.Ceil(float64(height) / float64(chunkSize)))

	arr := make([][]Cell, countX)
	for i := 0; i < countX; i++ {
		arr[i] = make([]Cell, countY)
	}

	return &GridMap{
		width:     width,
		height:    height,
		ChunkSize: chunkSize,
		CountX:    countX,
		CountY:    countY,
		Cells:     arr,
	}
}

func (gridMap *GridMap) UpdateGridMap(agents []*Agent) {
	gridMap.clearGridMap()
	gridMap.recalculateGridMap(agents)
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
		corX := int(agent.X / float64(gridMap.ChunkSize))
		corY := int(agent.Y / float64(gridMap.ChunkSize))
		gridMap.Cells[corX][corY].Agents = append(gridMap.Cells[corX][corY].Agents, agent)
		switch agent.Health {
		case Incubated, Healthy:
			gridMap.Cells[corX][corY].Healthy++
		case Infected, UnknownInfected:
			gridMap.Cells[corX][corY].Infected++
		case Cured:
			gridMap.Cells[corX][corY].Cured++
		}
	}
}

func (gridMap *GridMap) GetCell(x, y int) Cell {
	return gridMap.Cells[x][y]
}

func (gridMap *GridMap) GetNeighbours(agent *Agent) []*Agent {
	return gridMap.Cells[int(agent.X/float64(gridMap.ChunkSize))][int(agent.Y/float64(gridMap.ChunkSize))].Agents
}
