package main

import (
	"fmt"
	"os"
	"time"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

const (
	sampleRate beep.SampleRate = 44100
	size       int             = 100
)

func main() {
	// Load the audio file
	fmt.Print("Enter your file name (.wav): ")
	var name string
	fmt.Scan(&name)

	file, err := os.Open(name)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the audio file
	streamer, format, err := wav.Decode(file)
	if err != nil {
		fmt.Println("Error decoding audio file:", err)
		return
	}

	// Initialize the speaker with the audio format
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Create a 2D array to store the coordinates
	var coordinates [][]float64

	// Fill the coordinates array
	fillCoordinates(streamer, &coordinates)

	// Print the coordinates
	printCoordinates(coordinates)
}

func fillCoordinates(streamer beep.StreamSeeker, coordinates *[][]float64) {
	// Calculate the total number of samples
	numSamples := int(streamer.Len())

	// Read all samples and store the coordinates
	for i := 0; i < numSamples; i++ {
		// Read a sample of audio data
		sample := make([][2]float64, 1)
		streamer.Stream(sample)

		// Calculate the intensity for the sample
		intensity := calculateIntensity(sample[0])

		// Append the intensity to the coordinates array
		*coordinates = append(*coordinates, []float64{float64(i), intensity})
	}
}

func calculateIntensity(sample [2]float64) float64 {
	// Calculate intensity based on the audio sample
	return sample[0] + sample[1]
}

func printCoordinates(coordinates [][]float64) {
	// Create a 2D array 'space' representing the graph
	var space [size][size]rune // Assuming a fixed size of 100x100

	for i := range space {
		for j := range space[i] {
			space[i][j] = ' ' // Initialize all points as empty
		}
	}

	// Plot each point from coordinates onto the graph
	for _, coord := range coordinates {
		x, y := int(coord[0]/1000), int(coord[1]*10)
		if x >= 0 && x < size && y >= 0 && y < size {
			space[y][x] = '*'
		}
	}
	// Print space upside down
	for i := len(space) - 1; i >= 0; i-- {
		for j := range space[i] {
			if i < 3 {
				fmt.Printf("\x1b[34m%c\x1b[0m", space[i][j])
			}
			if i >= 3 && i <= 4 {
				fmt.Printf("\x1b[32m%c\x1b[0m", space[i][j])
			}
			if i > 4 {
				fmt.Printf("\x1b[31m%c\x1b[0m", space[i][j])
			}
		}
		fmt.Println()
	}

	// Plot using Gonum
	plotGraph(space)
}
func plotGraph(space [size][size]rune) {
	// Create a plot
	p := plot.New()
	// if err != nil {
	// 	fmt.Println("Error creating plot:", err)
	// 	return
	// }

	// Create a scatter plotter
	scatter, err := plotter.NewScatter(getXYs(space))
	if err != nil {
		fmt.Println("Error creating scatter plotter:", err)
		return
	}

	// Add the scatter plotter to the plot
	p.Add(scatter)

	// Save the plot to a file
	if err := p.Save(4, 4, "output.png"); err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}

	fmt.Println("Plot saved as output.png. Open the file with an image viewer.")
}

func getXYs(space [size][size]rune) plotter.XYs {
	var XYs plotter.XYs

	for i, row := range space {
		for j, value := range row {
			if value == '*' {
				XYs = append(XYs, plotter.XY{X: float64(j*10), Y: float64((size - 1 - i)*10)})
			}
		}
	}

	return XYs
}
