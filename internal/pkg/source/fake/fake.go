package fake

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/source"
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
		Name:      "Fake",
		Type:      reflect.TypeOf(fakeSource{}),
		CanHandle: func(url model.SourceURL) bool { return url.Scheme == "fake" },
		Creator: func(sourceURL model.SourceURL) (*interface{}, error) {
			var i interface{} = newSource(sourceURL)
			return &i, nil
		},
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

func newSource(sourceURL model.SourceURL) *source.Source {
	faker := getFakerArg(sourceURL)
	var src source.Source = fakeSource{
		BaseSource:      source.New('🃟', sourceURL),
		applicationName: model.ApplicationName(faker.AppName()),
		style:           getStyleArg(sourceURL),
		faker:           faker,
	}
	return &src
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
			var exception *model.Exception
			//nolint:gosec
			index := rand.Intn(len(nonErrorLevels))
			var lvl = nonErrorLevels[index]
			//nolint:gosec
			if rand.Intn(maxPercent) > errorThresholdPercent {
				stackTrace := exceptionGenerator()
				exception = &model.Exception{StackTrace: model.StackTrace(stackTrace)}
				lvl = level.Error
			}
			threadName := model.ThreadName(src.faker.Username())
			message := messageGenerator()
			srcChan <- source.Event{
				Source: src,
				Event: model.Event{
					ApplicationName: &src.applicationName,
					Timestamp:       model.Timestamp(time.Now()),
					Level:           lvl,
					Message:         model.Message(message),
					MethodName:      model.MethodName(src.faker.Animal()),
					//nolint:gosec
					LineNumber: model.LineNumber(rand.Intn(maxLineNumber)),
					ThreadName: &threadName,
					ClassName:  model.ClassName(src.faker.AppName()),
					Exception:  exception,
				},
			}
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
		return strings.Join(sentences, "\n")
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
