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

type Firmware struct {
	Version       string `json:"Version"`
	BuildDateTime string `json:"BuildDateTime"`
	Boot          int    `json:"Boot"`
	Hardware      string `json:"Hardware"`
	CpuFrequency  int    `json:"CpuFrequency"`
}