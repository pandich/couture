package fake

import (
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/polling"
	"couture/pkg/model"
	"github.com/brianvoe/gofakeit/v6"
	errors2 "github.com/pkg/errors"
	"io"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

var applicationName = model.ApplicationName(gofakeit.AppName())

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Type:        reflect.TypeOf(basePollingSource{}),
		CanHandle:   func(url model.SourceURL) bool { return url.Scheme == "fake" },
		Creator:     create,
		ExampleURLs: []string{"fake://(?seed=<seed_int>)"},
	}
}

// basePollingSource provides fake test data.
type basePollingSource struct {
	polling.Source
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
func newSource(sourceURL model.SourceURL) (*basePollingSource, error) {
	seed, err := sourceURL.QueryInt64("seed")
	if err != nil {
		return nil, errors2.Wrapf(err, "could not parse seed")
	}
	if seed != nil {
		gofakeit.Seed(*seed)
	}
	return &basePollingSource{polling.New(sourceURL, time.Second)}, nil
}

// Poll ...
//nolint:gosec,gomnd
func (source basePollingSource) Poll() ([]model.Event, error) {
	if rand.Intn(100) >= 90 {
		return []model.Event{}, io.EOF
	}
	var exception *model.Exception
	var level = []model.Level{
		model.LevelTrace,
		model.LevelDebug,
		model.LevelInfo,
		model.LevelWarn,
	}[rand.Intn(4)]
	const count = 10
	if rand.Intn(100) > 90 {
		stackTrace := strings.Join([]string{
			gofakeit.Sentence(count),
			gofakeit.Sentence(count),
			gofakeit.Sentence(count),
			gofakeit.Sentence(count),
		}, "\n")
		exception = &model.Exception{StackTrace: model.StackTrace(stackTrace)}
		level = model.LevelError
	}

	threadName := model.ThreadName(gofakeit.Username())
	return []model.Event{{
		ApplicationName: &applicationName,
		Timestamp:       model.Timestamp(time.Now().Truncate(time.Hour)),
		Level:           level,
		Message:         model.Message(gofakeit.HipsterParagraph(1, 4, count, "\n")),
		MethodName:      model.MethodName(gofakeit.Animal()),
		LineNumber:      model.LineNumber(uint64(rand.Intn(200))),
		ThreadName:      &threadName,
		ClassName:       model.ClassName(gofakeit.AppName()),
		Exception:       exception,
	}}, nil
}
