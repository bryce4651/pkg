package kfk

import (
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type KfkConsumCfg struct {
	Brokers     []string `yaml:"brokers"`
	GroupTopics []string `yaml:"group_topics"`
	GroupID     string   `yaml:"group_id"`
	BatchSize   int      `yaml:"batch_size"`
	MinBytes    int      `yaml:"min_bytes"`
	MaxBytes    int      `yaml:"max_bytes"`
	Sasl        SaslCfg  `yaml:"sasl"`
}

type KfkProducerCfg struct {
	Brokers                []string `yaml:"brokers"`
	Topic                  string   `yaml:"topic"`
	Acks                   int      `yaml:"acks"`
	Async                  bool     `yaml:"async"`
	Timeout                int      `yaml:"timeout"`
	Sasl                   SaslCfg  `yaml:"sasl"`
	AllowAutoTopicCreation bool     `yaml:"allow_auto_topic_creation"`
	BatchSize              int      `yaml:"batch_size"`
	BatchBytes             int      `yaml:"batch_bytes"`
}

type SaslCfg struct {
	Enabled  bool   `yaml:"enabled"` //密码开关
	Algo     string `yaml:"algo"`    // 身份验证加密机制
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (sc *SaslCfg) Mechanism() (sasl.Mechanism, error) {
	switch sc.Algo {
	case "plain":
		return plain.Mechanism{
			Username: sc.Username,
			Password: sc.Password,
		}, nil
	case "scram-sha-256":
		return scram.Mechanism(scram.SHA256, sc.Username, sc.Password)
	case "scram-sha-512":
		return scram.Mechanism(scram.SHA512, sc.Username, sc.Password)
	}
	return nil, nil
}
