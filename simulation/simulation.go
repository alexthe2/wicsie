package simulation

import (
	"log"
	"wicsie/agents"
	"wicsie/drawing"
)

type Simulation struct {
	step   int
	agents []*agents.Agent

	spreading agents.Change
}

func CreateSimulation(amount int, width int, height int, movement func() agents.Movement, spreading agents.Change) *Simulation {
	sim := new(Simulation)
	sim.step = 0
	sim.spreading = spreading

	sim.agents = make([]*agents.Agent, amount)
	for i := 0; i < amount; i++ {
		sim.agents[i] = agents.CreateAgentAtRandomPosition(float64(width), float64(height), movement())
	}

	return sim
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
