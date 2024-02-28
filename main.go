package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"os"
	"time"
)

const (
	sampleRate   beep.SampleRate = 44100
)

const  size  int = 100 
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
	// Print the coordinates
	
	// var space = [100][2]float64{};
	
	// for i, coord := range coordinates {
		// fmt.Printf("[%d] Time: %.2f, Intensity: %.2f\n", i, coord[0]/float64(sampleRate), coord[1])
		// space.append(coord[0], coord[1])
		// fmt.Printf("{},{}",coord[0],coord[1])

	// }
	// for i := 0; i < 10; i++ {
	// 	fmt.Printf("*************************")
	// }

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
		// fmt.Println(x,y)
		if x >= 0 && x < size && y >= 0 && y < size {
			space[y][x] = '*'
		}
	}
	//Print space upside down
	for i := len(space) - 1; i >= 0; i-- {
		for j := range space[i] {
			if(i<3){
				fmt.Printf("\x1b[34m%c\x1b[0m", space[i][j])
			}
			if(i>=3 && i<=4){
				fmt.Printf("\x1b[32m%c\x1b[0m", space[i][j])
			}
			if(i>4){
				fmt.Printf("\x1b[31m%c\x1b[0m", space[i][j])
			}
			// fmt.Printf("\x1b[31m%c\x1b[0m", space[i][j])
		}
		fmt.Println()
	}
}
