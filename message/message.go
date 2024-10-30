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

// Package message implements functions to model/manipulate messages specified by
// the method interface of the APT software package manager. For more information
// about the APT method interface see, http://www.fifi.org/doc/libapt-pkg-doc/method.html/ch2.html#s2.3.
package message

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Header models the first line of a message specified by the APT method
// interface.
type Header struct {
	Status      int
	Description string
}

// Field models the lines of a message specified by the APT method interface
// that follow a Header line.
type Field struct {
	Name  string
	Value string
}

// Message models an entore message specified by the APT method interface. It
// includes a Header and a list of Fields.
type Message struct {
	Header *Header
	Fields []*Field
}

// FromBytes takes a byte representation of a Message and unmarshals it into a
// Message.
func FromBytes(b []byte) (*Message, error) {
	m, err := parse(string(b))
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetFieldValue returns the Value property of the Field with the given name.
// If no field is found with the given name, it returns a zero length string.
// This is useful for Fields that appear only once in a given Message.
func (m *Message) GetFieldValue(name string) (string, bool) {
	for _, f := range m.Fields {
		if f.Name == name {
			return f.Value, true
		}
	}
	return "", false
}

// GetFieldList returns a slice of Fields with the given name. This is useful
// when looking for a collection of fields with a given name from the same
// Message, e.g. 'Config-Item'.
func (m *Message) GetFieldList(name string) []*Field {
	fields := []*Field{}
	for _, f := range m.Fields {
		if f.Name == name {
			fields = append(fields, f)
		}
	}
	return fields
}

// String returns a string representation of a Message formatted according to
// the APT method interface.
func (m *Message) String() string {
	buf := &bytes.Buffer{}
	for _, f := range m.Fields {
		buf.WriteString(f.String())
		buf.WriteString("\n")
	}
	return fmt.Sprintf("%s\n%s", m.Header.String(), buf.String())
}

// String returns a string representation of a Header formatted according to
// the APT method interface.
func (h *Header) String() string {
	return fmt.Sprintf("%d %s", h.Status, h.Description)
}

// String returns a string representation of a Field formatted according to the
// APT method interface.
func (f *Field) String() string {
	return fmt.Sprintf("%s: %s", f.Name, f.Value)
}

var (
	errMsgMissingRequiredLines = errors.New("message missing required number of lines")
)

// parse splits a string message by line, and then constructs a Message from a
// Header and slice of Fields.
func parse(value string) (Message, error) {
	lines := strings.Split(strings.TrimSpace(value), "\n")
	if len(lines) < 2 {
		return Message{}, errMsgMissingRequiredLines
	}
	headerLine := lines[0]
	fieldLines := lines[1:]

	header, err := parseHeader(headerLine)
	if err != nil {
		return Message{}, err
	}
	fields := parseFields(fieldLines)
	return Message{Header: header, Fields: fields}, nil
}

// parseHeader splits a string header by white space and constructs a Header
// based on the status code and description.
//
// Lines might look like the following:
//
// 102 Status
// 200 URI Start
// 201 URI Done
// 601 Configuration
func parseHeader(line string) (*Header, error) {
	tokens := strings.Split(strings.TrimSpace(line), " ")
	status := strings.TrimSpace(tokens[0])
	statusCode, err := strconv.Atoi(status)
	if err != nil {
		return nil, err
	}
	descTkns := make([]string, len(tokens[1:]))
	for idx, descTkn := range tokens[1:] {
		descTkns[idx] = strings.TrimSpace(descTkn)
	}
	return &Header{Status: statusCode, Description: strings.Join(descTkns, " ")}, nil
}

func parseFields(lines []string) []*Field {
	fields := []*Field{}
	for _, l := range lines {
		fields = append(fields, parseField(l))
	}
	return fields
}

// parseField splits a string field by colon and constructs a Field based on
// the name and value.
//
// Lines might look like the following:
//
// URI:s3://my-s3-repository/project-a/dists/trusty/main/binary-amd64/Packages
// Config-Item: Aptitude::Get-Root-Command=sudo:/usr/bin/sudo
func parseField(line string) *Field {
	tokens := strings.Split(strings.TrimSpace(line), ":")

	// The line may have additional colons, so the value needs to be any tokens
	// after the first joined with a colon.
	valueTkns := make([]string, len(tokens[1:]))
	for idx, valueTkn := range tokens[1:] {
		valueTkns[idx] = strings.TrimSpace(valueTkn)
	}
	return &Field{Name: tokens[0], Value: strings.Join(valueTkns, ":")}
}
