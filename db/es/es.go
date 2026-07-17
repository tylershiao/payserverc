package es

import (
	"context"

	"github.com/olivere/elastic/v7"
)

// EsClient es连接对象
type EsClient struct {
	client  *elastic.Client
	running bool
	enable  bool
}

func InitEsClient(urls []string, sniff bool) (*EsClient, error) {
	ES := &EsClient{}

	client, err := elastic.NewClient(elastic.SetURL(urls...), elastic.SetSniff(sniff))
	if err != nil {
		return nil, err
	}

	_, _, err = client.Ping(urls[0]).Timeout("1s").Do(context.TODO())
	if err != nil {
		return nil, err
	}

	ES.client = client
	ES.running = true
	ES.enable = true

	// InitBulkQueue(client)

	return ES, nil
}

func (ES *EsClient) UnInit() {
	ES.running = false
}
