package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	sdk "sd_sdk"
)

func TestDefaultConsumer(t *testing.T) {
	go MockServerRun()

	c, err := sdk.InitDefaultConsumer("http://127.0.0.1:8080/p?project=xingletao", 1000 /*"192.168.2.193:9092"*/)
	if err != nil {
		fmt.Println(err)
		return
	}

	sa := sdk.InitSensorsAnalytics(c, "xingletao", false)
	defer sa.Close()

	distinctId := "829c4c8959f54d81a54398a4757100705100021"
	event := "order"
	properties := make(map[string]interface{})
	properties["$order_id"] = "c4c8959"
	properties["$item_id"] = "710070510002"
	properties["$platform_rebate"] = 20
	properties["$customer_commission"] = 2
	properties["$order_source"] = 1
	properties["$payment_amount"] = 10
	properties["$order_time"] = int(time.Now().Unix())

	err = sa.Track(distinctId, event, properties, true)
	if err != nil {
		t.Fatal("default consumer track failed", err)
		return
	}
	sa.Close()

	t.Log("Default consumer ok")
}

func Benchmark_TestDefaultConsumer(b *testing.B) {
	c, err := sdk.InitDefaultConsumer("http://127.0.0.1:8080/p?project=xingletao", 1000, "192.168.2.193:9092")
	if err != nil {
		fmt.Println(err)
		return
	}

	sa := sdk.InitSensorsAnalytics(c, "xingletao", false)
	defer sa.Close()

	for i := 0; i < b.N; i++ {
		r := strconv.Itoa(i)
		distinctId := "829c4c8959f54d81a54398a47571007051000" + r
		event := "order"
		properties := make(map[string]interface{})
		properties["$order_id"] = "c4c8959" + r
		properties["$item_id"] = "710070510002" + r
		properties["$platform_rebate"] = 20
		properties["$customer_commission"] = 2
		properties["$order_source"] = 1
		properties["$payment_amount"] = 10
		properties["$order_time"] = time.Now().Format("2006-01-02 15:04:05")

		err = sa.Track(distinctId, event, properties, true)
		if err != nil {
			b.Fatal("default consumer track failed", err)
			return
		}
	}

}
