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

// Package message implements functions to model/manipulate messages specified by
// the method interface of the APT software package manager. For more information
// about the APT method interface see, http://www.fifi.org/doc/libapt-pkg-doc/method.html/ch2.html#s2.3.
package message

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Header models the first line of a message specified by the APT method interface.
type Header struct {
	Status      int
	Description string
}

// Field models the lines of a message specified by the APT method interface that follow a Header line.
type Field struct {
	Name  string
	Value string
}

// Message models an entore message specified by the APT method interface. It includes a Header
// and a list of Fields.
type Message struct {
	Header *Header
	Fields []*Field
}

// FromBytes takes a byte representation of a Message and unmarshals it into a Message.
func FromBytes(b []byte) (*Message, error) {
	m := &Message{}
	err := m.unmarshalText(b)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetFieldValue returns the Value property of the Field with the given name. If no field is found
// with the given name, it returns a zero length string. This is useful for Fields that appear only
// once in a given Message.
func (m *Message) GetFieldValue(name string) string {
	for _, field := range m.Fields {
		if field.Name == name {
			return field.Value
		}
	}
	return ""
}

// GetFieldList returns a slice of Fields with the given name. This is useful when looking for a
// collection of fields with a given name from the same Message, e.g. 'Config-Item'.
func (m *Message) GetFieldList(name string) []*Field {
	fields := []*Field{}
	for _, field := range m.Fields {
		if field.Name == name {
			fields = append(fields, field)
		}
	}
	return fields
}

// String returns a string representation of a Message formatted according to the APT method interface.
func (m *Message) String() string {
	buffer := &bytes.Buffer{}
	for _, field := range m.Fields {
		buffer.WriteString(field.String())
		buffer.WriteString("\n")
	}
	return fmt.Sprintf("%s\n%s", m.Header.String(), buffer.String())
}

// String returns a string representation of a Header formatted according to the APT method interface.
func (h *Header) String() string {
	return fmt.Sprintf("%d %s", h.Status, h.Description)
}

// String returns a string representation of a Field formatted according to the APT method interface.
func (f *Field) String() string {
	return fmt.Sprintf("%s: %s", f.Name, f.Value)
}

func (m *Message) marshalText() (text []byte, err error) {
	t := m.String()
	return []byte(t), nil
}

func (m *Message) unmarshalText(text []byte) error {
	var err error
	*m, err = parse(string(text))
	return err
}

// parse splits a string message by line, and then constructs a Message
// from a Header and slice of Fields.
func parse(value string) (Message, error) {
	trimmed := strings.TrimSpace(value)
	lines := strings.Split(trimmed, "\n")
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
func parseHeader(headerLine string) (*Header, error) {
	headerLine = strings.TrimSpace(headerLine)
	headerParts := strings.Split(headerLine, " ")
	statusString := strings.TrimSpace(headerParts[0])
	statusCode, err := strconv.Atoi(statusString)
	if err != nil {
		return nil, err
	}
	descriptionParts := make([]string, len(headerParts[1:]))
	for idx, descPart := range headerParts[1:] {
		descriptionParts[idx] = strings.TrimSpace(descPart)
	}
	description := strings.Join(descriptionParts, " ")
	return &Header{Status: statusCode, Description: description}, nil
}

func parseFields(fieldLines []string) []*Field {
	fields := []*Field{}
	for _, fieldLine := range fieldLines {
		field := parseField(fieldLine)
		fields = append(fields, field)
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
func parseField(fieldLine string) *Field {
	fieldLine = strings.TrimSpace(fieldLine)
	fieldParts := strings.Split(fieldLine, ":")
	fieldName := fieldParts[0]
	valueParts := make([]string, len(fieldParts[1:]))
	for idx, valuePart := range fieldParts[1:] {
		valueParts[idx] = strings.TrimSpace(valuePart)
	}
	fieldValue := strings.Join(valueParts, ":")
	return &Field{Name: fieldName, Value: fieldValue}
}
