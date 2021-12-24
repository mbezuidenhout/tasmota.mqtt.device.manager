module example

go 1.17

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/mbezuidenhout/tasmota.mqtt.device.manager/v2 v2.0.0-20211221074425-3b2da9af68cb
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0 // indirect
)

replace github.com/mbezuidenhout/tasmota.mqtt.device.manager/v2 => ../
