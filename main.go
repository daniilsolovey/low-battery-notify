package main

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
	"time"

	"github.com/distatus/battery"
	"github.com/reconquest/pkg/log"
)

const (
	showingTime = "500"
	sleepTime   = 5 * time.Second
)

var PERFORMANCE_MODE string

func main() {
	batteryMode := false
	powerMode := false
	for {
		// battery data:
		var batteryStatus string
		var batteryPercent float64
		battery, err := battery.Get(0)
		if err != nil {
			log.Error(err)
			batteryStatus = "battery: error"
		} else {
			batteryPercent = math.Floor(battery.Current / battery.Full * 100)
			batteryStatus = "low battery!!!!!! " + fmt.Sprintf(
				"%.0f", batteryPercent,
			) +
				" %"
		}

		if int(batteryPercent) <= 15 && battery.State.String() == "Discharging" {
			var info []string
			info = append(info, batteryStatus)
			notify := exec.Command(
				"notify-send", "-t",
				showingTime, "info",
				strings.Join(info, "\n"),
			)
			err = notify.Run()
			if err != nil {
				log.Error(err)
			}
		}

		if battery.State.String() == "Discharging" && batteryMode == false {
			//set battery mode
			err := setBatteryMode()
			if err != nil {
				log.Error(err)
			}
			log.Info("set battery mode")
			batteryMode = true
			powerMode = false
		}

		if battery.State.String() != "Discharging" && powerMode == false {
			//set power mode
			err := setPowerMode()
			if err != nil {
				log.Error(err)
			}
			log.Info("set power mode")
			powerMode = true
			batteryMode = false
		}

		time.Sleep(sleepTime)
	}
}

func setPowerMode() error {
	setMaxFrequency := exec.Command(
		"/bin/sh", "-c", "sudo cpupower frequency-set --max 5.2GHz",
	)
	setMinFrequency := exec.Command(
		"/bin/sh", "-c", "sudo cpupower frequency-set --min 4.4GHz",
	)

	err := setMinFrequency.Run()
	if err != nil {
		return err
	}

	err = setMaxFrequency.Run()
	if err != nil {
		return err
	}

	return nil
}

func setBatteryMode() error {
	setMaxFrequencyForDischarging := exec.Command(
		"/bin/sh", "-c", "sudo cpupower frequency-set --max 3.2GHz",
	)

	err := setMaxFrequencyForDischarging.Run()
	if err != nil {
		return err
	}

	return nil
}
