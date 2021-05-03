package elasticsearch

import (
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/polling"
	"couture/pkg/model"
	"encoding/json"
	"fmt"
	errors2 "github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v3"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// TODO only partially works
const eventsPerFetch = 100

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Type: reflect.TypeOf(elasticsearchSource{}),
		CanHandle: func(url model.SourceURL) bool {
			_, ok := map[string]bool{
				scheme + "+http":            true,
				schemeAliasShort + "+http":  true,
				scheme + "+https":           true,
				schemeAliasShort + "+https": true,
			}[url.Scheme]
			return ok
		},
		Creator: create,
		ExampleURLs: []string{
			"elasticsearch://...",
			"es://...",
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

// elasticsearchSource provides elasticsearch test data.
type elasticsearchSource struct {
	polling.Source
	query     string
	scrollID  string
	indexName string
	esClient  *elastic.Client
}

// create CloudFormation source casted to an *interface{}.
func create(sourceURL model.SourceURL) (*interface{}, error) {
	src, err := newSource(sourceURL)
	if err != nil {
		return nil, err
	}
	var i interface{} = src
	return &i, nil
}

// newSource ...
func newSource(sourceURL model.SourceURL) (*elasticsearchSource, error) {
	sourceURL.Scheme = strings.Replace(sourceURL.Scheme, scheme+"+", "", 1)
	sourceURL.Scheme = strings.Replace(sourceURL.Scheme, schemeAliasShort+"+", "", 1)
	esClient, err := elastic.NewClient(elastic.SetURL(sourceURL.Scheme + "://" + sourceURL.Host))
	if err != nil {
		return nil, err
	}
	u := url.URL(sourceURL)
	var queryPieces []string
	for k, v := range u.Query() {
		queryPieces = append(queryPieces, fmt.Sprintf(`%s=(%s)`, k, strings.Join(v, "|")))
	}
	query := fmt.Sprintf(
		`{"query":{"query_string":{"query":%s}}}`,
		strconv.Quote(strings.Join(queryPieces, " ")),
	)
	return &elasticsearchSource{
		Source:    polling.New(sourceURL, time.Second),
		esClient:  esClient,
		query:     query,
		indexName: strings.Trim(sourceURL.Path, "/"),
	}, nil
}

// Poll ...
func (source *elasticsearchSource) Poll() ([]model.Event, error) {
	result, err := source.esClient.Scroll(source.indexName).
		KeepAlive(keepAliveOneMinute).
		ScrollId(source.scrollID).
		Size(eventsPerFetch).
		Query(elastic.NewQueryStringQuery(source.URL().RawQuery)).
		SortBy(elastic.NewFieldSort(model.TimestampField)).
		Pretty(true).
		Do()
	if err != nil {
		return nil, errors2.Wrapf(err, "name=%s", source.indexName)
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
			return nil, err
		}
		events = append(events, evt)
	}
	source.scrollID = result.ScrollId
	return events, nil
}
