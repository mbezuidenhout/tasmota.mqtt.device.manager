package tasmota

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NewDevice will create a new Device
func NewDevice(topic, fullTopic string, mqttClient mqtt.Client) *Device {
	o := &Device{
		Topic:      topic,
		Fulltopic:  fullTopic,
		mqttClient: mqttClient,
	}
	subscribeTopics := make(map[string]byte)
	for _, t := range []string{"tele", "stat"} {
		subscribeTopics[strings.Replace(strings.Replace(fullTopic, "%prefix%", t, 1), "%topic%", topic, 1)+"+"] = 0
	}
	mqttClient.SubscribeMultiple(subscribeTopics, o.MessageHandler)
	return o
}

func (d *Device) MessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("%s: %s\n", msg.Topic(), msg.Payload())
	if strings.HasSuffix(msg.Topic(), "LWT") {
		if string(msg.Payload()) == "Online" {
			d.Online = true
			d.SendCmnd("STATUS", "5")
			d.SendCmnd("TIMEZONE", "")
			d.SendCmnd("STATUS", "2")
		} else {
			d.Online = false
		}
	} else if strings.HasSuffix(msg.Topic(), "STATUS5") || strings.HasSuffix(msg.Topic(), "STATUS2") || strings.HasSuffix(msg.Topic(), "STATE") || strings.HasSuffix(msg.Topic(), "RESULT") {
		err := d.unmarshalPayload(msg.Payload())
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (d *Device) SendCmnd(cmnd string, payload string) {
	mqttTopic := strings.Replace(strings.Replace(d.Fulltopic, "%prefix%", "cmnd", 1), "%topic%", d.Topic, 1) + cmnd
	d.mqttClient.Publish(mqttTopic, 1, false, payload)
}

func (d *Device) unmarshalPayload(payload []byte) error {
	// Append timezone to all date time strings
	r1 := regexp.MustCompile(`"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})"`)
	repl := r1.ReplaceAllString(string(payload), `"$1`+d.Timezone+`"`)
	return json.Unmarshal([]byte(repl), &d)
}
