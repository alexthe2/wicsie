package simulation

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"wicsie/agents"
)

type DataWriter struct {
	fileName string
	LOGGER   *log.Logger
}

func CreateDataWriter(fileName string) *DataWriter {
	file, err := os.Create(fileName)
	logger := log.New(os.Stdout, "[DW] ", log.Ltime)
	if err != nil {
		logger.Fatalf("Could not create file %s: %s", fileName, err)
	}
	writer := csv.NewWriter(file)
	err = writer.Write([]string{"Step", "Healthy", "Infected", "Cured"})
	if err != nil {
		logger.Fatalf("Could not write header to file %s: %s", fileName, err)
	}

	writer.Flush()
	err = file.Close()
	if err != nil {
		logger.Fatalf("Could not close file %s: %s", fileName, err)
	}

	return &DataWriter{fileName, logger}
}
func (d *DataWriter) Write(step int, grid agents.GridMap) {
	file, err := os.OpenFile(d.fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		d.LOGGER.Fatalf("Could not open file %s: %s", d.fileName, err)
	}
	writer := csv.NewWriter(file)
	err = writer.Write([]string{
		strconv.Itoa(step),
		strconv.Itoa(grid.Stats[0]),
		strconv.Itoa(grid.Stats[1]),
		strconv.Itoa(grid.Stats[2]),
	})
	if err != nil {
		d.LOGGER.Fatalf("Could not write to file %s: %s", d.fileName, err)
	}
	writer.Flush()
	err = file.Close()
	if err != nil {
		d.LOGGER.Fatalf("Could not close file %s: %s", d.fileName, err)
	}
}
