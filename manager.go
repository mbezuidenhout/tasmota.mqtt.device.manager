package tasmota

import (
	"fmt"
	"regexp"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var topicPrefixes = map[string]struct{}{
	"tele": {},
	"cmnd": {},
	"stat": {},
}

type Manager struct {
	MQTTClientOptions mqtt.ClientOptions
	MQTTclient        mqtt.Client
	devices           map[string]*Device
	topics            []string
}

// NewManager will create a new Manager with mqtt options from mqttClientOptions parameter
func NewManager(mqttClientOptions mqtt.ClientOptions) *Manager {
	o := &Manager{
		MQTTClientOptions: mqttClientOptions,
		topics:            []string{"%prefix%/%topic%/"},
		devices:           make(map[string]*Device),
	}
	o.MQTTclient = mqtt.NewClient(&o.MQTTClientOptions)
	o.Connect()
	return o
}

func (m *Manager) IsConnected() bool {
	return m.MQTTclient.IsConnected()
}

func (m *Manager) Disconnect() {
	m.MQTTclient.Disconnect(3)
	fmt.Println("MQTT disconnected")
}

func (m *Manager) Connect() {
	if m.MQTTclient.IsConnected() {
		return
	}
	if token := m.MQTTclient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	subscribeTopics := make(map[string]byte)
	for _, t := range m.topics {
		subscribeTopics[strings.Replace(strings.Replace(t, "%prefix%", "+", 1), "%topic%", "+", 1)+"LWT"] = 0
	}
	token := m.MQTTclient.SubscribeMultiple(subscribeTopics, m.MessageHandler)
	token.Wait()
}

func (m *Manager) ReConnect() {
	if m.MQTTclient.IsConnected() {
		m.MQTTclient.Disconnect(3)
	}
	m.Connect()
}

func (m *Manager) MessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("%s: %s\n", msg.Topic(), msg.Payload())
	if strings.HasSuffix(msg.Topic(), "LWT") {
		if fullTopic, ok := findFullTopic(m.topics, msg.Topic()); ok {
			topic := getTopic(fullTopic, msg.Topic())
			if _, ok = m.devices[topic]; !ok {
				//m.devices[topic] = &Device{topic: topic, fullTopic: fullTopic}
				m.devices[topic] = NewDevice(topic, fullTopic, m.MQTTclient)
			}
		}
	}

}

func (m *Manager) GetDevices() map[string]*Device {
	return m.devices
}

func getTopic(fullTopic string, topic string) string {
	regexString := strings.Replace(strings.Replace(fullTopic, "%prefix%", "(?P<prefix>.*?)", 1), "%topic%", "(?P<topic>.*?)", 1) + ".*$"
	myExp := regexp.MustCompile(regexString)
	match := myExp.FindStringSubmatch(topic)
	result := make(map[string]string)
	for i, name := range myExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result["topic"]
}

func findFullTopic(topics []string, topic string) (string, bool) {
	for _, v := range topics {
		regexString := strings.Replace(strings.Replace(v, "%prefix%", "(?P<prefix>.*?)", 1), "%topic%", "(?P<topic>.*?)", 1) + ".*$"
		myExp := regexp.MustCompile(regexString)
		match := myExp.FindStringSubmatch(topic)
		if match != nil {
			result := make(map[string]string)
			for i, name := range myExp.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			_, ok := topicPrefixes[result["prefix"]]
			if ok {
				return v, true
			} else {
				continue
			}
		}
	}
	return "", false
}
