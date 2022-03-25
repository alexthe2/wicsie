package agents

import (
	"math/rand"
)

type RandomMovement struct {
	moveX, moveY float64

	timeToChange int
}

func CreateRandomMovement() *RandomMovement {
	movement := RandomMovement{}
	movement.generateRandomMovement()

	return &movement
}

func (movement *RandomMovement) Move(_ []Agent, _ Status) (float64, float64) {
	movement.timeToChange--
	if movement.timeToChange == 0 {
		movement.generateRandomMovement()
	}

	return movement.moveX, movement.moveY
}

func (movement *RandomMovement) generateRandomMovement() {
	movement.moveX = rand.Float64() - 0.5
	movement.moveY = rand.Float64() - 0.5

	movement.timeToChange = rand.Intn(1000) + 1000
}
