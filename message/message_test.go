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
	if fakeMsg != actual {
		t.Errorf("Expected m.String() to equal: %s Actual: %s", fakeMsg, actual)
	}
}

func TestMarshalUnmarshallFakeMsg(t *testing.T) {
	h := &Header{Status: 700, Description: "Fake Description"}
	f := []*Field{
		&Field{Name: "Foo", Value: "bar"},
		&Field{Name: "Baz", Value: "false"},
		&Field{Name: "Filename", Value: "apt-transport.deb"},
	}
	m := &Message{Header: h, Fields: f}
	bytes, err := m.marshalText()
	if err != nil {
		t.Error(err)
	}
	newMessage := &Message{}
	newMessage.unmarshalText(bytes)
	if newMessage.Header.Description != "Fake Description" {
		t.Error("Failed to marshal -> unmarshal successfully")
	}
}

func TestUnmarshalConfigurationMsg(t *testing.T) {
	m := &Message{}
	m.unmarshalText([]byte(configMsg))

	expectedCount := 12
	count := len(m.Fields)
	if count != expectedCount {
		t.Errorf("Expected Fields to contain %d items, but had %d", expectedCount, count)
	}

	expected := "Config-Item: APT::Architecture=amd64"
	actual := m.Fields[0].String()
	if actual != expected {
		t.Errorf("Incorrect field %s\n Expected: %s", actual, expected)
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

	actual := m.GetFieldValue("Foo")
	expected := "bar"
	if actual != expected {
		t.Errorf("Incorrect field value '%s' Expected: %s", actual, expected)
	}
	actual = m.GetFieldValue("Baz")
	expected = "qux"
	if actual != expected {
		t.Errorf("Incorrect field value '%s' Expected: %s", actual, expected)
	}
	actual = m.GetFieldValue("Filename")
	expected = "apt-transport.deb"
	if actual != expected {
		t.Errorf("Incorrect field value '%s' Expected: %s", actual, expected)
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
		t.Errorf("Incorrect number of fields '%d' Expected: %d", actualLength, expectedLength)
	}
}

func TestUnmarshalAcquireMsg(t *testing.T) {
	m := &Message{}
	m.unmarshalText([]byte(acqMsg))

	expectedCount := 2
	count := len(m.Fields)
	if count != expectedCount {
		t.Errorf("Expected Fields to contain %d items, but had %d", expectedCount, count)
	}

	field := m.GetFieldValue("Filename")
	if field != "/var/cache/apt/archives/partial/python-bernhard_0.2.3-1_all.deb" {
		t.Errorf("Incorrect field %s", field)
	}
}

func TestUnmarshalFieldsWithMissingSpaces(t *testing.T) {
	m := &Message{}
	m.unmarshalText([]byte(acqMsgNoSpaces))

	expectedCount := 4
	count := len(m.Fields)
	if count != expectedCount {
		t.Errorf("Expected Fields to contain %d items, but had %d", expectedCount, count)
	}

	field := m.Fields[1]
	if field.String() != "Filename: Packages.downloaded" {
		t.Errorf("Incorrect field %s", field)
	}
}
