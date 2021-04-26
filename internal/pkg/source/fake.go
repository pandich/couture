package source

import (
	"couture/pkg/couture/model"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"time"
)

var (
	FakeSource Source = fakeSource{}
)

type (
	fakeSource struct {
	}
)

func (f fakeSource) ProvideEvent() (*model.Event, error) {
	if rand.Intn(100) >= 90 {
		return nil, nil
	}
	methodName := model.MethodName(gofakeit.Animal())
	lineNumber := model.LineNumber(rand.Intn(200))
	threadName := model.ThreadName(gofakeit.Username())
	className := model.ClassName(gofakeit.AppName())
	var stackTrace *model.StackTrace = nil
	if rand.Intn(1) == 1 {
		s := model.StackTrace(gofakeit.HipsterSentence(5))
		stackTrace = &s
	}
	return model.NewEvent(
		time.Now(),
		model.LevelInfo,
		model.Message(gofakeit.Sentence(50)),
		&methodName,
		&lineNumber,
		&threadName,
		&className,
		stackTrace,
	), nil
}
