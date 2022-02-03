package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mbezuidenhout/tasmota.mqtt.device.manager/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Host        string `yaml:"host"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Customtopic string `yaml:"customtopic"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	config, err := NewConfig("config.yml")
	if err != nil {
		return
	}
	mqttClientOptions := mqtt.NewClientOptions()
	mqttClientOptions.SetUsername(config.Username).SetPassword(config.Password).AddBroker(config.Host)
	mqttClientOptions.SetClientID("TMDM_DEV")

	m := tasmota.NewManager(*mqttClientOptions)
	m.AddTopic(config.Customtopic)
	m.Connect()

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
