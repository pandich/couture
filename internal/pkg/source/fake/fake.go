package fake

import (
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"couture/pkg/model/level"
	"github.com/brianvoe/gofakeit/v6"
	errors2 "github.com/pkg/errors"
	"io"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Type:      reflect.TypeOf(fakeSource{}),
		CanHandle: func(url model.SourceURL) bool { return url.Scheme == "fake" },
		Creator: func(sourceURL model.SourceURL) (*interface{}, error) {
			src, err := newSource(sourceURL)
			if err != nil {
				return nil, err
			}
			var i interface{} = src
			return &i, nil
		},
		ExampleURLs: []string{},
	}
}

// fakeSource provides fake test data.
type fakeSource struct {
	source.Polling
	applicationName model.ApplicationName
}

// newSource ...
func newSource(sourceURL model.SourceURL) (*source.Pollable, error) {
	seed, err := sourceURL.QueryInt64("seed")
	if err != nil {
		return nil, errors2.Wrapf(err, "could not parse seed\n")
	}
	if seed != nil {
		gofakeit.Seed(*seed)
	}
	var src source.Pollable = fakeSource{
		Polling: source.NewPollable(
			sourceURL,
			time.Second,
		),
		applicationName: model.ApplicationName(gofakeit.AppName()),
	}
	return &src, nil
}

// Poll ...
//nolint:gosec,gomnd
func (src fakeSource) Poll() ([]model.Event, error) {
	if rand.Intn(100) >= 90 {
		return []model.Event{}, io.EOF
	}
	var exception *model.Exception
	var lvl = []level.Level{
		level.Trace,
		level.Debug,
		level.Info,
		level.Warn,
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
		lvl = level.Error
	}

	threadName := model.ThreadName(gofakeit.Username())
	return []model.Event{{
		ApplicationName: &src.applicationName,
		Timestamp:       model.Timestamp(time.Now().Truncate(time.Hour)),
		Level:           lvl,
		Message:         model.Message(gofakeit.HipsterParagraph(1, 4, count, "\n")),
		MethodName:      model.MethodName(gofakeit.Animal()),
		LineNumber:      model.LineNumber(uint64(rand.Intn(200))),
		ThreadName:      &threadName,
		ClassName:       model.ClassName(gofakeit.AppName()),
		Exception:       exception,
	}}, nil
}
