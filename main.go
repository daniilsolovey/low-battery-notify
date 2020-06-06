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
	showingTime = "20000"
)

func main() {
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
		log.Warning("battery ", battery.State.String())

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

		time.Sleep(2 * time.Minute)
	}
}
