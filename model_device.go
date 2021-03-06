/*
 * Tasmota Device Manager API
 *
 * Device manager for Tasmota devices via MQTT [Source](https://github.com/mbezuidenhout/tdm).
 *
 * API version: 0.1.0
 * Contact: marius.bezuidenhout@gmail.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package tasmota

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Device struct {
	// Unique device topic
	Topic string `json:"Topic"`
	// Full topic format
	Fulltopic string `json:"FullTopic"`
	Name      string `json:"DeviceName"`
	Module    string `json:"Module"`
	// Device status
	Online     bool        `json:"Online"`
	mqttClient mqtt.Client `json:"-"`
	UptimeSec  uint        `json:"UptimeSec,omitempty"`
	LoadAvg    uint        `json:"LoadAvg,omitempty"`
	Timezone   string      `json:"Timezone,omitempty"`
	Wifi       Wifi        `json:"Wifi,omitempty"`
	Network    Network     `json:"StatusNET,omitempty"`
	Firmware   Firmware    `json:"StatusFWR,omitempty"`
	// Sensors are first added to a group and then by name or number.
	// If you have multiple Si7021 temperature sensors they should be like
	// Sensors["Si7021"]["0"] = {"Temperature": 10.0 }
	// Sensors["Si7021"]["1"] = {"Temperature": 12.3 }
	Sensors map[string]map[string]interface{} `json:"-"`
}
