package agents

import (
	"log"
	"math/rand"
	"os"
	"wicsie/constants"
)

type NextStage struct {
	NextHealth         Status
	TimeUntilNextState int
}

type GridSpread struct {
	tracked map[*Agent]NextStage

	gridMap *GridMap
	LOGGER  *log.Logger
}

func CreateGridSpread(gridMap *GridMap) *GridSpread {
	return &GridSpread{
		tracked: make(map[*Agent]NextStage),
		gridMap: gridMap,
		LOGGER:  log.New(os.Stdout, "[GRSPR] ", log.Ltime),
	}
}

func (gs *GridSpread) ModifyHealth(allAgents []*Agent) {
	gs.handleTracked()
	gs.handleNew(allAgents)
}

func (gs *GridSpread) Infect(agent *Agent) {
	if _, exists := gs.tracked[agent]; exists {
		return
	}

	gs.tracked[agent] = NextStage{
		NextHealth:         Incubated,
		TimeUntilNextState: calculateTime(constants.KBaseTimeUntilIncubation, constants.KVarianceInTimeUntilIncubation),
	}
}

func (gs *GridSpread) handleTracked() {
	for k := range gs.tracked {
		gs.moveHealth(k)
	}
}

func (gs *GridSpread) moveHealth(agent *Agent) {
	track := gs.tracked[agent]

	track.TimeUntilNextState--
	if track.TimeUntilNextState <= 0 {
		agent.Health = track.NextHealth
		switch track.NextHealth {
		case Incubated:
			track.NextHealth = Infected
			track.TimeUntilNextState = calculateTime(constants.KBaseTimeUntilInfection, constants.KVarianceInTimeUntilInfection)

		case Infected, UnknownInfected:
			track.NextHealth = Cured
			track.TimeUntilNextState = calculateTime(constants.KBaseTimeUntilRecovery, constants.KVarianceInTimeUntilRecovery)

		case Cured:
			track.NextHealth = Healthy
			track.TimeUntilNextState = calculateTime(constants.KBaseTimeUntilProtectionAfterCovidOver, constants.KVarianceInTimeUntilProtectionAfterCovidOver)

		case Healthy:
			delete(gs.tracked, agent)
			return

		}
	}

	gs.tracked[agent] = track
}

func (gs *GridSpread) handleNew(agents []*Agent) {
	for _, agent := range agents {
		if agent.Health == Infected || agent.Health == UnknownInfected {
			neighbours := gs.gridMap.GetNeighbours(agent)
			for _, partner := range neighbours {
				if partner.Health == Healthy && rand.Float64() < constants.KProbabilityOfInfection {
					gs.Infect(partner)
				}
			}
		}
	}
}

func calculateTime(base int, variance int) int {
	if variance == 0 {
		return base
	}

	return base + rand.Intn(variance*2) - variance
}
