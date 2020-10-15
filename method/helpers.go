package method

import (
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

func s3EndpointURL(region string) (*url.URL, error) {
	resolver := endpoints.DefaultResolver()

	endpoint, err := resolver.EndpointFor(endpoints.S3ServiceID, region, endpoints.StrictMatchingOption)
	if err != nil {
		return nil, fmt.Errorf("resolving S3 endpoint for region %s: %w", region, err)
	}

	return url.Parse(endpoint.URL)
}
