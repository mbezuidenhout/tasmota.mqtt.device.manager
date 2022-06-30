package tasmota

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/mitchellh/mapstructure"
)

func (d *Device) unmarshalSensorPayload(payload []byte) error {
	if len(d.Sensors["Zigbee"]) == 0 {
		d.Sensors["Zigbee"] = make(map[string]interface{})
	}
	r1 := regexp.MustCompile(`{\"ZbInfo\":{\"([x0-9A-F]+)\":(.*)}}`)
	matches1 := r1.FindAllStringSubmatch(string(payload), -1)
	if len(matches1) == 1 && len(matches1[0]) == 3 {
		var err error
		sensorMap := make(map[string](interface{}))
		err = json.Unmarshal([]byte(matches1[0][2]), &sensorMap)
		if err == nil {
			switch model := sensorMap["ModelId"]; model {
			case "TH01":
				x := ZigbeeTH01{}
				mapstructure.Decode(sensorMap, &x)
				d.Sensors["Zigbee"][matches1[0][1]] = x
			case "DS01":
				x := ZigbeeDS01{}
				mapstructure.Decode(sensorMap, &x)
				d.Sensors["Zigbee"][matches1[0][1]] = x
			case "WB01":
				x := ZigbeeWB01{}
				mapstructure.Decode(sensorMap, &x)
				d.Sensors["Zigbee"][matches1[0][1]] = x
			case "MS01":
				x := ZigbeeMS01{}
				mapstructure.Decode(sensorMap, &x)
				d.Sensors["Zigbee"][matches1[0][1]] = x
			default:
				d.Sensors["Zigbee"][matches1[0][1]] = sensorMap
			}
		}
		return err
	}
	r2 := regexp.MustCompile(`{\"ZbReceived\":{\"([x0-9A-F]+)\":(.*)}}`)
	matches2 := r2.FindAllStringSubmatch(string(payload), -1)
	if len(matches2) == 1 && len(matches2[0]) == 3 && d.Sensors["Zigbee"][matches2[0][1]] != nil {
		var err error
		sensorMap := make(map[string](interface{}))
		err = json.Unmarshal([]byte(matches2[0][2]), &sensorMap)
		if err == nil {
			if sensor, ok := d.Sensors["Zigbee"][matches2[0][1]]; ok { // Sensor exists in array
				switch sensor.(type) {
				case ZigbeeTH01:
					y := d.Sensors["Zigbee"][matches2[0][1]].(ZigbeeTH01)
					y.LastSeen = 0
					y.LastSeenEpoch = time.Now().UTC().Unix()
					if sensorMap["Temperature"] != nil {
						y.Temperature = float32(sensorMap["Temperature"].(float64))
					}
					if sensorMap["Humidity"] != nil {
						y.Humidity = float32(sensorMap["Humidity"].(float64))
					}
					if sensorMap["BatteryPercentage"] != nil {
						y.BatteryPercentage = int(sensorMap["BatteryPercentage"].(float64))
					}
					if sensorMap["LinkQuality"] != nil {
						y.LinkQuality = int(sensorMap["LinkQuality"].(float64))
					}
					d.Sensors["Zigbee"][matches2[0][1]] = y
				case ZigbeeWB01:
					y := d.Sensors["Zigbee"][matches2[0][1]].(ZigbeeWB01)
					y.LastSeen = 0
					y.LastSeenEpoch = time.Now().UTC().Unix()
					if sensorMap["BatteryPercentage"] != nil {
						y.BatteryPercentage = int(sensorMap["BatteryPercentage"].(float64))
					}
					if sensorMap["LinkQuality"] != nil {
						y.LinkQuality = int(sensorMap["LinkQuality"].(float64))
					}
					if sensorMap["Power"] != nil {
						y.Power = int(sensorMap["Power"].(float64))
					}
				case ZigbeeMS01:
					y := d.Sensors["Zigbee"][matches2[0][1]].(ZigbeeMS01)
					y.LastSeen = 0
					y.LastSeenEpoch = time.Now().UTC().Unix()
					if sensorMap["BatteryPercentage"] != nil {
						y.BatteryPercentage = int(sensorMap["BatteryPercentage"].(float64))
					}
					if sensorMap["LinkQuality"] != nil {
						y.LinkQuality = int(sensorMap["LinkQuality"].(float64))
					}
					if sensorMap["Occupancy"] != nil {
						y.Occupancy = int(sensorMap["Occupancy"].(float64))
					}
				case ZigbeeDS01:
					y := d.Sensors["Zigbee"][matches2[0][1]].(ZigbeeDS01)
					y.LastSeen = 0
					y.LastSeenEpoch = time.Now().UTC().Unix()
					if sensorMap["BatteryPercentage"] != nil {
						y.BatteryPercentage = int(sensorMap["BatteryPercentage"].(float64))
					}
					if sensorMap["LinkQuality"] != nil {
						y.LinkQuality = int(sensorMap["LinkQuality"].(float64))
					}
					if sensorMap["Contact"] != nil {
						y.Contact = int(sensorMap["Contact"].(float64))
					}
				}
				/*
					y.LastSeen = 0
					y.LastSeenEpoch = time.Now().UTC().Unix()
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
				*/
				return err
			} else { // First time this sensor is seen
				d.SendCmnd("ZbInfo", "")
			}
		}
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
