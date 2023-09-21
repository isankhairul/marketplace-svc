package elastic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alitto/pond"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"golang.org/x/exp/slices"
	"log"
	"marketplace-svc/pkg/util"
	"os"
	"time"
)

func (e elasticClient) BulkIndex(body []interface{}, indexName string, filename string, flush bool) error {
	if len(body) == 0 {
		return errors.New("empty body")
	}

	if flush {
		file, err := os.ReadFile("storage/data/elasticsearch/settings/" + filename + ".json")
		if err != nil {
			e.GetLogger().Error(errors.New("unable to read file: %v" + err.Error()))
			return err
		}
		// check exist index
		_, err = e.Ec.Indices.Exists([]string{indexName})
		if err == nil {
			// then delete
			_, _ = e.Ec.Indices.Delete([]string{indexName})
		}
		// create index with settings and mappings
		resp, err := esapi.IndicesCreateRequest{
			Index: indexName,
			Body:  bytes.NewReader(file),
		}.Do(context.Background(), e.GetClient())
		fmt.Println("IndicesCreateResponse: ", resp, " error: ", err)
	}

	// PARTITION
	var chunkPartition = 100
	if len(body)/chunkPartition == 0 {
		chunkPartition = 1
	}
	e.GetLogger().Info("indexName: " + indexName + ", totalBody: " + fmt.Sprint(len(body)) + ", chunkPartition: " + fmt.Sprint(chunkPartition))

	arrChunkBody := util.ChunkSlice(body, chunkPartition)
	pool := pond.New(10, len(body), pond.IdleTimeout(30*time.Second))

	for _, arrBodys := range arrChunkBody {
		arrBody := arrBodys
		pool.Submit(func() {
			var ndJSON bytes.Buffer
			enc := json.NewEncoder(bufio.NewWriter(&ndJSON))
			var arrBodyUpdated []interface{}

			// loop for add header index and _id
			for _, item := range arrBody {
				itemJSON, _ := json.Marshal(item)
				var mapItem map[string]interface{}
				_ = json.Unmarshal(itemJSON, &mapItem)
				id := fmt.Sprint(mapItem["id"])
				if id == "" {
					continue
				}
				mapHeader := map[string]interface{}{
					"index": map[string]string{"_id": id, "_index": indexName},
				}
				arrBodyUpdated = append(arrBodyUpdated, mapHeader)
				arrBodyUpdated = append(arrBodyUpdated, item)
			}

			if len(arrBodyUpdated) > 0 {
				// processing convert to ndjson
				for _, item := range arrBodyUpdated {
					_ = enc.Encode(item)
				}
				resp, err := esapi.BulkRequest{
					Index: indexName,
					Body:  bytes.NewReader(ndJSON.Bytes()),
				}.Do(context.Background(), e.GetClient())

				if err != nil {
					e.GetLogger().Error(errors.New("error BulkIndexing: " + err.Error()))
				}
				if !slices.Contains([]int{200, 201}, resp.StatusCode) {
					e.GetLogger().Info("response error bulkRequest: " + fmt.Sprint(resp))
				}
			}
		})
	}
	pool.StopAndWait()

	return nil
}

func (e elasticClient) Index(body []interface{}, indexName string) error {
	if len(body) == 0 {
		return errors.New("empty body")
	}
	e.GetLogger().Info("indexName: " + indexName + ", totalBody: " + fmt.Sprint(len(body)))

	// init NewBulkIndexer
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     e.GetClient(),
		Index:      indexName, // The default index name
		NumWorkers: 5,         // The number of worker goroutines
	})

	// processing
	for _, item := range body {
		itemJSON, _ := json.Marshal(item)
		var mapItem map[string]interface{}
		_ = json.Unmarshal(itemJSON, &mapItem)
		id := fmt.Sprint(mapItem["id"])
		if id == "" {
			return errors.New("field id is required")
		}

		err = indexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: id,
				Body:       bytes.NewReader(itemJSON),
				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			})
		if err != nil {
			return err
		}
	}

	err = indexer.Close(context.Background())
	if err != nil {
		return err
	}

	return nil
}
