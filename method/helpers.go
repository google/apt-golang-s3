// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
