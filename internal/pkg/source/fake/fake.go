//nolint:gomnd,gosec
package fake

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	hipsterStyle = "hipster"
	loremStyle   = "lorem"
	hackerStyle  = "hacker"
	defaultStyle = "default"
)

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Name:        "Fake",
		Type:        reflect.TypeOf(fakeSource{}),
		CanHandle:   func(url model.SourceURL) bool { return url.Scheme == "fake" },
		Creator:     newSource,
		ExampleURLs: []string{},
	}
}

// fakeSource provides fake test data.
type fakeSource struct {
	source.BaseSource
	applicationName model.Application
	style           string
	faker           *gofakeit.Faker
}

func newSource(_ *time.Time, sourceURL model.SourceURL) (*source.Source, error) {
	faker := getFakerArg(sourceURL)
	var src source.Source = fakeSource{
		BaseSource:      source.New('🃟', sourceURL),
		applicationName: model.Application(faker.AppName()),
		style:           getStyleArg(sourceURL),
		faker:           faker,
	}
	return &src, nil
}

func getStyleArg(sourceURL model.SourceURL) string {
	var style, ok = sourceURL.QueryKey("style")
	if !ok {
		style = defaultStyle
	}
	return style
}

// Start ...
func (src fakeSource) Start(
	wg *sync.WaitGroup,
	running func() bool,
	srcChan chan source.Event,
	_ chan model.SinkEvent,
	_ chan source.Error,
) error {
	const maxPercent = 100
	const errorThresholdPercent = 90
	const maxLineNumber = 200
	nonErrorLevels := []level.Level{level.Trace, level.Debug, level.Info, level.Warn}

	messageGenerator, errorGenerator := src.textGenerators()
	jsonMessageGenerater, jsonErrorGenerator := src.jsonGenerators()

	go func() {
		defer wg.Done()
		for running() {
			var errString string
			index := rand.Intn(len(nonErrorLevels))
			var lvl = nonErrorLevels[index]
			if rand.Intn(maxPercent) > errorThresholdPercent {
				if rand.Intn(10) > 2 {
					errString = jsonErrorGenerator()
				} else {
					errString = errorGenerator()
				}
				lvl = level.Error
			}
			var message string
			if rand.Intn(10) > 2 {
				message = jsonMessageGenerater()
			} else {
				message = messageGenerator()
			}
			hasError := errString != ""
			event := source.Event{
				Source: src,
				Event: fmt.Sprintf(
					src.getFormat(hasError),
					time.Now().Format(time.RFC3339),
					lvl,
					message,
					src.applicationName,
					src.faker.Animal(),
					rand.Uint64()%maxLineNumber,
					src.faker.Adverb(),
					src.faker.AppName(),
					errString,
				),
			}
			srcChan <- event
		}
	}()
	return nil
}

func (src fakeSource) getSentenceGenerator() func(int) string {
	var sentenceGenerator func(int) string
	switch src.style {
	case hipsterStyle:
		sentenceGenerator = src.faker.HipsterSentence
	case loremStyle:
		sentenceGenerator = src.faker.LoremIpsumSentence
	case hackerStyle:
		sentenceGenerator = func(_ int) string { return src.faker.HackerPhrase() }
	case defaultStyle:
		fallthrough
	default:
		sentenceGenerator = src.faker.Sentence
	}
	return sentenceGenerator
}

func (src fakeSource) jsonGenerators() (func() string, func() string) {
	messageGenerator := func() string {
		ba, err := src.faker.JSON(&gofakeit.JSONOptions{
			Type:     "object",
			RowCount: 10,
			Fields: []gofakeit.Field{
				{Name: "first_name", Function: "firstname"},
				{Name: "last_name", Function: "lastname"},
				{Name: "email", Function: "email"},
			},
		})
		if err != nil {
			panic(err)
		}

		return strings.ReplaceAll(string(ba), `"`, `\"`)
	}
	errorGenerator := func() string {
		ba, err := src.faker.JSON(&gofakeit.JSONOptions{
			Type:     "object",
			RowCount: 10,
			Fields: []gofakeit.Field{
				{Name: "first_name", Function: "firstname"},
				{Name: "last_name", Function: "lastname"},
				{Name: "email", Function: "email"},
			},
		})
		if err != nil {
			panic(err)
		}
		return strings.ReplaceAll(string(ba), `"`, `\"`)
	}
	return messageGenerator, errorGenerator
}

func (src fakeSource) getFormat(hasError bool) string {
	var format = "{" +
		`"@version":1,` +
		`"@timestamp":"%s",` +
		`"level":"%s",` +
		`"message":"%s",` +
		`"application":"%s",` +
		`"method":"%s",` +
		`"line_number":%d,` +
		`"thread_name":"%s",` +
		`"class":"%s"`
	if hasError {
		format += `,"exception":{"stacktrace":"%s"}`
	} else {
		format += "%s"
	}
	format += "}"
	return format
}
func (src fakeSource) textGenerators() (func() string, func() string) {
	const paragraphLength = 4
	const sentenceLength = 10

	sentenceGenerator := src.getSentenceGenerator()

	var messageGenerator = func() string {
		var sentences []string
		for i := 0; i < paragraphLength; i++ {
			sentences = append(sentences, sentenceGenerator(sentenceLength))
		}
		return strings.Join(sentences, "")
	}
	var errorGenerator = func() string {
		var sentences []string
		for i := 0; i < paragraphLength; i++ {
			sentences = append(sentences, sentenceGenerator(sentenceLength))
		}
		return strings.Join(sentences, "\\n")
	}
	return messageGenerator, errorGenerator
}

func getFakerArg(sourceURL model.SourceURL) *gofakeit.Faker {
	var seed, _ = sourceURL.QueryInt64("seed")
	if seed == nil {
		epoch := time.Now().Unix()
		seed = &epoch
	}
	return gofakeit.New(*seed)
}
