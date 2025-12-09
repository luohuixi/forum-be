package model

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/spf13/viper"
)

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
	KafkaReaderInit(topic)
	KafkaWriterInit(topic)
}

// ------------------- Reader -------------------

func KafkaReaderInit(topic string) {
	if KafkaReader == nil {
		// SASL PLAIN 配置
		mechanism := plain.Mechanism{
			Username: viper.GetString("kafka.username"),
			Password: viper.GetString("kafka.password"),
		}

		dialer := &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           nil,
		}

		KafkaReader = &kafkaReader{
			Self: kafka.NewReader(kafka.ReaderConfig{
				Brokers:   []string{viper.GetString("kafka.addr")},
				Topic:     topic,
				Partition: 0,
				Dialer:    dialer,
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

// ------------------- Writer -------------------

func KafkaWriterInit(topic string) {
	if KafkaWriter == nil {
		mechanism := plain.Mechanism{
			Username: viper.GetString("kafka.username"),
			Password: viper.GetString("kafka.password"),
		}

		dialer := &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           nil,
		}

		KafkaWriter = &kafkaWriter{
			Self: kafka.NewWriter(kafka.WriterConfig{
				Brokers:  []string{viper.GetString("kafka.addr")},
				Topic:    topic,
				Balancer: &kafka.LeastBytes{},
				Dialer:   dialer,
			}),
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
