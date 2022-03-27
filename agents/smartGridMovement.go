package agents

import (
	"math"
	"math/rand"
)

const chanceForSmallChange = 0.2
const chanceForBigChange = 0.65
const chanceForExtremeChange = 0.87
const chanceForSuperFly = 0.97

type SmartGridMovement struct {
	homeX, homeY float64

	destinationX, destinationY float64
	speed                      float64

	heatMap [][]int
}

func CreateSmartGridMovement(heatMap [][]int) *SmartGridMovement {
	movement := SmartGridMovement{
		heatMap: heatMap,
		homeX:   -1,
		homeY:   -1,
	}

	return &movement
}

func (movement *SmartGridMovement) Move(agents []Agent, agent Agent) (float64, float64) {
	if (movement.homeX == -1) || (movement.homeY == -1) {
		movement.homeX = agent.X
		movement.homeY = agent.Y

		movement.findDestination(agent)
	}

	if (int(movement.homeX-agent.X) <= 1) && (int(movement.homeY-agent.Y) <= 1) {
		movement.findDestination(agent)
	}

	dx := movement.destinationX - agent.X
	dy := movement.destinationY - agent.Y

	if math.Abs(dx) < 0.1 {
		dx = 0
	}
	if math.Abs(dy) < 0.1 {
		dy = 0
	}

	if (dx == 0 && dy == 0) || movement.heatMap[int(agent.X)][int(agent.Y)] == 0 {
		movement.goToHome()
		return 0, 0
	}

	return applySpeedViaNormalization(dx, dy, movement.speed)
}

func (movement *SmartGridMovement) findDestination(agent Agent) {
	distanceChange := rand.Float64()

	if distanceChange > chanceForSuperFly {
		movement.destinationX = agent.X + rand.Float64()*3000 - 1500
		movement.destinationY = agent.Y + rand.Float64()*3000 - 1500
		movement.speed = rand.Float64()*10 + 2
		return
	} else if distanceChange > chanceForExtremeChange {
		movement.destinationX = agent.X + rand.Float64()*1500 - 750
		movement.destinationY = agent.Y + rand.Float64()*1500 - 750
	} else if distanceChange > chanceForBigChange {
		movement.destinationX = agent.X + rand.Float64()*500 - 250
		movement.destinationY = agent.Y + rand.Float64()*500 - 250
	} else if distanceChange > chanceForSmallChange {
		movement.destinationX = agent.X + rand.Float64()*100 - 50
		movement.destinationY = agent.Y + rand.Float64()*100 - 50
	} else {
		movement.destinationX = agent.X
		movement.destinationY = agent.Y
	}

	movement.speed = rand.Float64() * 4
}

func (movement *SmartGridMovement) goToHome() {
	movement.destinationX = movement.homeX
	movement.destinationY = movement.homeY

	movement.speed = rand.Float64()*0.5 + 0.5
}

func applySpeedViaNormalization(dx, dy, speed float64) (float64, float64) {
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	if distance > 0 {
		dx = dx / distance * speed
		dy = dy / distance * speed
	}

	return dx, dy
}
