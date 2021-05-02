package source

import (
	"couture/internal/pkg/model"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//init registers the type with the typeRegistry.
func init() {
	typeRegistry[reflect.TypeOf(Fake{})] = NewFake
	registry = append(registry, Fake{})
}

//NewFake provides a configured Fake source.
func NewFake(srcUrl url.URL) interface{} {
	return Fake{baseSource: baseSource{srcUrl: srcUrl}}
}

//Fake provides fake data.
type Fake struct {
	baseSource
}

func (source Fake) CanHandle(url url.URL) bool {
	return url.Scheme == "fake"
}

func (source Fake) String() string {
	return "fake"
}

func (source Fake) GoString() string {
	return "Ξ " + source.String()
}

func (source Fake) Poll() (model.Event, error) {
	if rand.Intn(100) >= 90 {
		return model.Event{}, model.ErrNoMoreEvents
	}
	var exception *model.Exception
	var level = []model.Level{
		model.LevelTrace,
		model.LevelDebug,
		model.LevelInfo,
	}[rand.Intn(3)]
	if rand.Intn(100) > 90 {
		stackTrace := strings.Join([]string{
			gofakeit.Sentence(10),
			gofakeit.Sentence(10),
			gofakeit.Sentence(10),
			gofakeit.Sentence(10),
		}, "\n")
		exception = &model.Exception{StackTrace: model.StackTrace(stackTrace)}
		level = model.LevelError
	}

	return model.Event{
		Timestamp:  model.Timestamp(time.Now().Truncate(time.Hour)),
		Level:      level,
		Message:    model.Message(gofakeit.HipsterParagraph(1, 4, 10, "\n")),
		MethodName: model.MethodName(gofakeit.Animal()),
		LineNumber: model.LineNumber(strconv.FormatInt(int64(rand.Intn(200)), 10)),
		ThreadName: model.ThreadName(gofakeit.Username()),
		ClassName:  model.ClassName(gofakeit.AppName()),
		Exception:  exception,
	}, nil
}
