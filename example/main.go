package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mbezuidenhout/tasmota.mqtt.device.manager/v2"
)

func main() {
	mqttClientOptions := mqtt.NewClientOptions()
	mqttClientOptions.SetUsername("root").SetPassword("$4!KPx^*K@5*2p").AddBroker("tcp://mqtt.lan:1883")
	mqttClientOptions.SetClientID("TMDM_DEV")

	m := tasmota.NewManager(*mqttClientOptions)

	if m.MQTTclient.IsConnected() {
		fmt.Println("MQTT is connected")
		defer m.Disconnect()
	}

	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				devices := m.GetDevices()
				fmt.Printf("There are now %d devices online", len(devices))
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	devices := m.GetDevices()

	fmt.Println(devices)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {

		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
