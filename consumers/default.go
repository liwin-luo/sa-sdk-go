package consumers

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"

	"sd_sdk/structs"
)

type DefaultConsumer struct {
	Url     string
	Timeout time.Duration

	Asynchronous bool
	brokerReader *kafka.Reader
	brokerWriter *kafka.Writer
}

func InitDefaultConsumer(url string, timeout int, broker ...string) (*DefaultConsumer, error) {
	c := &DefaultConsumer{Url: url, Timeout: time.Duration(timeout) * time.Millisecond}
	if len(broker) != 0 {
		c.Asynchronous = true
		c.brokerReader = brokerReader(broker, "sd_default", c.async)
		c.brokerWriter = brokerWriter(broker, "sd_default")
	}
	return c, nil
}

func (c *DefaultConsumer) async(data string) error {
	return send(c.Url, data, c.Timeout, false)
}

func (c *DefaultConsumer) sendAsync(data []byte) error {
	return c.brokerWriter.WriteMessages(context.TODO(), kafka.Message{
		Topic:     "sd_default",
		Partition: 0,
		Value:     data,
	})
}

func (c *DefaultConsumer) Send(data structs.EventData) error {
	jdata, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if c.Asynchronous {
		return c.sendAsync(jdata)
	}

	return send(c.Url, string(jdata), c.Timeout, false)
}

func (c *DefaultConsumer) ItemSend(item structs.Item) error {
	itemData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return send(c.Url, string(itemData), c.Timeout, false)
}

func (c *DefaultConsumer) Flush() error {
	return nil
}

func (c *DefaultConsumer) Close() error {
	if c.Asynchronous {
		_ = c.brokerReader.Close()
		_ = c.brokerWriter.Close()
	}
	return nil
}
