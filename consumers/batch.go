package consumers

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"

	"sd_sdk/structs"
)

const (
	BATCH_DEFAULT_MAX = 50
)

type BatchConsumer struct {
	Url        string
	Max        int
	DataBuffer []structs.EventData
	ItemBuffer []structs.Item
	Timeout    time.Duration

	Asynchronous bool
	brokerReader *kafka.Reader
	brokerWriter *kafka.Writer
}

func InitBatchConsumer(url string, max, timeout int, broker ...string) (*BatchConsumer, error) {
	if max > BATCH_DEFAULT_MAX {
		max = BATCH_DEFAULT_MAX
	}

	c := &BatchConsumer{Url: url, Max: max, Timeout: time.Duration(timeout) * time.Millisecond}
	c.DataBuffer = make([]structs.EventData, 0, max)
	if len(broker) != 0 {
		c.Asynchronous = true
		c.brokerReader = brokerReader(broker, "sd_batch", c.async)
		c.brokerWriter = brokerWriter(broker, "sd_batch")
	}
	return c, nil
}

func (c *BatchConsumer) async(data string) error {
	return send(c.Url, data, c.Timeout, true)
}

func (c *BatchConsumer) Send(data structs.EventData) error {
	c.DataBuffer = append(c.DataBuffer, data)
	if (len(c.DataBuffer) + len(c.ItemBuffer)) < c.Max {
		return nil
	}

	return c.Flush()
}

func (c *BatchConsumer) sendAsync(data []byte) error {
	return c.brokerWriter.WriteMessages(context.TODO(), kafka.Message{
		Topic:     "sd_batch",
		Partition: 0,
		Value:     data,
	})
}

func (c *BatchConsumer) Flush() error {
	// 刷新 Event 数据
	if len(c.DataBuffer) != 0 {
		jdata, err := json.Marshal(c.DataBuffer)
		if err != nil {
			return err
		}

		err = send(c.Url, string(jdata), c.Timeout, true)

		c.DataBuffer = c.DataBuffer[:0]
		return err
	}

	// 刷新 Item 数据
	if len(c.ItemBuffer) != 0 {
		itemData, err := json.Marshal(c.ItemBuffer)
		if err != nil {
			return err
		}

		if c.Asynchronous {
			err = c.sendAsync(itemData)
		} else {
			err = send(c.Url, string(itemData), c.Timeout, true)
		}

		c.ItemBuffer = c.ItemBuffer[:0]
		return err
	}
	return nil
}

func (c *BatchConsumer) Close() error {
	return c.Flush()
}

func (c *BatchConsumer) ItemSend(item structs.Item) error {
	c.ItemBuffer = append(c.ItemBuffer, item)
	if (len(c.DataBuffer) + len(c.ItemBuffer)) < c.Max {
		return nil
	}
	if c.Asynchronous {
		_ = c.brokerReader.Close()
		_ = c.brokerWriter.Close()
	}
	return c.Flush()
}
