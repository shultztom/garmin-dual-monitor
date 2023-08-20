package main

import (
	"garmin-dual-monitor/bt"
	"garmin-dual-monitor/utils"
)

func main() {
	age := utils.GetAgePrompt()

	bt.StartHrm(age)
}
