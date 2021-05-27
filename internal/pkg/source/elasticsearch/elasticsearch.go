package elasticsearch

import (
	"bytes"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"encoding/json"
	"github.com/bnkamalesh/errors"
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
func (src elasticSearch) Start(
	wg *sync.WaitGroup,
	running func() bool,
	srcChan chan source.Event,
	errChan chan source.Error,
) error {
	go func() {
		const eofSleepTime = 100 * time.Millisecond

		defer wg.Done()
		for running() {
			result, err := src.scroll()
			if err != nil {
				if errors.Is(err, io.EOF) {
					time.Sleep(eofSleepTime)
				} else {
					errChan <- source.Error{Source: src, Error: err}
				}
			} else if result != nil || result.Hits != nil || result.Hits.Hits != nil {
				for _, hit := range result.Hits.Hits {
					src.processEvent(srcChan, errChan, *hit.Source)
				}
			}
		}
	}()
	return nil
}

func (src *elasticSearch) scroll() (*elastic.SearchResult, error) {
	result, err := src.esClient.Scroll(src.indexName).
		KeepAlive(keepAliveOneMinute).
		ScrollId(src.scrollID).
		Size(eventsPerFetch).
		Query(src.query).
		SortBy(elastic.NewFieldSort(model.TimestampField)).
		Pretty(true).
		Do()
	if err == nil {
		src.scrollID = result.ScrollId
	}
	return result, err
}

func (src *elasticSearch) processEvent(
	srcChan chan source.Event,
	errChan chan source.Error,
	hit json.RawMessage,
) {
	holder := eventHolder{}
	err := json.Unmarshal(hit, &holder)
	if err != nil {
		errChan <- source.Error{Source: src, Error: err}
		return
	}
	var evt = model.Event{}
	err = json.Unmarshal([]byte(holder.Event), &evt)
	if err != nil || evt.Timestamp == model.Timestamp(time.Time{}) {
		applicationName := model.ApplicationName(src.indexName)
		var formatted bytes.Buffer
		err := json.Indent(&formatted, []byte(holder.Event), "", "   ")
		var message = holder.Event
		if err == nil {
			message = formatted.String()
		}
		evt = model.Event{
			Timestamp:       model.Timestamp(time.Now()),
			Level:           level.Info,
			Message:         model.Message(message),
			ApplicationName: &applicationName,
			MethodName:      "search",
			LineNumber:      model.NoLineNumber,
			ThreadName:      nil,
			ClassName:       "-",
			Exception:       nil,
		}
	}
	srcChan <- source.Event{Source: src, Event: evt}
}
