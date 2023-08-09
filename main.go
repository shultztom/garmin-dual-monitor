package main

import (
	"bufio"
	"fmt"
	"github.com/inancgumus/screen"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"tinygo.org/x/bluetooth"
)

//func connectAddress() string {
//	uuid := "66ebfb9a-bfbd-94c9-1c21-6baa0eec41f0"
//	return uuid
//}

var (
	adapter = bluetooth.DefaultAdapter

	heartRateServiceUUID        = bluetooth.ServiceUUIDHeartRate
	heartRateCharacteristicUUID = bluetooth.CharacteristicUUIDHeartRateMeasurement

	age = 0
)

func getHeartRateZones() {
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

func main() {
	// Info from https://www.rei.com/learn/expert-advice/how-to-train-with-a-heart-rate-monitor.html
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

	age = int(ageConv)

	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	ch := make(chan bluetooth.ScanResult, 1)

	// Start scanning.
	println("scanning...")
	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		println("found device:", result.Address.String(), result.RSSI, result.LocalName())
		if strings.Contains(result.LocalName(), "HRM-Dual") {
			adapter.StopScan()
			ch <- result
		}
	})

	var device *bluetooth.Device
	select {
	case result := <-ch:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		if err != nil {
			println(err.Error())
			return
		}

		println("connected to ", result.Address.String())
	}

	// get services
	println("discovering services/characteristics")
	srvcs, err := device.DiscoverServices([]bluetooth.UUID{heartRateServiceUUID})
	must("discover services", err)

	if len(srvcs) == 0 {
		panic("could not find heart rate service")
	}

	srvc := srvcs[0]

	println("found service", srvc.UUID().String())

	chars, err := srvc.DiscoverCharacteristics([]bluetooth.UUID{heartRateCharacteristicUUID})
	if err != nil {
		println(err)
	}

	if len(chars) == 0 {
		panic("could not find heart rate characteristic")
	}

	char := chars[0]
	println("found characteristic", char.UUID().String())

	char.EnableNotifications(func(buf []byte) {
		screen.Clear()
		getHeartRateZones()
		fmt.Println()
		// start listening for updates and render
		fmt.Printf("Heart Rate %d\n\n", uint8(buf[1]))

	})

	select {}

}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
