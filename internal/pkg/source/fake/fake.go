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
	applicationName model.ApplicationName
	style           string
	faker           *gofakeit.Faker
}

func newSource(sourceURL model.SourceURL) (*source.Source, error) {
	faker := getFakerArg(sourceURL)
	var src source.Source = fakeSource{
		BaseSource:      source.New('🃟', sourceURL),
		applicationName: model.ApplicationName(faker.AppName()),
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
//nolint:gosec
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

	messageGenerator, exceptionGenerator := src.makers()

	go func() {
		defer wg.Done()
		for running() {
			var exception string
			index := rand.Intn(len(nonErrorLevels))
			var lvl = nonErrorLevels[index]
			if rand.Intn(maxPercent) > errorThresholdPercent {
				exception = exceptionGenerator()
				lvl = level.Error
			}
			message := messageGenerator()
			var format = "{" +
				`"@timestamp":"%s",` +
				`"level":"%s",` +
				`"message":"%s",` +
				`"application":"%s",` +
				`"method":"%s",` +
				`"line_number":%d,` +
				`"thread_name":"%s",` +
				`"class":"%s"`
			if exception != "" {
				format += `,"exception":{"stacktrace":"%s"}`
			} else {
				format += "%s"
			}
			format += "}"
			event := source.Event{
				Source: src,
				Event: fmt.Sprintf(
					format,
					time.Now().Format(time.RFC3339),
					lvl,
					message,
					src.applicationName,
					src.faker.Animal(),
					//nolint:gosec
					rand.Uint64()%maxLineNumber,
					src.faker.Adverb(),
					src.faker.AppName(),
					exception,
				),
				Schema: "logstash",
			}
			srcChan <- event
		}
	}()
	return nil
}

func (src fakeSource) getSentenceMaker() func(int) string {
	var sentenceMaker func(int) string
	switch src.style {
	case hipsterStyle:
		sentenceMaker = src.faker.HipsterSentence
	case loremStyle:
		sentenceMaker = src.faker.LoremIpsumSentence
	case hackerStyle:
		sentenceMaker = func(_ int) string { return src.faker.HackerPhrase() }
	case defaultStyle:
		fallthrough
	default:
		sentenceMaker = src.faker.Sentence
	}
	return sentenceMaker
}

func (src fakeSource) makers() (func() string, func() string) {
	const paragraphLength = 4
	const sentenceLength = 10

	sentenceMaker := src.getSentenceMaker()

	var messageGenerator = func() string {
		var sentences []string
		for i := 0; i < paragraphLength; i++ {
			sentences = append(sentences, sentenceMaker(sentenceLength))
		}
		return strings.Join(sentences, "")
	}
	var exceptionGenerator = func() string {
		var sentences []string
		for i := 0; i < paragraphLength; i++ {
			sentences = append(sentences, sentenceMaker(sentenceLength))
		}
		return strings.Join(sentences, "\\n")
	}
	return messageGenerator, exceptionGenerator
}

func getFakerArg(sourceURL model.SourceURL) *gofakeit.Faker {
	var seed, _ = sourceURL.QueryInt64("seed")
	if seed == nil {
		epoch := time.Now().Unix()
		seed = &epoch
	}
	return gofakeit.New(*seed)
}
