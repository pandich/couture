package s3

import (
	"context"
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	errors2 "github.com/pkg/errors"
	"reflect"
	"strings"
	"sync"
)

const scheme = "s3"

// Metadata ...
func Metadata() source.Metadata {
	return source.Metadata{
		Name: "AWS S3",
		Type: reflect.TypeOf(s3Source{}),
		CanHandle: func(url model.SourceURL) bool {
			_, ok := map[string]bool{scheme: true}[url.Scheme]
			return ok
		},
		Creator:     newSource,
		ExampleURLs: []string{fmt.Sprintf("%s://<path>", scheme)},
	}
}

// s3Source a S3 log poller.
type s3Source struct {
	aws.Source
	s3 *s3.Client
}

// newSource S3 source.
func newSource(sourceURL model.SourceURL) (*source.Source, error) {
	awsSource, err := aws.New('☂', &sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad S3 URL: %+v\n", sourceURL)
	}
	src := s3Source{
		Source: awsSource,
		s3:     s3.NewFromConfig(awsSource.Config()),
	}
	var p source.Source = &src
	return &p, nil
}

// Start ...
func (src *s3Source) Start(
	wg *sync.WaitGroup,
	running func() bool,
	srcChan chan source.Event,
	_ chan model.SinkEvent,
	errChan chan source.Error,
) error {
	writer := chanWriteAt{
		src:  src,
		out:  srcChan,
		lock: sync.Mutex{},
	}
	downloader := manager.NewDownloader(src.s3, func(d *manager.Downloader) {
		const partSize = 16 * 1024
		d.PartSize = partSize
		d.Concurrency = 1
	})
	bucket := src.URL().Host
	key := src.URL().Path[1:]
	go func() {
		defer wg.Done()
		for running() {
			if _, err := downloader.Download(context.TODO(), &writer, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &key,
			}); err != nil {
				errChan <- source.Error{
					SourceURL: src.URL(),
					Error:     err,
				}
			}
		}
	}()
	return nil
}

type chanWriteAt struct {
	src       source.Source
	out       chan source.Event
	remainder string
	lock      sync.Mutex
}

// WriteAt ...
func (b *chanWriteAt) WriteAt(p []byte, _ int64) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if len(p) == 0 {
		return 0, nil
	}
	s := b.remainder + string(p)
	b.remainder = ""
	var pieces = strings.Split(s, "\n")
	if s[len(s)-1] != '\n' {
		pieces, b.remainder = pieces[0:len(pieces)-1], pieces[len(pieces)-1]
	}
	for _, s := range pieces {
		b.out <- source.Event{
			Source: b.src,
			Event:  s,
		}
	}
	return len(p), nil
}
