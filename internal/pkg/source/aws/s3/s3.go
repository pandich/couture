package s3

import (
	"context"
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	errors2 "github.com/pkg/errors"
	"reflect"
	"sync"
	"time"
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
	s3     *s3.Client
	bucket string
	key    string
}

// newSource S3 source.
func newSource(_ *time.Time, sourceURL model.SourceURL) (*source.Source, error) {
	awsSource, err := aws.New('☂', &sourceURL)
	if err != nil {
		return nil, errors2.Wrapf(err, "bad S3 URL: %+v\n", sourceURL)
	}
	src := s3Source{
		Source: awsSource,
		s3:     s3.NewFromConfig(awsSource.Config()),
		bucket: sourceURL.Host,
		key:    sourceURL.Path[1:],
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
	_ chan source.Error,
) error {
	const partSize = 16 * 1024

	request := &s3.GetObjectInput{Bucket: &src.bucket, Key: &src.key}

	writer := sink.NewChanWriterAt(src, srcChan)
	downloader := manager.NewDownloader(src.s3, func(d *manager.Downloader) {
		d.PartSize = partSize
		d.Concurrency = 1
	})
	ctx := context.Background()
	if _, err := downloader.Download(ctx, writer, request); err != nil {
		return err
	}

	go func() {
		defer wg.Done()
		_, cancel := context.WithCancel(ctx)
		defer cancel()
		for running() {
		}
	}()
	return nil
}
