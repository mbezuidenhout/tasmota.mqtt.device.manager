package tasmota

import (
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type statusNet struct {
	Hostname   string  `json:"Hostname"`
	IPAddress  string  `json:"IPAddress"`
	Gateway    string  `json:"Gateway"`
	Subnetmask string  `json:"Subnetmask"`
	DNSServer1 string  `json:"DNSServer1"`
	DNSServer2 string  `json:"DNSServer2"`
	Mac        string  `json:"Mac"`
	Webserver  int     `json:"Webserver"`
	WifiPower  float32 `json:"WifiPower"`
}

type Device struct {
	topic      string
	fullTopic  string
	mqttClient mqtt.Client
	online     bool
	StatusNet  statusNet `json:"StatusNET"`
	Uptime     string    `json:"Uptime"`
}

// NewDevice will create a new Device
func NewDevice(topic, fullTopic string, mqttClient mqtt.Client) *Device {
	o := &Device{
		topic:      topic,
		fullTopic:  fullTopic,
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
			d.online = true
			d.GetStatusNet()
		} else {
			d.online = false
		}
	} else if strings.HasSuffix(msg.Topic(), "STATUS5") || strings.HasSuffix(msg.Topic(), "STATE") {
		d.unmarshalPayload(msg.Payload())
	}
}

func (d *Device) GetStatusNet() {
	status := strings.Replace(strings.Replace(d.fullTopic, "%prefix%", "cmnd", 1), "%topic%", d.topic, 1) + "STATUS"
	d.mqttClient.Publish(status, 1, false, "5")
}

func (d *Device) unmarshalPayload(payload []byte) {
	json.Unmarshal(payload, d)
}
