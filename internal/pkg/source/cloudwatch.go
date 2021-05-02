package source

import (
	"couture/internal/pkg/model"
	"fmt"
	"net/url"
	"reflect"
)

//init registers the type with the typeRegistry.
func init() {
	typeRegistry[reflect.TypeOf(CloudWatch{})] = func(srcUrl url.URL) interface{} {
		return CloudWatch{baseSource{srcUrl: srcUrl}}
	}
	registry = append(registry, CloudWatch{})
}

//CloudWatch provides CloudWatch data.
type CloudWatch struct {
	baseSource
}

func (source CloudWatch) CanHandle(url url.URL) bool {
	return url.Scheme == "cloudwatch"
}

func (source CloudWatch) String() string {
	return source.srcUrl.Path
}

func (source CloudWatch) GoString() string {
	return "☁︎ " + source.String()
}

func (source CloudWatch) Poll() (model.Event, error) {
	// TODO implement CloudWatch source.
	return model.Event{}, fmt.Errorf("not implemented")
}
