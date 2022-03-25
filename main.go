package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
	agents2 "wicsie/agents"
	"wicsie/drawing"
	"wicsie/simulation"
)

func main() {
	appendix := flag.String("appendix", "", "the appendix in which the pictures should be saved, outAppendix")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	width := 2000
	height := 1100

	const agents = 2000
	const steps = 2000

	createMovement := func() agents2.Movement {
		return agents2.CreateRandomMovement(300)
	}

	simu := simulation.CreateSimulation(agents, width, height, createMovement, agents2.CreateOnTouchSpreading())
	simu.InitInfect(0.01)
	board := drawing.CreateBoard(width, height)

	for i := 0; i < steps; i++ {
		simu.Step()
		simu.DrawToBoard(board)
		board.SaveBoard(fmt.Sprintf("out%s/board%d.png", *appendix, i))
	}
}
