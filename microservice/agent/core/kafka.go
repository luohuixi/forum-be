package core

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

const (
	NumPartitions     = 1 //分区数量
	ReplicationFactor = 1 //分区的副本数量
)

var KafkaWriter *kafka.Writer
var KafkaTopic = "forum_agent"

func KafkaInit() {
	if v := viper.GetString("kafka.topic"); v != "" {
		KafkaTopic = v
	}

	addr := viper.GetString("kafka.addr")
	if err := ensureTopic(addr, KafkaTopic); err != nil {
		log.Fatalf("ensure kafka topic failed: %v", err)
	}

	//mechanism := plain.Mechanism{
	//	Username: viper.GetString("kafka.username"),
	//	Password: viper.GetString("kafka.password"),
	//}
	//
	//dialer := &kafka.Dialer{SASLMechanism: mechanism, TLS: nil}
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Topic:    KafkaTopic,
		Balancer: &kafka.LeastBytes{},
		//Transport: &kafka.Transport{
		//	Dial: dialer.DialFunc,
		//	SASL: dialer.SASLMechanism,
		//	TLS:  dialer.TLS,
		//},
	}

	if KafkaWriter == nil {
		log.Fatal("kafka writer init failed")
	}
}

func KafkaPublish(msg []byte) error {
	if KafkaWriter == nil {
		KafkaInit()
	}
	return KafkaWriter.WriteMessages(context.Background(), kafka.Message{Value: msg})
}

func KafkaReader() *kafka.Reader {
	if KafkaWriter == nil {
		KafkaInit()
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{viper.GetString("kafka.addr")},
		Topic:   KafkaTopic,
		GroupID: viper.GetString("kafka.group_id"),
	})

	return reader
}

func ensureTopic(addr, topic string) error {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err == nil && len(partitions) > 0 {
		return nil
	}

	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     NumPartitions,
		ReplicationFactor: ReplicationFactor,
	})
}
