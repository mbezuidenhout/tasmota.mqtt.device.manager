package tasmota

import (
	"encoding/json"
	"fmt"
	"regexp"
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
	UptimeSec  uint      `json:"UptimeSec"`
	LoadAvg    uint      `json:"LoadAvg"`
	Timezone   string    `json:"Timezone"`
	StatusNet  statusNet `json:"StatusNET"`
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
			d.SendCmnd("STATUS", "5")
			d.SendCmnd("TIMEZONE", "")
		} else {
			d.online = false
		}
	} else if strings.HasSuffix(msg.Topic(), "STATUS5") || strings.HasSuffix(msg.Topic(), "STATE") || strings.HasSuffix(msg.Topic(), "RESULT") {
		d.unmarshalPayload(msg.Payload())
	}
}

func (d *Device) SendCmnd(cmnd string, payload string) {
	mqttTopic := strings.Replace(strings.Replace(d.fullTopic, "%prefix%", "cmnd", 1), "%topic%", d.topic, 1) + cmnd
	d.mqttClient.Publish(mqttTopic, 1, false, payload)
}

func (d *Device) unmarshalPayload(payload []byte) error {
	// Append timezone to all date time strings
	r1 := regexp.MustCompile(`"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})"`)
	repl := r1.ReplaceAllString(string(payload), `"$1`+d.Timezone+`"`)
	return json.Unmarshal([]byte(repl), &d)
}
