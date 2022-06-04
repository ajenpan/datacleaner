package writer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type ESWriter struct {
	ES *elasticsearch.Client
	BI esutil.BulkIndexer

	IndexName       string
	CountSuccessful uint64
}

func NewEsWriter(indexName string, conf *elasticsearch.Config) (*ESWriter, error) {
	if conf == nil {
		conf = &elasticsearch.Config{}
	}
	es, err := elasticsearch.NewClient(*conf)
	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}

	if res.IsError() {
		return nil, fmt.Errorf("error: %s", res.String())
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,        // The default index name
		Client:        es,               // The Elasticsearch client
		FlushBytes:    5e+6,             // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})

	if err != nil {
		return nil, err
	}

	return &ESWriter{ES: es, IndexName: indexName, BI: bi}, nil
}

func (w *ESWriter) Close() error {
	return w.BI.Close(context.Background())
}

func (w *ESWriter) Write(line map[string]interface{}) error {

	b, err := json.Marshal(line)
	if err != nil {
		return err
	}

	err = w.BI.Add(
		context.Background(),
		esutil.BulkIndexerItem{
			// Action field configures the operation to perform (index, create, delete, update)
			Action: "create",
			// Body is an `io.Reader` with the payload
			Body: bytes.NewReader(b),
			// OnSuccess is called for each successful operation
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				atomic.AddUint64(&w.CountSuccessful, 1)
			},
			// OnFailure is called for each failed operation
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
				}
			},
		})
	return err
}
