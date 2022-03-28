package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"gopkg.in/yaml.v3"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
	"wicsie/agents"
	"wicsie/constants"
	"wicsie/drawing"
	"wicsie/heatMapDecoder"
	"wicsie/simulation"
)

type config struct {
	PopulationMap string  `yaml:"populationMap"`
	MaskMap       string  `yaml:"maskMap"`
	Behaviour     string  `yaml:"behaviour"`
	Steps         int     `yaml:"steps"`
	Weight        float64 `yaml:"weight"`
}

type spreadingPoint struct {
	X           int     `yaml:"x"`
	Y           int     `yaml:"y"`
	Probability float64 `yaml:"probability"`
}

func (c *config) readConfig() {
	file, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		log.Fatalf("[Config] Error reading config file: %v", err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Fatalf("[Config] Error parsing config file: %v", err)
	}
}

func readProbabilities() []*spreadingPoint {
	file, err := ioutil.ReadFile("config/diseaseStart.yml")
	if err != nil {
		log.Fatalf("[Config] Error reading probabilities file: %v", err)
	}
	var probabilities []*spreadingPoint
	err = yaml.Unmarshal(file, &probabilities)
	if err != nil {
		log.Fatalf("[Config] Error parsing probabilities file: %v", err)
	}
	return probabilities
}

func initSystem() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	rand.Seed(time.Now().UnixNano())
}

func createSimulation(cfg config) (*simulation.Simulation, *agents.GridMap, int, int) {
	heatMap, _, _, width, height := heatMapDecoder.LoadAndDecode(fmt.Sprintf("config/%s", cfg.PopulationMap))
	legend := heatMapDecoder.ReadPredefined()
	grid := agents.CreateGridMap(width, height, constants.KChunkSize)
	script := agents.DecodeFile(fmt.Sprintf("config/%s", cfg.Behaviour))

	movementCreation := func() agents.Movement {
		return agents.CreateScriptMovement(grid, script, heatMap)
	}

	return simulation.CreateSimulation(simulation.Config{
		Weight:      cfg.Weight,
		Width:       float64(width),
		Height:      float64(height),
		Movement:    movementCreation,
		Spreading:   agents.CreateGridSpread(grid),
		HeatMap:     heatMap,
		LegendIndex: legend,
	}), grid, width, height
}

func preInfectSystem(simulation *simulation.Simulation) {
	points := readProbabilities()
	for _, point := range points {
		simulation.InfectAtPosition(float64(point.X), float64(point.Y), point.Probability)
	}
}

func runSimulation(simulation *simulation.Simulation, board *drawing.Board, grid *agents.GridMap, config config) {
	for i := 0; i < config.Steps; i++ {
		simulation.Step()
		grid.UpdateGridMap(simulation.GetAgents())

		board.DrawGridMap(*grid)
		board.SaveBoard(fmt.Sprintf("out/raw/boardgrid%d.png", i))
	}
}

func main() {
	var config config
	config.readConfig()

	initSystem()

	mask, err := gg.LoadImage(fmt.Sprintf("config/%s", config.MaskMap))
	if err != nil {
		log.Fatalf("[Config] Could not load mask: %v", err)
	}

	simu, grid, w, h := createSimulation(config)
	preInfectSystem(simu)
	board := drawing.CreateBoard(w, h, mask, 1)
	runSimulation(simu, board, grid, config)
}
