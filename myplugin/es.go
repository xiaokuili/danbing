package myplugin

import (
	"bytes"
	"context"
	"danbing/conf"
	"danbing/plugin"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type EsWriter struct {
	cli     *elasticsearch.Client
	Connect *conf.Connect
	Query   *conf.Query
}

func (writer *EsWriter) Init(tq *conf.Query, tc *conf.Connect) {
	url := fmt.Sprintf("http://%s:%d", tc.Host, tc.Port)
	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		Username: tc.Username,
		Password: tc.Password,
	}
	cli, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
	}
	res, err := cli.Info()
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode != 200 {
		fmt.Println(res.String())

	}
	writer.cli = cli
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

func (writer *EsWriter) Writer(result []map[string]interface{}) {

	for i := 0; i < len(result); i++ {
		d := result[i]
		var docID string
		for i := 0; i < len(writer.Query.Columns); i++ {
			c := writer.Query.Columns[i]
			if c.PrimaryField {
				docID = docID + toStr(d[c.Name])
			}

		}
		fmt.Println(docID)

		data, err := json.Marshal(d)
		if err != nil {
			fmt.Println(err)
		}

		req := esapi.IndexRequest{
			Index:      writer.Query.Table,
			DocumentID: docID,
			Body:       bytes.NewReader(data),
			Refresh:    "true",
		}
		res, err := req.Do(context.Background(), writer.cli)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
		if res.IsError() {
			log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
		} else {
			// Deserialize the response into a map.
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and indexed document version.
				log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			}
		}
	}

}

// TODO: init必须手动维护
func init() {
	plugin.Register(&EsWriter{})
}
