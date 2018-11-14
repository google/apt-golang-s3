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
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Header struct {
	Status      int
	Description string
}

type Field struct {
	Name  string
	Value string
}

type Message struct {
	Header *Header
	Fields []*Field
}

func FromBytes(b []byte) (*Message, error) {
	m := &Message{}
	err := m.unmarshalText(b)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Message) GetFieldValue(name string) string {
	for _, field := range m.Fields {
		if field.Name == name {
			return field.Value
		}
	}
	return ""
}

func (m *Message) GetFieldList(name string) []*Field {
	fields := []*Field{}
	for _, field := range m.Fields {
		if field.Name == name {
			fields = append(fields, field)
		}
	}
	return fields
}

func (m *Message) String() string {
	buffer := &bytes.Buffer{}
	for _, field := range m.Fields {
		buffer.WriteString(field.String())
		buffer.WriteString("\n")
	}
	return fmt.Sprintf("%s\n%s", m.Header.String(), buffer.String())
}

func (h *Header) String() string {
	return fmt.Sprintf("%d %s", h.Status, h.Description)
}

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

// Lines might look like the following
//102 Status
//200 URI Start
//201 URI Done
//601 Configuration
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

// Lines might look like the following
//URI:s3://my-s3-repository/project-a/dists/trusty/main/binary-amd64/Packages
//Config-Item: Aptitude::Get-Root-Command=sudo:/usr/bin/sudo
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
