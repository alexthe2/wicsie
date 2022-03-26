package heatMapDecoder

import (
	"image"
	"image/color"
	"io"
	"log"
	"os"
)

func Decode(file io.Reader) ([][]int, map[int]color.Color, int, int) {
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding image: %v", err)
	}
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	heatMap := make([][]int, height)
	for i := 0; i < height; i++ {
		heatMap[i] = make([]int, width)
	}

	colorMap := make(map[int]color.Color)
	//[1,2,4,4,4]
	//[5,3,1,1,3]
	knownColors := make(map[color.Color]int)
	index := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			clr := img.At(x, y)
			//clr --> [R,G,B,A]
			if _, ok := knownColors[clr]; !ok {
				knownColors[clr] = index
				colorMap[index] = clr
				index++
			}

			heatMap[y][x] = knownColors[clr]
		}
	}

	return heatMap, colorMap, width, height
}

func LoadAndDecode(file string) ([][]int, map[int]color.Color, int, int) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(f)
	return Decode(f)
}
