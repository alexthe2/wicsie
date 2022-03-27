package simulation

import (
	"log"
	"math"
	"math/rand"
	"wicsie/agents"
	"wicsie/drawing"
	"wicsie/heatMapDecoder"
)

type Simulation struct {
	step   int
	agents []*agents.Agent

	spreading agents.Change
}

type Config struct {
	Weight        float64
	Width, Height float64
	Movement      func() agents.Movement
	Spreading     agents.Change

	HeatMap     [][]int
	LegendIndex []heatMapDecoder.LegendIndex
}

func CreateSimulation(config Config) *Simulation {
	sim := new(Simulation)
	sim.step = 0

	sim.spreading = config.Spreading
	sim.agents = make([]*agents.Agent, 0)

	for y := 0.0; y < config.Height; y += 1 {
		for x := 0.0; x < config.Width; x += 1 {
			n := int(config.Weight * float64(heatMapDecoder.GetMultiplier(config.LegendIndex, config.HeatMap[int(y)][int(x)])))
			for i := 0; i < n; i++ {
				if rand.Float64() < 0.5 {
					sim.agents = append(sim.agents, agents.CreateAgent(x+rV(), y+rV(), config.Width, config.Height, config.Movement()))
				}
			}
		}
	}

	return sim
}

func (sim *Simulation) InitInfect(probability float64) {
	for _, agent := range sim.agents {
		if rand.Float64() < probability {
			sim.spreading.Infect(agent)
		}
	}
}

func (sim *Simulation) Step() {
	agentCopy := make([]agents.Agent, len(sim.agents))
	for i := 0; i < len(sim.agents); i++ {
		agentCopy[i] = *sim.agents[i]
	}

	for i := 0; i < len(sim.agents); i++ {
		sim.agents[i].Move(agentCopy)
	}

	sim.spreading.ModifyHealth(sim.agents)

	log.Printf("Finished step %d", sim.step)
	sim.step++
}

func (sim *Simulation) DrawToBoard(board *drawing.Board) {
	board.DrawAgents(sim.agents)
}

func (sim *Simulation) GetAgents() []*agents.Agent {
	return sim.agents
}

func rV() float64 {
	return rand.Float64()*2 - 1
}

func (sim *Simulation) InfectAtPosition(x, y, prop float64) {
	for _, agent := range sim.agents {
		if math.Abs(agent.X-x) < 1 && math.Abs(agent.Y-y) < 1 {
			sim.spreading.Infect(agent)
		}
	}
}
