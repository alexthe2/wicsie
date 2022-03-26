package main

import (
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/png"
	"log"
	"math/rand"
	"time"
	"wicsie/agents"
	"wicsie/drawing"
	"wicsie/heatMapDecoder"
	"wicsie/simulation"
)

func main() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	rand.Seed(time.Now().UnixNano())

	heatMap, heatChunkMap, colorMap, width, height := heatMapDecoder.LoadAndDecode("population.png")
	mask, err := gg.LoadImage("europe.png")
	if err != nil {
		log.Fatalf("Could not load mask: %v", err)
	}

	fmt.Printf("%v\n", colorMap)
	legend := heatMapDecoder.ReadPredefined()

	appendix := flag.String("appendix", "", "the appendix in which the pictures should be saved, outAppendix")
	flag.Parse()

	const steps = 1000

	grid := agents.CreateGridMap(width, height, 3)

	createMovement := func() agents.Movement {
		return agents.CreateGridMovement(100, grid, heatChunkMap)
	}

	simu := simulation.CreateSimulation(simulation.Config{
		Weight:    .2,
		Width:     float64(width),
		Height:    float64(height),
		Movement:  createMovement,
		Spreading: agents.CreateNoSpreading(),

		HeatMap:     heatMap,
		LegendIndex: legend,
	})

	simu.InitInfect(0.01)
	board := drawing.CreateBoard(width, height, mask, 1)

	for i := 0; i < steps; i++ {
		grid.UpdateGridMap(simu.GetAgents())
		board.DrawGridMap(*grid)
		board.SaveBoard(fmt.Sprintf("out%s/boardgrid%d.png", *appendix, i))

		simu.Step()
		simu.DrawToBoard(board)
		board.SaveBoard(fmt.Sprintf("out%s/board%d.png", *appendix, i))
	}

}
