package manager

import (
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws/cloudformation"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/internal/pkg/source/elasticsearch"
	"couture/internal/pkg/source/fake"
	"couture/internal/pkg/source/ssh"
	"couture/internal/pkg/source/tail"
)

// SourceMetadata is a list of sourceMetadata sourceMetadata.
var SourceMetadata = source.MetadataGroup{
	fake.Metadata(),
	cloudwatch.Metadata(),
	cloudformation.Metadata(),
	elasticsearch.Metadata(),
	tail.Metadata(),
	ssh.Metadata(),
}
