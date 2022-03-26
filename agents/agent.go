package agents

import (
	"math"
	"math/rand"
)

type Agent struct {
	X, Y          float64
	width, height float64

	Health   Status
	Movement Movement
}

func CreateAgent(x, y, width, height float64, movement Movement) *Agent {
	return &Agent{
		X:        x,
		Y:        y,
		width:    width,
		height:   height,
		Health:   Healthy,
		Movement: movement,
	}
}

func CreateAgentAtRandomPosition(width, height float64, movement Movement) *Agent {
	return &Agent{
		X:        float64(rand.Intn(int(width))),
		Y:        float64(rand.Intn(int(height))),
		width:    width,
		height:   height,
		Health:   Healthy,
		Movement: movement,
	}
}

func (agent *Agent) Move(agents []Agent) {
	dx, dy := agent.Movement.Move(agents, agent.Health)
	agent.X = math.Mod(agent.X+dx, agent.width)
	agent.Y = math.Mod(agent.Y+dy, agent.height)
}
