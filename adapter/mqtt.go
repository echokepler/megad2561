package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/echokepler/megad2561/core"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type MQTTClient struct {
	megadID    *string
	commander  core.ExecuteCommander
	connection mqtt.Client
}

type MQTTClientOptions struct {
	Address  string
	ClientID *string
	Password string
}

type MegadPortInMessage struct {
	Port  int    `json:"port"`
	Mode  int8   `json:"m"`
	Value string `json:"value"`
	Click int8   `json:"click,omitempty"`
	Count int32  `json:"cnt"`
}

func NewMqttClient(opts MQTTClientOptions) (*MQTTClient, error) {
	address := fmt.Sprintf("tcp://%v", opts.Address)

	mqttOpts := mqtt.NewClientOptions().AddBroker(address).SetClientID("MegadGO")
	mqttOpts.SetPingTimeout(1 * time.Second)
	mqttOpts.SetPassword(opts.Password)
	mqttOpts.SetUsername(*opts.ClientID)

	client := mqtt.NewClient(mqttOpts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &MQTTClient{
		connection: client,
		megadID:    opts.ClientID,
	}, nil
}

func (mc *MQTTClient) DoCommand(commander core.CommandReader) error {
	return mc.publish("cmd", commander.Convert())
}

func (mc *MQTTClient) SubscribePortIn(ID uint8, handler core.MessageHandlerCallback) error {
	topic := fmt.Sprintf("%v/%v", *mc.megadID, ID)

	if token := mc.connection.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var message core.MegadPortInMessage

		err := json.Unmarshal(msg.Payload(), &message)
		if err != nil {
			fmt.Println(err)
		}

		if message.Port == int(ID) {
			handler(message)
		}
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mc *MQTTClient) publish(topic string, msg string) error {
	token := mc.connection.Publish(fmt.Sprintf("%v/%v", *mc.megadID, topic), 0, false, msg)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	return nil
}
