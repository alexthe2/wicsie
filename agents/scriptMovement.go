package agents

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
)

type Behaviour struct {
	AreaXStart, AreaXEnd, AreaYStart, AreaYEnd int

	Min, Max             float64
	ChanceForSmallMove   float64
	MaximalSmallMove     float64
	MaximalSmallMoveTime int

	ChanceForBigMove   float64
	MaximalBigMove     float64
	MaximalBigMoveTime int

	ChanceForReturningHome float64
	ChanceForExtremeMove   float64
}

type CurrentMove struct {
	x, y     float64
	timeLeft int
}

type ScriptMovement struct {
	gridMap  *GridMap
	colorMap [][]int

	behaviours []Behaviour

	currentMove CurrentMove
}

func DecodeFile(fileName string) []Behaviour {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("[SCRIPT] Could not open file %s: %s", fileName, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("[SCRIPT] Could not close file %s: %s", fileName, err)
		}
	}(file)

	behaviours := make([]Behaviour, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if line[0] == '#' {
			continue
		}
		behaviour := Behaviour{}
		_, err := fmt.Sscanf(line, "WITHIN AREA ((%d,%d),(%d,%d)) FOR CASES BETWEEN %f AND %f CHANCE FOR SMALL MOVE IS %f WHERE SMALL MOVE IS %f OVER %d DAYS AND CHANCE FOR BIG MOVE IS %f WHERE BIG MOVE IS %f OVER %d DAYS WHILE GOING HOME IS %f AND EXTREME MOVE CHANCE IS %f", &behaviour.AreaXStart, &behaviour.AreaYStart, &behaviour.AreaXEnd, &behaviour.AreaYEnd, &behaviour.Min, &behaviour.Max, &behaviour.ChanceForSmallMove, &behaviour.MaximalSmallMove, &behaviour.MaximalSmallMoveTime, &behaviour.ChanceForBigMove, &behaviour.MaximalBigMove, &behaviour.MaximalBigMoveTime, &behaviour.ChanceForReturningHome, &behaviour.ChanceForExtremeMove)
		if err != nil {
			log.Fatalf("[SCRIPT] Could not parse line %s: %s", line, err)
		}
		behaviours = append(behaviours, behaviour)
	}

	return behaviours
}

func CreateScriptMovement(gridMap *GridMap, behaviours []Behaviour, colorMap [][]int) *ScriptMovement {
	return &ScriptMovement{
		gridMap:     gridMap,
		behaviours:  behaviours,
		colorMap:    colorMap,
		currentMove: CurrentMove{},
	}
}

func (movement *ScriptMovement) Move(agentsAround []Agent, me Agent) (float64, float64) {
	if movement.colorMap[int(me.X)%504][int(me.Y)%599] == 0 {
		movement.currentMove.timeLeft = 0
		return 0, 0
	}

	if movement.currentMove.timeLeft > 0 {
		movement.currentMove.timeLeft--
		return movement.currentMove.x, movement.currentMove.y
	}

	cell := movement.gridMap.GetCellForAgent(me)
	covidPercentage := float64(cell.Cured+cell.Infected+cell.Healthy) / math.Max(float64(cell.Infected), 1)
	for _, behaviour := range movement.behaviours {
		if behaviour.Min <= covidPercentage && covidPercentage <= behaviour.Max {
			if randFloat() < behaviour.ChanceForSmallMove {
				movement.currentMove.x = randFloat()*behaviour.MaximalSmallMove*2 - behaviour.MaximalSmallMove
				movement.currentMove.y = randFloat()*behaviour.MaximalSmallMove*2 - behaviour.MaximalSmallMove
				movement.currentMove.timeLeft = behaviour.MaximalSmallMoveTime/2 + rand.Intn(behaviour.MaximalSmallMoveTime)/2
				return movement.currentMove.x, movement.currentMove.y
			}
			if randFloat() < behaviour.ChanceForBigMove {
				movement.currentMove.x = randFloat()*behaviour.MaximalBigMove*2 - behaviour.MaximalBigMove
				movement.currentMove.y = randFloat()*behaviour.MaximalBigMove*2 - behaviour.MaximalBigMove
				movement.currentMove.timeLeft = behaviour.MaximalBigMoveTime/2 + rand.Intn(behaviour.MaximalBigMoveTime)/2
				return movement.currentMove.x, movement.currentMove.y
			}
			if randFloat() < behaviour.ChanceForReturningHome {
				movement.currentMove.x = 0
				movement.currentMove.y = 0
				movement.currentMove.timeLeft = 0
				return movement.currentMove.x, movement.currentMove.y
			}
		}
	}

	return randFloat()*.2 - .1, randFloat()*.2 - .1
}

func randFloat() float64 {
	return rand.Float64()
}
