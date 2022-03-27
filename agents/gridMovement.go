package agents

import (
	"math"
	"math/rand"
	"wicsie/constants"
)

type GridMovement struct {
	moveX, moveY float64

	timeToChange int
	ttcMax       int

	grid         *GridMap
	heatChunkMap [][]int
}

func CreateGridMovement(ttcMax int, grid *GridMap, heatChunkMap [][]int) *GridMovement {
	movement := GridMovement{
		ttcMax:       ttcMax,
		timeToChange: 0,
		grid:         grid,
		heatChunkMap: heatChunkMap,
	}

	return &movement
}

func (movement *GridMovement) Move(_ []Agent, agent Agent) (float64, float64) {
	movement.timeToChange--
	if movement.timeToChange <= 0 {
		movement.generateMovementBehaviour(agent)
	}

	return movement.moveX, movement.moveY
}

func (movement *GridMovement) generateMovementBehaviour(agent Agent) {
	xChunkAgent := int(math.Ceil(float64(agent.X) / float64(constants.KChunkSize)))
	yChunkAgent := int(math.Ceil(float64(agent.Y) / float64(constants.KChunkSize)))
	decisionProbability := rand.Float64()
	xNeighbours := [8]int{0, 1, 1, 1, 0, -1, -1, -1}
	yNeighbours := [8]int{1, 1, 0, -1, -1, -1, 0, 1}

	//We are staying in the same cell (probability 70%)
	if decisionProbability < 0.70 {
		movement.moveX = rand.Float64() - 0.5
		movement.moveY = rand.Float64() - 0.5
	} else {
		//We are going to other cell
		maxScore := 0
		iBest := 0
		jBest := 0
		for _, i := range xNeighbours {
			for _, j := range yNeighbours {
				if xChunkAgent+i > 0 && yChunkAgent+j > 0 { //} && xChunkAgent+i < width && yChunkAgent+j < height {
					chunk := movement.grid.Cells[xChunkAgent+i][yChunkAgent+j]
					score := 0
					if movement.heatChunkMap[xChunkAgent+i][yChunkAgent+j] == 1 {
						score = -100
						//fmt.Println("score: ", score)
						//fmt.Println("maxScore: ", maxScore)
					} else {
						score = chunk.Healthy + chunk.Cured - chunk.Infected
					}
					if score > maxScore {
						maxScore = score
						//fmt.Println("newScore: ", score)
						iBest = i
						jBest = j
					}
				}

			}
		}
		movement.moveX += constants.KChunkSize * float64(iBest)
		movement.moveY += constants.KChunkSize * float64(jBest)
	}

	movement.timeToChange = rand.Intn(movement.ttcMax) + 10

}
