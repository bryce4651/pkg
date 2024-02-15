package es

type EsCfg struct {
	Addr     string  `yaml:"addr"`
	User     string  `yaml:"user"`
	Password string  `yaml:"password"`
	Index    string  `yaml:"index"`
	BulkCfg  BulkCfg `yaml:"bulk_cfg"`
}

type BulkCfg struct {
	BulkActionsNum int `yaml:"bulk_actions_num"`
	WorkNum        int `yaml:"work_num"`
}
