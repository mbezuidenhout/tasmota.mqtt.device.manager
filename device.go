package tasmota

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Device struct {
	Topic     string
	topic     string
	fullTopic string
}

// NewDevice will create a new Device
func NewDevice(topic, fullTopic string) *Device {
	o := &Device{
		topic:     topic,
		fullTopic: fullTopic,
	}
	return o
}

func (d *Device) MessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("%s: %s\n", msg.Topic(), msg.Payload())
}
