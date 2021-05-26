package elasticsearch

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"encoding/json"
	"github.com/bnkamalesh/errors"
	errors2 "github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v3"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"
)

const eventsPerFetch = 100

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Name: "ElasticSearch",
		Type: reflect.TypeOf(elasticSearch{}),
		CanHandle: func(url model.SourceURL) bool {
			_, ok := map[string]bool{
				scheme + "+http":            true,
				schemeAliasShort + "+http":  true,
				scheme + "+https":           true,
				schemeAliasShort + "+https": true,
			}[url.Scheme]
			return ok
		},
		Creator: func(sourceURL model.SourceURL) (*interface{}, error) {
			src, err := newSource(sourceURL)
			if err != nil {
				return nil, err
			}
			var i interface{} = src
			return &i, nil
		},
		ExampleURLs: []string{
			"elasticsearch+http://...",
			"elasticsearch+https://...",
			"es+http://...",
			"es+https://...",
		},
	}
}

const keepAliveOneMinute = "1m"

const (
	scheme           = "elasticsearch"
	schemeAliasShort = "es"
)

type eventHolder struct {
	Event string `json:"log"`
}

// elasticSearch provides elasticsearch test data.
type elasticSearch struct {
	source.BaseSource
	query     elastic.Query
	scrollID  string
	indexName string
	esClient  *elastic.Client
}

// newSource ...
func newSource(sourceURL model.SourceURL) (*source.Source, error) {
	normalizeURL(&sourceURL)

	esClient, err := newElasticsearchClient(sourceURL)
	if err != nil {
		return nil, err
	}

	indexName := strings.Trim(sourceURL.Path, "/")
	query := elastic.NewQueryStringQuery(sourceURL.RawQuery)

	var src source.Source = elasticSearch{
		BaseSource: source.New('፨', sourceURL),
		esClient:   esClient,
		query:      query,
		indexName:  indexName,
	}
	return &src, nil
}

func newElasticsearchClient(sourceURL model.SourceURL) (*elastic.Client, error) {
	esClient, err := elastic.NewClient(elastic.SetURL(sourceURL.Scheme + "://" + sourceURL.Host))
	if err != nil {
		return nil, err
	}
	return esClient, nil
}

func normalizeURL(sourceURL *model.SourceURL) {
	sourceURL.Scheme = strings.Replace(sourceURL.Scheme, scheme+"+", "", 1)
	sourceURL.Scheme = strings.Replace(sourceURL.Scheme, schemeAliasShort+"+", "", 1)
}

// Start ...
func (src elasticSearch) Start(wg *sync.WaitGroup, running func() bool, out chan source.Event) error {
	go func() {
		defer wg.Done()
		for running() {
			result, err := src.esClient.Scroll(src.indexName).
				KeepAlive(keepAliveOneMinute).
				ScrollId(src.scrollID).
				Size(eventsPerFetch).
				Query(src.query).
				SortBy(elastic.NewFieldSort(model.TimestampField)).
				Pretty(true).
				Do()
			if err != nil {
				if errors.Is(err, io.EOF) {
					continue
				}
				panic(errors2.Wrapf(err, "name=%s\n", src.indexName))
			}
			if result == nil || result.Hits == nil || result.Hits.Hits == nil {
				continue
			}
			for _, hit := range result.Hits.Hits {
				holder := eventHolder{}
				err = json.Unmarshal(*hit.Source, &holder)
				if err != nil {
					panic(err)
				}
				var evt = model.Event{}
				err = json.Unmarshal([]byte(holder.Event), &evt)
				if err != nil {
					evt = model.Event{
						Timestamp:       model.Timestamp(time.Now()),
						Level:           level.Info,
						Message:         model.Message(holder.Event),
						ApplicationName: nil,
						MethodName:      "-",
						LineNumber:      model.NoLineNumber,
						ThreadName:      nil,
						ClassName:       "-",
						Exception:       nil,
					}
				}
				out <- source.Event{Source: src, Event: evt}
			}
			src.scrollID = result.ScrollId
		}
	}()
	return nil
}
