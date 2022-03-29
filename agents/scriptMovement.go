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
	Alpha								 float64
	DistFact   					 float64
	Beta                 float64
	TimeFact             float64

	ChanceForReturningHome float64
	ChanceForExtremeMove   float64
	Scaling								 int
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
		_, err := fmt.Sscanf(line, "WITHIN AREA ((%d,%d),(%d,%d)) FOR CASES BETWEEN %f AND %f ALPHA %f DISTANCE FACTOR %f BETA %f TIME FACTOR %f STAY HOME CHANCE %f EXTREME MOVE CHANCE %f SCALING %d", &behaviour.AreaXStart, &behaviour.AreaYStart, &behaviour.AreaXEnd, &behaviour.AreaYEnd, &behaviour.Min, &behaviour.Max, &behaviour.Alpha, &behaviour.DistFact, &behaviour.Beta, &behaviour.TimeFact, &behaviour.ChanceForReturningHome, &behaviour.ChanceForExtremeMove, &behaviour.Scaling)
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
			if randFloat() < behaviour.ChanceForExtremeMove {
				movement.currentMove.x = randFloat()*504 - me.X
				movement.currentMove.y = randFloat()*599 - me.Y
				movement.currentMove.timeLeft = 0
			} else if randFloat() < behaviour.ChanceForReturningHome {
				movement.currentMove.x = 0
				movement.currentMove.y = 0
				movement.currentMove.timeLeft = 0
				return movement.currentMove.x, movement.currentMove.y
			} else {
				distanceInKm := getRandomDistanceTravelled(behaviour.Alpha, behaviour.DistFact) * float64(behaviour.Scaling)
				timeSpent := getRandomTimeSpent(behaviour.Beta, behaviour.TimeFact) * float64(behaviour.Scaling)
				direction := randFloat()* math.Pi * 2

				// Magic number
				pixelSizeInKm := 4.3
				distance := distanceInKm / pixelSizeInKm

				movement.currentMove.x = distance * math.Cos(direction)
				movement.currentMove.y = distance * math.Sin(direction)
				movement.currentMove.timeLeft = int(timeSpent)
				return movement.currentMove.x, movement.currentMove.y
			}
		}
	}

	return randFloat()*.2 - .1, randFloat()*.2 - .1
}

func getRandomDistanceTravelled(alpha float64, distFact float64) float64 {
	return math.Pow(distFact * randFloat(), 1 / (-1-alpha))
}

func getRandomTimeSpent(beta float64, timeFact float64) float64 {
	return math.Pow(timeFact * randFloat(), 1 / (-1-beta))
}

func randFloat() float64 {
	return rand.Float64()
}
