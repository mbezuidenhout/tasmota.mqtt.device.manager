package tasmota

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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
	// Limit number of different sensors to 32
	o.Sensors = make(map[string]map[string]interface{})
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
			d.SendCmnd("Module", "")
			d.SendCmnd("DeviceName", "")
		} else {
			d.Online = false
		}
	} else if strings.HasSuffix(msg.Topic(), "STATUS5") || strings.HasSuffix(msg.Topic(), "STATUS2") || strings.HasSuffix(msg.Topic(), "STATE") || strings.HasSuffix(msg.Topic(), "RESULT") {
		err := d.unmarshalDevicePayload(msg.Payload())
		if err != nil {
			fmt.Println(err)
		}
	} else if strings.HasSuffix(msg.Topic(), "SENSOR") {
		err := d.unmarshalSensorPayload(msg.Payload())
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (d *Device) SendCmnd(cmnd string, payload string) {
	mqttTopic := strings.Replace(strings.Replace(d.Fulltopic, "%prefix%", "cmnd", 1), "%topic%", d.Topic, 1) + cmnd
	d.mqttClient.Publish(mqttTopic, 1, false, payload)
}

func (d *Device) unmarshalDevicePayload(payload []byte) error {
	// Append timezone to all date time strings
	r1 := regexp.MustCompile(`"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})"`)
	repl := r1.ReplaceAllString(string(payload), `"$1`+d.Timezone+`"`)

	// Change module string to only have the name.
	r2 := regexp.MustCompile(`{\"Module\":{\"\d+\":\"([^"]+)\"}}`)
	repl = r2.ReplaceAllString(string(repl), `{"Module":"$1"}`)

	return json.Unmarshal([]byte(repl), &d)
}

func (d *Device) unmarshalSensorPayload(payload []byte) error {
	r3 := regexp.MustCompile(`{\"ZbInfo\":{\"([x0-9A-F]+)\":(.*)}}`)
	matches := r3.FindAllStringSubmatch(string(payload), -1)
	if len(matches) == 1 && len(matches[0]) == 3 {
		if len(d.Sensors["Zigbee"]) == 0 {
			d.Sensors["Zigbee"] = make(map[string]interface{})
		}
		x := ZigbeeTH01{}
		err := json.Unmarshal([]byte(matches[0][2]), &x)
		if err == nil {
			d.Sensors["Zigbee"][matches[0][1]] = x
		}
		return err
	} else {
		return nil
	}

}

func (d *Device) SetName(name string) {
	fmt.Printf("Setting device %s name to %s\n", d.Topic, name)
	d.SendCmnd("DeviceName", name)
}

func (d *Device) SetTimezone(timezone int) {
	fmt.Printf("Setting device %s timezone to %d\n", d.Topic, timezone)
	d.SendCmnd("DeviceName", strconv.Itoa(timezone))
}

func (d *Device) GetSensorTypes() []string {
	var keys = []string{}
	for k := range d.Sensors {
		keys = append(keys, k)
	}
	return keys
}

func (d *Device) GetSensor(sensorType string) map[string]interface{} {
	return d.Sensors[sensorType]
}
