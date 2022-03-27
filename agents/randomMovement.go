package agents

import (
	"math/rand"
)

const chanceForCrazy = 0.001
const chanceForJustGoingHomeImmediately = 0.001

type RandomMovement struct {
	moveX, moveY  float64
	width, height float64
	homeX, homeY  float64

	timeToChange int
	ttcMax       int

	heatMap [][]int
}

func CreateRandomMovement(ttcMax int, heatMap [][]int, width, height float64) *RandomMovement {
	movement := RandomMovement{ttcMax: ttcMax, heatMap: heatMap, homeX: -1, homeY: -1}
	movement.generateRandomMovement()

	return &movement
}

func (movement *RandomMovement) Move(_ []Agent, agent Agent) (float64, float64) {
	//if movement.heatMap[int(math.Mod(agent.X+movement.moveX, movement.width))][int(math.Mod(agent.Y+movement.height, movement.height))] == 0 {
	//	movement.generateRandomMovement()
	//	return 0, 0
	//}

	if movement.homeX == -1 {
		movement.homeX = agent.X
		movement.homeY = agent.Y
	}

	if rand.Float64() < chanceForJustGoingHomeImmediately {
		movement.generateRandomMovement()
		return movement.homeX - agent.X, movement.homeY - agent.Y
	}

	movement.timeToChange--
	if movement.timeToChange == 0 {
		movement.generateRandomMovement()
	}

	return movement.moveX, movement.moveY
}

func (movement *RandomMovement) generateRandomMovement() {
	movement.moveX = rand.Float64() - 0.5
	movement.moveY = rand.Float64() - 0.5

	if rand.Float64() < chanceForCrazy {
		movement.moveX *= 50
		movement.moveY *= 50
	}

	movement.timeToChange = rand.Intn(movement.ttcMax) + 10
}
