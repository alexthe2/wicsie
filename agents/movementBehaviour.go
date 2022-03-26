package agents

import (
	"math"
	"math/rand"
	"wicsie/constants"
)

type MovementBehaviour struct {
	moveX, moveY float64

	timeToChange int
	ttcMax       int

	grid *GridMap
}

func CreateMovementBehaviour(ttcMax int, grid *GridMap, agent Agent) *MovementBehaviour {
	movement := MovementBehaviour{
		ttcMax: ttcMax,
		grid:   grid,
	}
	movement.generateMovementBehaviour(agent)

	return &movement
}

func (movement *MovementBehaviour) Move(_ []Agent, agent Agent) (float64, float64) {
	movement.timeToChange--
	if movement.timeToChange == 0 {
		movement.generateMovementBehaviour(agent)
	}

	return movement.moveX, movement.moveY
}

func (movement *MovementBehaviour) generateMovementBehaviour(agent Agent) {
	xChunkAgent := int(math.Ceil(float64(agent.X) / float64(constants.KChunkSize)))
	yChunkAgent := int(math.Ceil(float64(agent.Y) / float64(constants.KChunkSize)))
	decisionProbability := rand.Float64()
	xNeighbours := [8]int{0, 1, 1, 1, 0, -1, -1, -1}
	yNeighbours := [8]int{1, 1, 0, -1, -1, -1, 0, 1}

	//We are staying in the same cell
	if decisionProbability < 0.55 {
		movement.moveX = rand.Float64() - 0.5
		movement.moveY = rand.Float64() - 0.5
	} else {
		maxScore := 0
		iBest := 0
		jBest := 0
		for _, i := range xNeighbours {
			for _, j := range yNeighbours {
				chunk := movement.grid.Cells[xChunkAgent+i][yChunkAgent+j]
				score := chunk.Healthy + chunk.Cured - chunk.Infected
				if score > maxScore {
					maxScore = score
					iBest = i
					jBest = j
				}
			}
		}
		movement.moveX += constants.KChunkSize * float64(iBest)
		movement.moveY += constants.KChunkSize * float64(jBest)
	}

	movement.timeToChange = rand.Intn(movement.ttcMax) + 10

}
