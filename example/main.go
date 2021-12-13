package main

import (
	"fmt"

	tasmota "github.com/mbezuidenhout/tasmota.mqtt.device.manager"
)

func main() {
	device := tasmota.Device{"name", "ip"}

	fmt.Println(device.GetUptime)
}
