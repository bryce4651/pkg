package kfk

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/bryce4651/pkg/log"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w     *kafka.Writer
	ch    chan *kafka.Message
	errCh chan error
	done  bool
}

func NewProducer(cfg *KfkProducerCfg) (*Producer, error) {
	mechanism, err := cfg.Sasl.Mechanism()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		Dialer: &kafka.Dialer{
			// ClientID:        "",
			Timeout:         time.Duration(cfg.Timeout * int(time.Second)),
			Deadline:        time.Time{},
			LocalAddr:       nil,
			DualStack:       true,
			FallbackDelay:   0,
			KeepAlive:       0,
			Resolver:        nil,
			TLS:             &tls.Config{},
			SASLMechanism:   mechanism,
			TransactionalID: "",
		},
		Balancer:          nil,
		MaxAttempts:       0,
		QueueCapacity:     0,
		BatchSize:         0,
		BatchBytes:        0,
		BatchTimeout:      0,
		ReadTimeout:       0,
		WriteTimeout:      0,
		RebalanceInterval: 0,
		IdleConnTimeout:   0,
		RequiredAcks:      cfg.Acks,
		Async:             cfg.Async,
		CompressionCodec:  nil,
		Logger:            nil,
		ErrorLogger:       nil,
	})
	w.AllowAutoTopicCreation = cfg.AllowAutoTopicCreation

	return &Producer{
		w:  w,
		ch: make(chan *kafka.Message, 1000),
	}, nil
}

func (p *Producer) AsynStart() {
	tk := time.NewTicker(time.Millisecond * 450)
	var msgList []kafka.Message
	for {
		select {
		case msg, ok := <-p.ch:
			if !ok && msg == nil {
				log.Warn("===> close producer")
				p.done = true
				p.Send(msgList...)
				return
			}
			if msg == nil {
				continue
			}
			msgList = append(msgList, *msg)
			if len(msgList) >= 1000 {
				p.Send(msgList...)
				msgList = msgList[:0]
			}
		case <-tk.C:
			if len(msgList) == 0 {
				continue
			}
			p.Send(msgList...)
			msgList = msgList[:0]
		}

	}
}

func (p *Producer) Send(msgList ...kafka.Message) {
	log.Infof("send producer: %d", len(msgList))
	err := p.w.WriteMessages(context.Background(), msgList...)
	if err != nil {
		log.Error(err)
		return
	}
}

func (p *Producer) Pub(key, value []byte) {
	p.ch <- &kafka.Message{
		Key:   key,
		Value: value,
	}
}

func (p *Producer) Close() {
	close(p.ch)
	for range time.NewTicker(time.Microsecond * 100).C {
		if p.done {
			log.Warn("====> close producer <====")
			break
		}
	}
	if err := p.w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
