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
		heading := int16(0)
		direction := int16(-1)
		gobot.Every(100*time.Millisecond, func() {
			r2q5.Dome(heading)
			if (heading < -160) {
				heading = -160
			} else if (heading == -160) {
				direction *= -1
				heading += direction
			} else if (heading > 180) {
				heading = 180
			} else if (heading == 180) {
				direction *= -1
				heading += direction
			} else {
				heading += direction
			}
		})
	}

	robot := gobot.NewRobot("R2Q5",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{r2q5},
		work,
	)

	robot.Start()
}
