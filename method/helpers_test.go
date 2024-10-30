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
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestS3EndpointURL(t *testing.T) {
	specs := map[string]struct {
		expectedURL *url.URL
		expectError bool
	}{
		"us-east-1": {
			&url.URL{
				Scheme: "https",
				Host:   "s3.amazonaws.com",
			},
			false,
		},
		"us-west-2": {
			&url.URL{
				Scheme: "https",
				Host:   "s3.us-west-2.amazonaws.com",
			},
			false,
		},
		"us-gov-east-1": {
			&url.URL{
				Scheme: "https",
				Host:   "s3.us-gov-east-1.amazonaws.com",
			},
			false,
		},
		"cn-north-1": {
			&url.URL{
				Scheme: "https",
				Host:   "s3.cn-north-1.amazonaws.com.cn",
			},
			false,
		},
		"outer-space-0": {
			nil,
			true,
		},
	}

	for region, spec := range specs {
		t.Run(region, func(t *testing.T) {
			u, err := s3EndpointURL(region)
			if err != nil && !spec.expectError {
				t.Errorf("expected s3EndpointURL(%#v) not to return an error but got %#v", region, err)
			} else if err == nil && spec.expectError {
				t.Errorf("expected s3EndpointURL(%#v) to return an error but got none", region)
			}

			if diff := cmp.Diff(spec.expectedURL, u); diff != "" {
				t.Errorf("s3EndpointURL(%#v) mismatch (-want +got):\n%s", region, diff)
			}
		})
	}
}
