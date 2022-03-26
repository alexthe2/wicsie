package heatMapDecoder

import (
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"wicsie/constants"
)

func Decode(file io.Reader) ([][]int, [][]int, map[int]color.Color, int, int) {
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

	heatChunkMap := make([][]int, (height/constants.KChunkSize)*10)
	for i := 0; i < (height/constants.KChunkSize)*10; i++ {
		heatChunkMap[i] = make([]int, (width/constants.KChunkSize)*10)
	}
	//fmt.Println(len(heatChunkMap))
	//fmt.Println((height / constants.KChunkSize))
	//fmt.Println((width / constants.KChunkSize))

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

	for i := 0; i < height; i += constants.KChunkSize {
		for j := 0; j < width; j += constants.KChunkSize {
			counter := 0
			for x := i; x < (i + constants.KChunkSize); x++ {
				for y := j; y < (j + constants.KChunkSize); y++ {
					clr := img.At(x, y)
					if knownColors[clr] == 0 {
						counter++
					}
					//percentage := float64(counter) / float64(constants.KChunkSize*constants.KChunkSize)
					//fmt.Println("i/constants.KChunkSize: ", i/constants.KChunkSize)
					//fmt.Println("j/constants.KChunkSize: ", j/constants.KChunkSize)
					if counter > 0 /*float64(percentage) > 0.1*/ {
						heatChunkMap[i/constants.KChunkSize][j/constants.KChunkSize] = 1

					} else {
						heatChunkMap[i/constants.KChunkSize][j/constants.KChunkSize] = 0
					}
				}
			}
		}
	}

	return heatMap, heatChunkMap, colorMap, width, height
}

func LoadAndDecode(file string) ([][]int, [][]int, map[int]color.Color, int, int) {
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
