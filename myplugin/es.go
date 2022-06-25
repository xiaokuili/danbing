package myplugin

import (
	"bytes"
	"context"
	"danbing/conf"
	"danbing/plugin"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

type EsWriter struct {
	es      *elasticsearch.Client
	Connect *conf.Connect
	Query   *conf.Query
}

func (writer *EsWriter) Init(tq *conf.Query, tc *conf.Connect) {
	url := fmt.Sprintf("http://%s:%d", tc.Host, tc.Port)
	retryBackoff := backoff.NewExponentialBackOff()

	es, err := elasticsearch.NewClient(elasticsearch.Config{

		// Retry on 429 TooManyRequests statuses
		//
		RetryOnStatus: []int{502, 503, 504},

		// Configure the backoff function
		//
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},

		// Retry up to 5 attempts
		//
		MaxRetries: 5,
		Addresses: []string{
			url,
		},
		Username: tc.Username,
		Password: tc.Password,
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	writer.es = es
	writer.Query = tq
}

func (writer *EsWriter) Name() string {
	return "eswriter"
}

func (writer *EsWriter) Split(taskNum int) []plugin.WriterPlugin {
	plugins := make([]plugin.WriterPlugin, 0)
	for i := 0; i < taskNum; i++ {
		plugins = append(plugins, writer)
	}
	return plugins
}

func (writer *EsWriter) Close() {

}

func toStr(i interface{}) string {
	switch t := i.(type) {
	case string:
		return t
	}
	return ""
}

func CreateID(column []string, data map[string]interface{}) string {
	docID := ""
	for i := 0; i < len(column); i++ {
		f := column[i]
		docID = docID + toStr(data[f])
	}
	if docID == "" {
		panic("id不能为空")
	}
	return docID
}

func (writer *EsWriter) Writer(result []map[string]interface{}) {

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         writer.Query.Table,
		Client:        writer.es,
		NumWorkers:    3,
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}
	start := time.Now().UTC()

	for i := 0; i < len(result); i++ {
		docID := CreateID(writer.Query.Primary, result[i])
		d, err := json.Marshal(result[i])
		if err != nil {
			panic(err)
		}
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "index",

				DocumentID: docID,

				Body: bytes.NewReader(d),
			},
		)
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	biStats := bi.Stats()

	dur := time.Since(start)

	if biStats.NumFailed > 0 {

		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)

	}

}

// TODO: init必须手动维护
func init() {
	plugin.Register(&EsWriter{})
}
