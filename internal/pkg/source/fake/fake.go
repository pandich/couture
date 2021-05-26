package fake

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"github.com/brianvoe/gofakeit/v6"
	errors2 "github.com/pkg/errors"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	hipsterStyle = "hipster"
	loremStyle   = "lorem"
	defaultStyle = "default"
)

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Name:      "Fake",
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
	source.BaseSource
	applicationName  model.ApplicationName
	messageGenerator string
}

// newSource ...
func newSource(sourceURL model.SourceURL) (*source.Source, error) {
	// seed
	seed, err := sourceURL.QueryInt64("seed")
	if err != nil {
		return nil, errors2.Wrapf(err, "could not parse seed\n")
	}
	if seed != nil {
		gofakeit.Seed(*seed)
	}

	// source
	var messageGenerator, ok = sourceURL.QueryKey("style")
	if !ok {
		messageGenerator = defaultStyle
	}

	var src source.Source = fakeSource{
		BaseSource:       source.New('🃟', sourceURL),
		applicationName:  model.ApplicationName(gofakeit.AppName()),
		messageGenerator: messageGenerator,
	}
	return &src, nil
}

// Start ...
//nolint:gosec,gomnd
func (src fakeSource) Start(wg *sync.WaitGroup, running func() bool, out chan source.Event) error {
	go func() {
		defer wg.Done()
		for running() {
			if rand.Intn(100) >= 90 {
				continue
			}
			var exception *model.Exception
			var lvl = []level.Level{level.Trace, level.Debug, level.Info, level.Warn}[rand.Intn(4)]
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
			var message string
			switch src.messageGenerator {
			case hipsterStyle:
				message = gofakeit.HipsterParagraph(1, 4, count, "\n")
			case loremStyle:
				message = gofakeit.LoremIpsumParagraph(1, 4, count, "\n")
			case defaultStyle:
				fallthrough
			default:
				message = gofakeit.Paragraph(1, 4, count, "\n")
			}
			out <- source.Event{
				Source: src,
				Event: model.Event{
					ApplicationName: &src.applicationName,
					Timestamp:       model.Timestamp(time.Now()),
					Level:           lvl,
					Message:         model.Message(message),
					MethodName:      model.MethodName(gofakeit.Animal()),
					LineNumber:      model.LineNumber(uint64(rand.Intn(200))),
					ThreadName:      &threadName,
					ClassName:       model.ClassName(gofakeit.AppName()),
					Exception:       exception,
				},
			}
		}
	}()
	return nil
}
