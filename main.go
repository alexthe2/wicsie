package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
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

	heatMap, colorMap, width, height := heatMapDecoder.LoadAndDecode("population.png")
	fmt.Printf("%v\n", colorMap)
	legend := heatMapDecoder.ReadPredefined()

	appendix := flag.String("appendix", "", "the appendix in which the pictures should be saved, outAppendix")
	flag.Parse()

	const steps = 1000

	createMovement := func() agents.Movement {
		return agents.CreateRandomMovement(10)
	}

	simu := simulation.CreateSimulation(simulation.Config{
		Weight:    1,
		Width:     float64(width),
		Height:    float64(height),
		Movement:  createMovement,
		Spreading: agents.CreateOnTouchSpreading(),

		HeatMap:     heatMap,
		LegendIndex: legend,
	})

	simu.InitInfect(0.01)
	board := drawing.CreateBoard(width, height)
	for i := 0; i < steps; i++ {
		simu.Step()
		simu.DrawToBoard(board)
		board.SaveBoard(fmt.Sprintf("out%s/board%d.png", *appendix, i))
	}

}
