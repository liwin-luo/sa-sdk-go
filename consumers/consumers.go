package consumers

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"

	"sd_sdk/structs"
	"sd_sdk/utils"
)

type Consumer interface {
	Send(data structs.EventData) error
	Flush() error
	Close() error
	ItemSend(item structs.Item) error
}

func send(url string, data string, to time.Duration, list bool) error {
	pdata := ""

	if list {
		rdata, err := utils.GeneratePostDataList(data)
		if err != nil {
			return err
		}
		pdata = rdata
	} else {
		rdata, err := utils.GeneratePostData(data)
		if err != nil {
			return err
		}
		pdata = rdata
	}
	err := utils.DoRequest(url, pdata, to)
	if err != nil {
		return err
	}

	return nil
}

const (
	groupID = "sd_sdk"
)

func brokerReader(broker []string, topic string, message func(data string) error) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  broker,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	go func(r *kafka.Reader) {
		ctx := context.Background()
		for {
			m, err := r.FetchMessage(ctx)
			if err != nil {
				break
			}
			_ = message(string(m.Value))
			_ = r.CommitMessages(ctx, m)
		}
	}(reader)

	return reader
}

func brokerWriter(broker []string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  broker,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}
