package agents

import (
	"math/rand"
	"sync"
	"wicsie/constants"
)

type nextStage struct {
	nextHealth         Status
	timeUntilNextState int
}

type GridSpread struct {
	tracked map[*Agent]nextStage

	gridMap *GridMap
}

func CreateGridSpread(gridMap *GridMap) *GridSpread {
	return &GridSpread{
		tracked: make(map[*Agent]nextStage),
		gridMap: gridMap,
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

	gs.tracked[agent] = nextStage{
		nextHealth:         Incubated,
		timeUntilNextState: calculateTime(constants.KBaseTimeUntilIncubation, constants.KVarianceInTimeUntilIncubation),
	}
}

func (gs *GridSpread) parallelInfect(agent *Agent, mutex *sync.Mutex) {
	mutex.Lock()
	if _, exists := gs.tracked[agent]; exists {
		mutex.Unlock()
		return
	}

	gs.tracked[agent] = nextStage{
		nextHealth:         Incubated,
		timeUntilNextState: calculateTime(constants.KBaseTimeUntilIncubation, constants.KVarianceInTimeUntilIncubation),
	}
	mutex.Unlock()
}

func (gs *GridSpread) handleTracked() {
	for k := range gs.tracked {
		gs.moveHealth(k)
	}
}

func (gs *GridSpread) moveHealth(agent *Agent) {
	track := gs.tracked[agent]

	track.timeUntilNextState--
	if track.timeUntilNextState <= 0 {
		agent.Health = track.nextHealth
		switch track.nextHealth {
		case Incubated:
			track.nextHealth = Infected
			track.timeUntilNextState = calculateTime(constants.KBaseTimeUntilInfection, constants.KVarianceInTimeUntilInfection)

		case Infected, UnknownInfected:
			track.nextHealth = Cured
			track.timeUntilNextState = calculateTime(constants.KBaseTimeUntilRecovery, constants.KVarianceInTimeUntilRecovery)

		case Cured:
			track.nextHealth = Healthy
			track.timeUntilNextState = calculateTime(constants.KBaseTimeUntilProtectionAfterCovidOver, constants.KVarianceInTimeUntilProtectionAfterCovidOver)

		case Healthy:
			delete(gs.tracked, agent)
			return

		}
	}

	gs.tracked[agent] = track
}

func (gs *GridSpread) handleNew(agents []*Agent) {
	var mutex sync.Mutex

	for _, agent := range agents {
		if agent.Health == Infected || agent.Health == UnknownInfected {
			var wg sync.WaitGroup
			neighbours := gs.gridMap.GetNeighbours(agent)
			wg.Add(len(neighbours))
			for _, partner := range neighbours {
				go func(partner *Agent) {
					if partner.Health == Healthy && rand.Float64() < constants.KProbabilityOfInfection {
						gs.parallelInfect(partner, &mutex)
					}
					wg.Done()
				}(partner)
			}
			wg.Wait()
		}
	}
}

func calculateTime(base int, variance int) int {
	if variance == 0 {
		return base
	}

	return base + rand.Intn(variance*2) - variance
}
