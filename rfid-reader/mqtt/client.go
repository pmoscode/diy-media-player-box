package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const notConnected = "Mqtt Client not connected! Call 'connect' method..."

func (m Message) FromJson() string {
	marshal, err := json.Marshal(m.Value)
	if err != nil {
		return fmt.Sprintf("%v", m.Value)
	}
	return string(marshal)
}

func (m Message) ToStruct(target interface{}) {
	message := m.ToRawString()
	err := json.Unmarshal([]byte(message), &target)
	if err != nil {
		log.Println("Message is not a valid Json: ", message)
	}
}

func (m Message) ToRawString() string {
	return string(m.Value.([]uint8))
}

type Client struct {
	brokerIp string
	port     int
	topic    string
	clientId string
	client   *mqtt.Client
}

func (c *Client) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.brokerIp, c.port))
	opts.SetClientID(c.clientId)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Could not connect to broker: ", token.Error())

		return token.Error()
	}

	c.client = &client
	log.Println("Mqtt connected to", c.brokerIp)

	return nil
}

func (c *Client) Disconnect() {
	client := *c.client
	client.Disconnect(100)
}

func (c *Client) Publish(message *Message) {
	if c.client == nil {
		log.Println(notConnected)
	} else {
		client := *c.client
		token := client.Publish(message.Topic, 2, false, message.FromJson())
		token.Wait()
	}
}

func (c *Client) Subscribe(topic string, fn func(message Message)) {
	if c.client == nil {
		log.Println(notConnected)
	} else {
		client := *c.client
		client.Subscribe(topic, 2, func(client mqtt.Client, msg mqtt.Message) {
			message := Message{
				Topic: msg.Topic(),
				Value: msg.Payload(),
			}
			fn(message)
		})
	}
}

func (c *Client) LoopForever() {
	if c.client == nil {
		log.Println(notConnected)
	} else {
		channel := make(chan os.Signal)
		signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-channel
			c.Disconnect()
			log.Println("Exiting...")
			os.Exit(1)
		}()

		for {
			time.Sleep(1 * time.Second)
		}
	}
}

func CreateClient(brokerIp string, port int, clientId string) *Client {
	return &Client{
		brokerIp: brokerIp,
		port:     port,
		clientId: clientId,
	}
}
