module example

go 1.17

require (
	github.com/eclipse/paho.mqtt.golang v1.4.1
	github.com/mbezuidenhout/tasmota.mqtt.device.manager/v2 v2.0.0-20220617123114-cab69f9c898e
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
)

replace github.com/mbezuidenhout/tasmota.mqtt.device.manager/v2 => ../
