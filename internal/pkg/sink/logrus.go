package sink

import (
	"couture/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

//NewLogrus provides a configured Logrus sink.
func NewLogrus(_ string) interface{} {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        time.RFC3339,
		PadLevelText:           true,
		DisableLevelTruncation: true,
	})
	log.SetLevel(logrus.TraceLevel)
	log.SetOutput(os.Stdout)
	return Logrus{log: log}
}

//Logrus uses the Logrus logger to display events.
type Logrus struct {
	log *logrus.Logger
}

func (s Logrus) Accept(event model.Event) {
	var level logrus.Level
	var err error
	level, err = logrus.ParseLevel(string(event.Level))
	if err != nil {
		level = logrus.InfoLevel
	}
	s.log.Log(level, event)
}
