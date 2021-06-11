package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/bnkamalesh/errors"
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/source"
	"go.uber.org/ratelimit"
	"gopkg.in/olivere/elastic.v3"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"
)

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
		Creator: newSource,
		ExampleURLs: []string{
			"elasticsearch+http://...",
			"elasticsearch+https://...",
			"es+http://...",
			"es+https://...",
		},
	}
}

const (
	scheme           = "elasticsearch"
	schemeAliasShort = "es"
)

type elasticSearch struct {
	source.BaseSource
	scroll            *elastic.ScrollService
	searchRateLimiter ratelimit.Limiter
}

func newSource(since *time.Time, sourceURL model.SourceURL) (*source.Source, error) {
	const maxRequestsPerMinute = 60
	const eventsPerFetch = 100
	const keepAliveOneMinute = "1m"

	normalizeURL(&sourceURL)

	indexName := strings.Trim(sourceURL.Path, "/")

	esClient, err := newElasticsearchClient(sourceURL)
	if err != nil {
		return nil, err
	}

	scroll := esClient.Scroll(indexName).
		KeepAlive(keepAliveOneMinute).
		Size(eventsPerFetch)

	var query elastic.Query
	if q, ok := sourceURL.QueryKey("q"); ok {
		query = elastic.NewQueryStringQuery(q)
	}

	if fieldName, ok := sourceURL.QueryKey("timestamp"); ok {
		scroll.SortBy(elastic.NewFieldSort(fieldName))
		if since != nil {
			dateQuery := elastic.NewRangeQuery(fieldName).From(since.Unix())
			if query != nil {
				query = elastic.NewBoolQuery().Must(query, dateQuery)
			} else {
				query = dateQuery
			}
		}
	}

	if query != nil {
		scroll.Query(query)
	}

	var src source.Source = elasticSearch{
		BaseSource: source.New('፨', sourceURL),
		scroll:     scroll,
		searchRateLimiter: ratelimit.New(
			maxRequestsPerMinute,
			ratelimit.Per(time.Minute),
			ratelimit.WithSlack(maxRequestsPerMinute),
		),
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
func (src elasticSearch) Start(
	wg *sync.WaitGroup,
	running func() bool,
	srcChan chan source.Event,
	_ chan model.SinkEvent,
	errChan chan source.Error,
) error {
	go func() {
		defer wg.Done()
		defer func() {
			err := src.scroll.Clear(context.TODO())
			if err != nil {
				errChan <- source.Error{
					SourceURL: src.URL(),
					Error:     err,
				}
			}
		}()

		for running() {
			const eofSleepTime = 100 * time.Millisecond

			src.searchRateLimiter.Take()
			result, err := src.scroll.DoC(context.TODO())
			if err != nil {
				if errors.Is(err, io.EOF) {
					time.Sleep(eofSleepTime)
				} else {
					errChan <- source.Error{SourceURL: src.URL(), Error: err}
				}
			} else if result != nil || result.Hits != nil || result.Hits.Hits != nil {
				for _, hit := range result.Hits.Hits {
					src.processEvent(srcChan, *hit.Source)
				}
			}
		}
	}()
	return nil
}

func (src *elasticSearch) processEvent(srcChan chan source.Event, hit json.RawMessage) {
	srcChan <- source.Event{Source: src, Event: string(hit)}
}
