package myplugin

import (
	"bytes"
	"context"
	"danbing/plugin"
	"danbing/task"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type EsWriter struct {
	cli     *elasticsearch.Client
	Connect *task.Connect
}

func (writer *EsWriter) Init(tq *task.Query, tc *task.Connect) {
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

func (writer *EsWriter) Writer(s string) {
	result := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(result); i++ {
		d := result[i]
		data, err := json.Marshal(d)
		if err != nil {
			fmt.Println(err)
		}
		req := esapi.IndexRequest{
			Index:      "danbing",
			DocumentID: strconv.Itoa(i + 1),
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
