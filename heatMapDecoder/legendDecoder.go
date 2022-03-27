package heatMapDecoder

import (
	"image/color"
)

type LegendIndex struct {
	Clr        color.Color
	Multiplier int
	Id         int
}

func ReadPredefined() []LegendIndex {
	return []LegendIndex{
		LegendIndex{Clr: color.RGBA{G: 255, B: 255, A: 255}},
		LegendIndex{Clr: color.RGBA{R: 255, G: 255, B: 255, A: 255}, Multiplier: 40, Id: 1},
		LegendIndex{Clr: color.RGBA{R: 255, B: 255, A: 255}, Multiplier: 60, Id: 2},
		LegendIndex{Clr: color.RGBA{R: 255, A: 255}, Multiplier: 120, Id: 3},
		LegendIndex{Clr: color.RGBA{R: 255, G: 255, A: 255}, Multiplier: 70, Id: 4},
		LegendIndex{Clr: color.RGBA{B: 255, A: 255}, Multiplier: 200, Id: 5},
		LegendIndex{Clr: color.RGBA{A: 255}, Multiplier: 800, Id: 6},
	}
}

func GetMultiplier(indexes []LegendIndex, id int) int {
	for _, indx := range indexes {
		if indx.Id == id {
			return indx.Multiplier
		}
	}
	return 0
}
