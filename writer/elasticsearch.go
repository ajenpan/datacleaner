package writer

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

type ESWriter struct {
	ES        *elasticsearch.Client
	indexName string
}

func NewEsWriter(indexName string, conf *elasticsearch.Config) *ESWriter {
	if conf == nil {
		conf = &elasticsearch.Config{}
	}
	es, err := elasticsearch.NewClient(*conf)
	if err != nil {
		return nil
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	// es.Indices.Create(indexName, nil)
	// var r map[string]interface{}
	// if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
	// 	log.Fatalf("Error parsing the response body: %s", err)
	// }
	// // Print client and server version numbers.
	// log.Printf("Client: %s", elasticsearch.Version)
	// log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	// log.Println(strings.Repeat("~", 37))

	return &ESWriter{ES: es, indexName: indexName}
}

func (w *ESWriter) Work(chunk []string) error {
	//TODO:
	// var buf bytes.Buffer
	// for _, v := range chunk {
	// 	buf.Write([]byte(v))
	// }
	return nil
}
