package utils

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func GetHeartRateZones(age int) {
	// Info from https://www.rei.com/learn/expert-advice/how-to-train-with-a-heart-rate-monitor.html
	// Calculate max heart rate
	mmr := 220 - age

	// Calculate zones
	z1min := math.Floor(float64(mmr) * 0.55)
	z1max := math.Floor(float64(mmr) * 0.65)

	z2min := math.Floor(float64(mmr) * 0.65)
	z2max := math.Floor(float64(mmr) * 0.78)

	z3min := math.Floor(float64(mmr) * 0.78)
	z3max := math.Floor(float64(mmr) * 0.85)

	z4min := math.Floor(float64(mmr) * 0.85)
	z4max := math.Floor(float64(mmr) * 0.90)

	z5min := math.Floor(float64(mmr) * 0.90)
	z5max := math.Floor(float64(mmr) * 1)

	fmt.Println("Zone 1: ", z1min, "-", z1max)
	fmt.Println("Zone 2: ", z2min, "-", z2max)
	fmt.Println("Zone 3: ", z3min, "-", z3max)
	fmt.Println("Zone 4: ", z4min, "-", z4max)
	fmt.Println("Zone 5: ", z5min, "-", z5max)
}

func GetAgePrompt() int {
	fmt.Print("Enter age: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("An error occurred while reading input. Please try again", err)
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	ageConv, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		log.Fatal("An error occurred while parsing age. Please try again", err)
	}

	age := int(ageConv)

	return age
}
