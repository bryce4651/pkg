package es

import (
	"context"

	"github.com/olivere/elastic/v7"
)

type ESWriter struct {
	index     string
	cli       *elastic.Client
	processor *elastic.BulkProcessor
}

func NewEsWriter(name string, cfg *EsCfg) (*ESWriter, error) {
	cli, err := elastic.NewClient(elastic.SetURL(cfg.Addr))
	if err != nil {
		return nil, err
	}
	// 创建一个BulkProcessor
	bulkProcessor, err := cli.BulkProcessor().
		Name(name).
		Workers(cfg.BulkCfg.Workers).
		BulkActions(cfg.BulkCfg.BulkActions).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &ESWriter{cli: cli, processor: bulkProcessor, index: cfg.Index}, nil
}

func (w *ESWriter) SetIndex(index string) {
	w.index = index
}

func (w *ESWriter) Add(doc interface{}, indexs ...string) {
	if len(indexs) == 0 {
		w.processor.Add(elastic.NewBulkIndexRequest().Index(w.index).Doc(doc))
		return
	}
	for _, index := range indexs {
		w.processor.Add(elastic.NewBulkCreateRequest().Index(index).Doc(doc))
	}
}

func (w *ESWriter) Close() error {
	if w.cli != nil {
		w.cli.Stop()
	}
	if w.processor == nil {
		return nil
	}
	return w.processor.Close()
}
