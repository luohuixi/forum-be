package model

import (
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)
import "context"

// KafkaWriter ...
type kafkaWriter struct {
	Self *kafka.Writer
}

// KafkaReader ...
type kafkaReader struct {
	Self *kafka.Reader
}

var KafkaReader *kafkaReader
var KafkaWriter *kafkaWriter

func InitKafka(topic string) {
	KafkaReader.Init(topic)
	KafkaWriter.Init(topic)

}

func (*kafkaReader) Init(topic string) {
	if KafkaReader == nil {
		KafkaReader = &kafkaReader{
			Self: kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{viper.GetString("kafka.addr")},
				Topic:     topic,
				Partition: 0,
			}),
		}
	}
}

func (r kafkaReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return r.Self.FetchMessage(ctx)
}

func (r kafkaReader) CommitMessage(ctx context.Context, msg kafka.Message) error {
	return r.Self.CommitMessages(ctx, msg)
}

func (*kafkaWriter) Init(topic string) {
	if KafkaReader != nil {
		KafkaWriter = &kafkaWriter{
			Self: &kafka.Writer{
				Addr:     kafka.TCP(viper.GetString("kafka.addr")),
				Topic:    topic,
				Balancer: &kafka.LeastBytes{},
			},
		}
	}
}

func (w kafkaWriter) PublishMessage(msg []byte) error {
	return w.Self.WriteMessages(context.Background(),
		kafka.Message{
			Value: msg,
		},
	)
}
