package tasmota

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func (d *Device) unmarshalSensorPayload(payload []byte) error {
	r1 := regexp.MustCompile(`{\"ZbInfo\":{\"([x0-9A-F]+)\":(.*)}}`)
	matches1 := r1.FindAllStringSubmatch(string(payload), -1)
	if len(matches1) == 1 && len(matches1[0]) == 3 {
		if len(d.Sensors["Zigbee"]) == 0 {
			d.Sensors["Zigbee"] = make(map[string]interface{})
		}
		x := ZigbeeTH01{}
		err := json.Unmarshal([]byte(matches1[0][2]), &x)
		if err == nil {
			d.Sensors["Zigbee"][matches1[0][1]] = x
		}
		return err
	}
	r2 := regexp.MustCompile(`{\"ZbReceived\":{\"([x0-9A-F]+)\":(.*)}}`)
	matches2 := r2.FindAllStringSubmatch(string(payload), -1)
	if len(matches2) == 1 && len(matches2[0]) == 3 {
		if len(d.Sensors["Zigbee"]) == 0 {
			d.Sensors["Zigbee"] = make(map[string]interface{})
		}
		x := ZigbeeTH01{}
		y := d.Sensors["Zigbee"][matches2[0][1]].(ZigbeeTH01)
		err := json.Unmarshal([]byte(matches2[0][2]), &x)
		y.LastSeen = 0
		if x.Name != y.Name {
			y.Name = x.Name
		}
		if x.Temperature != 0 {
			y.Temperature = x.Temperature
		}
		if x.Humidity != 0 {
			y.Humidity = x.Humidity
		}
		if x.BatteryPercentage != 0 {
			y.BatteryPercentage = x.BatteryPercentage
		}
		if x.LinkQuality != y.LinkQuality {
			y.LinkQuality = x.LinkQuality
		}
		d.Sensors["Zigbee"][matches2[0][1]] = y
		fmt.Println(d.Sensors["Zigbee"][matches2[0][1]].(ZigbeeTH01))
		return err
	}
	return nil
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
