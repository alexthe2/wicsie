package agents

import (
	"github.com/enriquebris/goconcurrentqueue"
	"log"
	"math/rand"
	"os"
	"sync"
	"wicsie/constants"
)

type ParallelGridSpread struct {
	trackedMap  sync.Map
	trackedList []*Agent

	gridMap *GridMap
	LOGGER  *log.Logger
}

func CreateParallelGridSpread(gridMap *GridMap) *ParallelGridSpread {
	return &ParallelGridSpread{
		trackedMap:  sync.Map{},
		trackedList: []*Agent{},
		gridMap:     gridMap,
		LOGGER:      log.New(os.Stdout, "[PGSPR] ", log.Ltime),
	}
}

func (gs *ParallelGridSpread) ModifyHealth(allAgents []*Agent) {
	gs.LOGGER.Println("Modifying health for tracked agents")
	gs.parallelHandleTracked()
	gs.parallelHandleUntracked()
	gs.LOGGER.Println("Finished modifying health for tracked agents")
	gs.LOGGER.Println("Finding new infections")
	gs.parallelHandleNew(allAgents)
	gs.LOGGER.Println("Finished finding new infections")
}

func (gs *ParallelGridSpread) parallelHandleUntracked() {
	tracked := make([]*Agent, 0)
	for _, agent := range gs.trackedList {
		if agent != nil {
			tracked = append(tracked, agent)
		}
	}

	gs.trackedList = tracked
}

func (gs *ParallelGridSpread) Infect(agent *Agent) {
	if _, exists := gs.trackedMap.Load(agent); exists {
		return
	}

	gs.trackedMap.Store(agent, NextStage{
		NextHealth:         Incubated,
		TimeUntilNextState: calculateTime(constants.KBaseTimeUntilIncubation, constants.KVarianceInTimeUntilIncubation),
	})

	gs.trackedList = append(gs.trackedList, agent)
}

func (gs *ParallelGridSpread) parallelHandleTracked() {
	for k := range gs.trackedList {
		if !gs.parallelMoveHealth(gs.trackedList[k]) {
			gs.trackedList[k] = nil
		}
	}
}

func (gs *ParallelGridSpread) parallelMoveHealth(agent *Agent) bool {
	trackA, ok := gs.trackedMap.Load(agent)
	if !ok {
		gs.LOGGER.Println("Agent not found in tracked map")
		return false
	}
	track := trackA.(NextStage)

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
			gs.trackedMap.Delete(agent)
			return false
		}
	}

	gs.trackedMap.Store(agent, track)
	return true
}

func (gs *ParallelGridSpread) parallelHandleNew(agents []*Agent) {
	queue := goconcurrentqueue.NewFIFO()

	for _, agent := range agents {
		if agent.Health == Infected || agent.Health == UnknownInfected {
			var wg sync.WaitGroup
			neighbours := gs.gridMap.GetNeighbours(agent)
			wg.Add(len(neighbours))
			for _, partner := range neighbours {
				go func(partner *Agent) {
					if partner.Health == Healthy && rand.Float64() < constants.KProbabilityOfInfection {
						err := queue.Enqueue(partner)
						if err != nil {
							gs.LOGGER.Println("Error enqueuing partner: ", err)
						}
					}
					wg.Done()
				}(partner)
			}
			wg.Wait()
		}
	}

	length := queue.GetLen()
	for i := 0; i < length; i++ {
		agent, err := queue.Dequeue()
		if err != nil {
			gs.LOGGER.Println("Error dequeuing agent: ", err)
		}
		if agent != nil {
			gs.Infect(agent.(*Agent))
		}
	}
}
