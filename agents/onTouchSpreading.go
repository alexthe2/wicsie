package agents

import (
	"math"
	"math/rand"
)

const consideredTouchBase = 20
const consideredTouchVariance = 10

const timeUntilIncubation = 6
const timeUntilSpreadBase = 12
const timeUntilSpreadVariance = 6
const timeUntilCureBase = 48
const timeUntilCureVariance = 24
const timeUntilCureOverBase = 300
const timeUntilCureOverVariance = 150

type nextBehaviour struct {
	agent          *Agent
	nextHealth     Status
	iterationsLeft int
}

type OnTouchSpreading struct {
	tracked []nextBehaviour
}

func CreateOnTouchSpreading() *OnTouchSpreading {
	return &OnTouchSpreading{
		tracked: []nextBehaviour{},
	}
}

func (spreading *OnTouchSpreading) ModifyHealth(agents []*Agent) {
	spreading.newInfections(agents)
	spreading.handleTracked()
}

func (spreading *OnTouchSpreading) Infect(agent *Agent) {
	infection := nextBehaviour{
		agent:          agent,
		nextHealth:     Incubated,
		iterationsLeft: timeUntilIncubation,
	}

	if !contains(spreading.tracked, infection) {
		spreading.tracked = append(spreading.tracked, infection)
	}
}

func (spreading *OnTouchSpreading) handleTracked() {
	for i, _ := range spreading.tracked {
		spreading.tracked[i].iterationsLeft--
		if spreading.tracked[i].iterationsLeft == 0 {
			spreading.tracked[i].agent.Health = spreading.tracked[i].nextHealth
			switch spreading.tracked[i].nextHealth {
			case Incubated:
				spreading.tracked[i].nextHealth = Infected
				spreading.tracked[i].iterationsLeft = timeUntilSpreadBase + varianceBoth(timeUntilSpreadVariance)

			case Infected:
				spreading.tracked[i].nextHealth = Cured
				spreading.tracked[i].iterationsLeft = timeUntilCureBase + varianceBoth(timeUntilCureVariance)

			case Cured:
				spreading.tracked[i].nextHealth = Healthy
				spreading.tracked[i].iterationsLeft = timeUntilCureOverBase + varianceBoth(timeUntilCureOverVariance)

			case Healthy:
				spreading.tracked[i].iterationsLeft = -1
			}
		}
	}

	tmp := spreading.tracked[:0]
	for _, tracked := range spreading.tracked {
		if tracked.iterationsLeft >= 0 {
			tmp = append(tmp, tracked)
		}
	}
	spreading.tracked = tmp
}

func (spreading *OnTouchSpreading) newInfections(agents []*Agent) {
	for _, agent := range agents {
		if agent.Health != UnknownInfected && agent.Health != Infected {
			continue
		}

		for _, partner := range agents {
			if partner.Health != Healthy || distance(*agent, *partner) > float64(consideredTouchBase+varianceBoth(consideredTouchVariance)) {
				continue
			}

			spreading.Infect(partner)
		}
	}
}

func distance(agent1 Agent, agent2 Agent) float64 {
	return math.Sqrt(math.Pow(agent1.X-agent2.X, 2) + math.Pow(agent1.Y-agent2.Y, 2))
}

func contains(s []nextBehaviour, e nextBehaviour) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func varianceBoth(num int) int {
	return rand.Intn(num*2+1) - num
}
