package sink

import (
	"github.com/pandich/couture/source"
	"io"
	"strings"
	"sync"
)

// NewChanWriterAt ...
func NewChanWriterAt(
	src source.Source,
	out chan source.Event,
) io.WriterAt {
	return &chanWriteAt{
		src:  src,
		out:  out,
		lock: sync.Mutex{},
	}
}

// chanWriteAt ...
type chanWriteAt struct {
	src       source.Source
	out       chan source.Event
	remainder string
	lock      sync.Mutex
}

// WriteAt ...
func (writer *chanWriteAt) WriteAt(buf []byte, _ int64) (n int, err error) {
	if len(buf) == 0 {
		return 0, nil
	}
	writer.lock.Lock()
	defer writer.lock.Unlock()

	s := writer.remainder + string(buf)
	writer.remainder = ""
	var pieces = strings.Split(s, "\n")

	if !strings.HasSuffix(s, "\n") {
		pieces, writer.remainder = pieces[0:len(pieces)-1], pieces[len(pieces)-1]
	}
	for _, s := range pieces {
		writer.out <- source.Event{Source: writer.src, Event: s}
	}
	return len(buf), nil
}
