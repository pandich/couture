package elasticsearch

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"couture/pkg/model/level"
	"encoding/json"
	"errors"
	errors2 "github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v3"
	"io"
	"reflect"
	"strings"
	"time"
)

const eventsPerFetch = 100

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
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
	source.Polling
	query     elastic.Query
	scrollID  string
	indexName string
	esClient  *elastic.Client
}

// newSource ...
func newSource(sourceURL model.SourceURL) (*source.Pollable, error) {
	normalizeURL(&sourceURL)

	esClient, err := newElasticsearchClient(sourceURL)
	if err != nil {
		return nil, err
	}

	indexName := strings.Trim(sourceURL.Path, "/")
	query := elastic.NewQueryStringQuery(sourceURL.RawQuery)

	var src source.Pollable = elasticSearch{
		Polling:   source.NewPollable(sourceURL, time.Second),
		esClient:  esClient,
		query:     query,
		indexName: indexName,
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

// Poll ...
func (source elasticSearch) Poll() ([]model.Event, error) {
	result, err := source.esClient.Scroll(source.indexName).
		KeepAlive(keepAliveOneMinute).
		ScrollId(source.scrollID).
		Size(eventsPerFetch).
		Query(source.query).
		SortBy(elastic.NewFieldSort(model.TimestampField)).
		Pretty(true).
		Do()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return []model.Event{}, nil
		}
		return nil, errors2.Wrapf(err, "name=%s\n", source.indexName)
	}

	var events []model.Event
	for _, hit := range result.Hits.Hits {
		holder := eventHolder{}
		err = json.Unmarshal(*hit.Source, &holder)
		if err != nil {
			return nil, err
		}
		evt := model.Event{}
		err = json.Unmarshal([]byte(holder.Event), &evt)
		if err != nil {
			return []model.Event{
				{
					Timestamp:       model.Timestamp(time.Now()),
					Level:           level.Info,
					Message:         model.Message(holder.Event),
					ApplicationName: nil,
					MethodName:      "-",
					LineNumber:      model.NoLineNumber,
					ThreadName:      nil,
					ClassName:       "-",
					Exception:       nil,
				},
			}, err
		}
		events = append(events, evt)
	}
	source.scrollID = result.ScrollId
	return events, nil
}
