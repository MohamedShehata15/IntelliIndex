package elasticsearch

import (
	"bytes"
	"context"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func (c *Client) IndexDocument(ctx context.Context, indexName, docID string, doc interface{}) error {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	docJSON := mustMarshalJSON(doc)
	res, err := c.PerformRequest(ctx, &esapi.IndexRequest{
		Index:      fullIndexName,
		DocumentID: docID,
		Body:       bytes.NewReader(docJSON),
	})
	if err != nil {
		return err
	}
	defer CloseBody(res.Body)
	return nil
}
