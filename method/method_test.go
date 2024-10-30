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
	"strings"
	"testing"
	"time"
)

const (
	capMsg = `100 Capabilities
Send-Config: true
Pipeline: true
Single-Instance: yes
`

	// The trailing blank line is intentional
	acqMsg = `600 URI Acquire
URI: s3://fake-access-key-id:fake-access-key-secret@s3.amazonaws.com/apt-repo-bucket/apt/generic/python-bernhard_0.2.3-1_all.deb
Filename: /tmp/python-bernhard_0.2.3-1_all.deb

600 URI Acquire
URI: s3://fake-access-key-id:fake-access-key-secret@s3.amazonaws.com/apt-repo-bucket/apt/generic/riemann-sumd_0.7.2-1_all.deb
Filename: /tmp/riemann-sumd_0.7.2-1_all.deb

`

	// The trailing blank line is intentional
	configMsg = `601 Configuration
Config-Item: Dir::Log=var/log/apt
Config-Item: Dir::Log::Terminal=term.log
Config-Item: Dir::Log::History=history.log
Config-Item: Dir::Ignore-Files-Silently::=~$
Config-Item: Acquire::cdrom::mount=/media/cdrom
Config-Item: Acquire::s3::region=us-east-2
Config-Item: Aptitude::Get-Root-Command=sudo:/usr/bin/sudo
Config-Item: Unattended-Upgrade::Allowed-Origins::=${distro_id}:${distro_codename}-security

`
)

func TestCapabilities(t *testing.T) {
	actual := capabilities().String()
	if actual != capMsg {
		t.Errorf("capabilities() = %s; expected %s", actual, capMsg)
	}
}

func TestReadInputFinishes(t *testing.T) {
	reader := strings.NewReader(acqMsg)
	method := New()
	go method.readInput(reader)

	msgs := 0
loop:
	for {
		select {
		case <-method.msgChan:
			msgs++
		case <-time.After(10 * time.Millisecond):
			if reader.Len() > 0 {
				t.Errorf("Found reader with %d bytes; expected reader to be empty", reader.Len())
			}
			break loop
		}
	}

	if msgs != 2 {
		t.Errorf("Found %d messages; expected %d", msgs, 2)
	}
}

func TestSettingRegion(t *testing.T) {
	reader := strings.NewReader(configMsg)
	method := New()
	go method.readInput(reader)

	//consume the messages on the channel
	for {
		bytes := <-method.msgChan
		method.handleBytes(bytes)
		if reader.Len() == 0 {
			break
		}
	}
	expected := "us-east-2"
	if method.region != expected {
		t.Errorf("method.region = %s; expected %s", method.region, expected)
	}
}

type locTest struct {
	url             string
	accessKey       string
	accessKeySecret string
}

var locTests = []locTest{
	{
		"s3://fake-access-key-id:fake-access-key-secret@s3.amazonaws.com/apt-repo-bucket/apt/generic/python-bernhard_0.2.3-1_all.deb",
		"fake-access-key-id",
		"fake-access-key-secret",
	},
	{
		"s3://fake-ac/cess-key-id:fake-ac/cess-key-secret@s3.amazonaws.com/apt-repo-bucket/apt/generic/python-bernhard_0.2.3-1_all.deb",
		"fake-ac/cess-key-id",
		"fake-ac/cess-key-secret", // secret contains a forward slash
	},
	{
		"s3://fake-ac%2Fcess-key-id:fake-ac%2Fcess-key-secret@s3.amazonaws.com/apt-repo-bucket/apt/generic/python-bernhard_0.2.3-1_all.deb",
		"fake-ac/cess-key-id",     // access key contains a forward slash that was encoded as %2F in the original url
		"fake-ac/cess-key-secret", // secret contains a forward slash that was encoded as %2F in the original url
	},
	{
		"s3://fake-access-key-id:@s3.amazonaws.com/apt-repo-bucket/apt/generic/python-bernhard_0.2.3-1_all.deb",
		"fake-access-key-id",
		"", // secret is blank
	},
	{
		"s3://:fake-access-key-secret@s3.amazonaws.com/apt-repo-bucket/apt/generic/python-bernhard_0.2.3-1_all.deb",
		"", // access key is blank
		"fake-access-key-secret",
	},
}

func TestCreateLocation(t *testing.T) {
	for _, lt := range locTests {
		l, err := newLocation(lt.url, "s3.amazonaws.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if l.uri.User.Username() != lt.accessKey {
			t.Errorf("unexpected accessKey: got %s, want %s", l.uri.User.Username(), lt.accessKey)
		}
		pass, _ := l.uri.User.Password()
		if pass != lt.accessKeySecret {
			t.Errorf("unexpected accessKeySecret: got %s, want %s", pass, lt.accessKeySecret)
		}
	}
}
