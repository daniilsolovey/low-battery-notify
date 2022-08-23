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
	SHOWING_TIME         = "500"
	SLEEP_TIME           = 5 * time.Second
	LOW_BATTERY_LEVEL    = 15
	MEDIUM_BATTERY_LEVEL = 50
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
		if err != nil && !strings.Contains(fmt.Sprint(err), "State:Invalid state `Not charging") {
			log.Error(err)
			batteryStatus = "battery: error"
		} else {
			batteryPercent = math.Floor(battery.Current / battery.Full * 100)
			batteryStatus = "low battery!!!!!! " + fmt.Sprintf(
				"%.0f", batteryPercent,
			) +
				" %"
		}

		if int(batteryPercent) <= LOW_BATTERY_LEVEL &&
			battery.State.String() == "Discharging" {
			var info []string
			info = append(info, batteryStatus)
			notify := exec.Command(
				"notify-send", "-t",
				SHOWING_TIME, "info",
				strings.Join(info, "\n"),
			)
			err = notify.Run()
			if err != nil {
				log.Error(err)
			}
		}

		if battery.State.String() == "Discharging" && !batteryMode {
			//set min power mode (for battery)
			err := setBatteryMode()
			if err != nil {
				log.Error(err)
			}

			log.Info("set min power mode (for battery)")
			batteryMode = true
			powerMode = false
		}

		if battery.State.String() != "Discharging" &&
			!powerMode &&
			int(batteryPercent) > MEDIUM_BATTERY_LEVEL {
			//set max power mode
			err := setPowerMode()
			if err != nil {
				log.Error(err)
			}

			log.Info("set max power mode")
			powerMode = true
			batteryMode = false
		}

		time.Sleep(SLEEP_TIME)
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
