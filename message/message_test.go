// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package message

import (
	"testing"
)

const (
	fakeMsg = `700 Fake Description
Foo: bar
Baz: false
Filename: apt-transport.deb
`
	configMsg = `601 Configuration
Config-Item: APT::Architecture=amd64
Config-Item: APT::Build-Essential::=build-essential
Config-Item: APT::Color::Highlight=%1b[32m
Config-Item: APT::Update::Post-Invoke-Success::=touch%20/var/lib/apt/periodic/update-success-stamp%202>/dev/null%20||%20true
Config-Item: Acquire::cdrom::mount=/media/cdrom
Config-Item: Aptitude::Get-Root-Command=sudo:/usr/bin/sudo
Config-Item: CommandLine::AsString=apt-get%20install%20riemann-sumd
Config-Item: DPkg::Post-Invoke::=if%20[%20-d%20/var/lib/update-notifier%20];%20then%20touch%20/var/lib/update-notifier/dpkg-run-stamp;%20fi;%20if%20[%20-e%20/var/lib/update-notifier/updates-available%20];%20then%20echo%20>%20/var/lib/update-notifier/updates-available;%20fi%20
Config-Item: DPkg::Pre-Install-Pkgs::=/usr/sbin/dpkg-preconfigure%20--apt%20||%20true
Config-Item: Dir::State=var/lib/apt/
Config-Item: Dir=/
Config-Item: Unattended-Upgrade::Allowed-Origins::=${distro_id}:${distro_codename}-security
`
	acqMsg = `600 URI Acquire
URI: s3://fake-key-id:fake-key-secret@s3.amazonaws.com/my-fake-s3-bucket/apt/generic/python-bernhard_0.2.3-1_all.deb
Filename: /var/cache/apt/archives/partial/python-bernhard_0.2.3-1_all.deb
`
	acqMsgNoSpaces = `600 URI Acquire
URI:s3://my-s3-repository/project-a/dists/trusty/main/binary-amd64/Packages
Filename:Packages.downloaded
Fail-Ignore:true
Index-File:true
`
)

func TestMessageString(t *testing.T) {
	h := &Header{Status: 700, Description: "Fake Description"}
	f := []*Field{
		&Field{Name: "Foo", Value: "bar"},
		&Field{Name: "Baz", Value: "false"},
		&Field{Name: "Filename", Value: "apt-transport.deb"},
	}
	m := &Message{Header: h, Fields: f}

	actual := m.String()
	if actual != fakeMsg {
		t.Errorf("m.String() = %s; expected %s", actual, fakeMsg)
	}
}

func TestParseConfigurationMsg(t *testing.T) {
	m, err := FromBytes([]byte(configMsg))
	if err != nil {
		t.Errorf("Failed to parse %s into a message", configMsg)
	}

	expectedCount := 12
	count := len(m.Fields)
	if count != expectedCount {
		t.Errorf("Expected Fields to contain %d items, but had %d", expectedCount, count)
	}

	expected := "Config-Item: APT::Architecture=amd64"
	actual := m.Fields[0].String()
	if actual != expected {
		t.Errorf("String() = %s;  expected: %s", actual, expected)
	}
}

func TestGetFieldValue(t *testing.T) {
	h := &Header{Status: 700, Description: "Fake Description"}
	f := []*Field{
		&Field{Name: "Foo", Value: "bar"},
		&Field{Name: "Baz", Value: "qux"},
		&Field{Name: "Filename", Value: "apt-transport.deb"},
	}
	m := &Message{Header: h, Fields: f}

	actual, _ := m.GetFieldValue("Foo")
	expected := "bar"
	if actual != expected {
		t.Errorf("m.GetFieldValue(\"Foo\") = %s; expected: %s", actual, expected)
	}
	actual, _ = m.GetFieldValue("Baz")
	expected = "qux"
	if actual != expected {
		t.Errorf("m.GetFieldValue(\"Baz\") = %s; expected: %s", actual, expected)
	}
	actual, _ = m.GetFieldValue("Filename")
	expected = "apt-transport.deb"
	if actual != expected {
		t.Errorf("m.GetFieldValue(\"Filename\") = %s; expected: %s", actual, expected)
	}
}

func TestGetFieldList(t *testing.T) {
	h := &Header{Status: 700, Description: "Fake Description"}
	f := []*Field{
		&Field{Name: "Config-Item", Value: "bar"},
		&Field{Name: "Config-Item", Value: "qux"},
		&Field{Name: "Filename", Value: "apt-transport.deb"},
	}
	m := &Message{Header: h, Fields: f}

	actualLength := len(m.GetFieldList("Config-Item"))
	expectedLength := 2
	if actualLength != expectedLength {
		t.Errorf("Incorrect number of fields '%d' expected: %d", actualLength, expectedLength)
	}
}

func TestParseAcquireMsg(t *testing.T) {
	m, err := FromBytes([]byte(acqMsg))
	if err != nil {
		t.Errorf("Failed to parse %s into a message", acqMsg)
	}

	expectedCount := 2
	count := len(m.Fields)
	if count != expectedCount {
		t.Errorf("Found %d fields; expected %d", count, expectedCount)
	}

	status := m.Header.Status
	expected := 600
	if status != expected {
		t.Errorf("Status = %d; expected %d", status, expected)
	}

	description := m.Header.Description
	expectedDesc := "URI Acquire"
	if description != expectedDesc {
		t.Errorf("Description = %s; expected %s", description, expectedDesc)
	}

	value, _ := m.GetFieldValue("Filename")
	expectedVal := "/var/cache/apt/archives/partial/python-bernhard_0.2.3-1_all.deb"

	if value != expectedVal {
		t.Errorf("m.GetFieldValue(\"Filename\") = %s; expected %s", value, expectedVal)
	}
}

func TestParseFieldsWithMissingSpaces(t *testing.T) {
	m, err := FromBytes([]byte(acqMsgNoSpaces))
	if err != nil {
		t.Errorf("Failed to parse %s into a message", acqMsgNoSpaces)
	}

	count := len(m.Fields)
	expected := 4
	if count != expected {
		t.Errorf("len(m.Fields) = %d; expected %d", count, expected)
	}

	field := m.Fields[0]
	expectedName := "URI"
	expectedVal := "s3://my-s3-repository/project-a/dists/trusty/main/binary-amd64/Packages"
	if field.Name != expectedName {
		t.Errorf("field.Name = %s; expected %s", field.Name, expectedName)
	}
	if field.Value != expectedVal {
		t.Errorf("field.Value = %s; expected %s", field.Value, expectedVal)
	}
}
