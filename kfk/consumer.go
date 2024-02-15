package kfk

import (
	"context"

	"github.com/bryce4651/pkg/log"
	"github.com/segmentio/kafka-go"
)

type Handler func(kafka.Message) error

type Consumer struct {
	cfg    *KfkConsumCfg
	reader *kafka.Reader
	errCh  chan *ConsumErr
}

type ConsumErr struct {
	Err error
	Msg kafka.Message
}

func NewConsumer(conf *KfkConsumCfg) (*Consumer, error) {
	var err error
	dialer := &kafka.Dialer{}
	if conf.Sasl.Enabled {
		dialer.SASLMechanism, err = conf.Sasl.Mechanism()
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     conf.Brokers,
		Dialer:      dialer,
		MinBytes:    conf.MinBytes, // 10KB
		MaxBytes:    conf.MaxBytes, // 10MB
		GroupID:     conf.GroupID,
		GroupTopics: conf.GroupTopics,
	})
	return &Consumer{
		cfg:    conf,
		reader: reader,
		errCh:  make(chan *ConsumErr, 100),
	}, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) Start(ctx context.Context, fn Handler) {
	go func() {
		for cErr := range c.errCh {
			log.Error(cErr.Err.Error())
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.errCh <- &ConsumErr{Err: err, Msg: m}
			continue
		}

		err = fn(m)
		if err != nil {
			c.errCh <- &ConsumErr{Err: err, Msg: m}
		}
	}
}
