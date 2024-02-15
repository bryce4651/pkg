package es

type EsCfg struct {
	Addr     string  `yaml:"addr" json:"addr"`
	User     string  `yaml:"user" json:"user"`
	Password string  `yaml:"password" json:"password"`
	Index    string  `yaml:"index" json:"index"`
	BulkCfg  BulkCfg `yaml:"bulk_cfg" json:"bulk_cfg"`
}

type BulkCfg struct {
	BulkActions int `yaml:"bulk_actions" json:"bulk_actions"`
	Workers     int `yaml:"workers" json:"workers"`
}
