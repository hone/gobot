// +build example
//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth address or name as the first param:

	go run examples/r2q5.go R2-1234

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero/r2q5"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	r2q5 := r2q5.NewDriver(bleAdaptor)

	work := func() {
		gobot.Every(3*time.Second, func() {
			r2q5.Macro(uint8(gobot.Rand(55)))
		})
	}

	robot := gobot.NewRobot("R2Q5",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{r2q5},
		work,
	)

	robot.Start()
}
