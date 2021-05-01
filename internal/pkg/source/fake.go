package source

import (
	"couture/internal/pkg/model"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"strconv"
	"time"
)

//NewFake provides a configured Fake source.
func NewFake(_ string) interface{} {
	return Fake{}
}

//Fake provides fake data.
type Fake struct {
}

func (f Fake) Poll() (model.Event, error) {
	if rand.Intn(100) >= 90 {
		return model.Event{}, model.ErrNoMoreEvents
	}
	var exception *model.Exception
	var level = []model.LogLevel{
		model.LevelTrace,
		model.LevelDebug,
		model.LevelInfo,
		model.LevelWarn,
		model.LevelError,
	}[rand.Intn(5)]
	if rand.Intn(100) > 95 {
		stackTrace := gofakeit.Sentence(200)
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
