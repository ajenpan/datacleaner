package writer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

type ESWriter struct{}

func NewEsWriter() *ESWriter {
	return &ESWriter{}
}

func (w *ESWriter) Run() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	fmt.Println(es)
}

func (w *ESWriter) Write(chunk []string) error {

	addresses := []string{"http://127.0.0.1:9200", "http://127.0.0.1:9201"}
	config := elasticsearch.Config{
		Addresses: addresses,
		Username:  "",
		Password:  "",
		CloudID:   "",
		APIKey:    "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	fmt.Println(err, "Error creating the client")
	// Index creates or updates a document in an index
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"title":   "你看到外面的世界是什么样的？",
		"content": "外面的世界真的很精彩",
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		fmt.Println(err, "Error encoding doc")
	}

	res, err := es.Index("demo", &buf, es.Index.WithDocumentID("doc"))

	if err != nil {
		fmt.Println(err, "Error Index response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
	return nil
}
