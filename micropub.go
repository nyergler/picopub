//go:generate jsonenums -type=EntryType
//go:generate stringer -type=EntryType

package picopub

import (
	"encoding/json"
	"net/url"
	"strings"
)

type EntryType int

const (
	Unknown EntryType = iota
	Entry
)

var (
	_EntryTypeNameToValue = map[string]EntryType{
		"Unknown": Unknown,
		"Entry":   Entry,
		"H-Entry": Entry,
		"h-entry": Entry,
	}

	_EntryTypeValueToName = map[EntryType]string{
		Unknown: "Unknown",
		Entry:   "Entry",
	}
)

func FromString(s string) EntryType {
	t, ok := _EntryTypeNameToValue[s]
	if ok {
		return t
	}
	t, ok = _EntryTypeNameToValue[strings.Title(s)]
	if ok {
		return t
	}
	return Unknown
}

type NestedProperty interface {
	// either array of strings or array of mf objects
	Type() EntryType
	Content() []string
	// Get() like url.Values?
	Object() *MicroformatObject
}

type Property struct {
	entryType  EntryType
	strContent []string
	content    []interface{}
}

func (p Property) Type() EntryType {
	return p.entryType
}

func (p Property) Content() []string {
	if len(p.content) > len(p.strContent) {
		ss := make([]string, len(p.content))
		for i, s := range p.content {
			ss[i] = s.(string)
		}
		p.strContent = ss
	}

	return p.strContent
}

func unmarshalMicroformat(r interface{}, o *MicroformatObject) error {
	m := r.(map[string]interface{})
	for k, v := range m {
		switch k {
		case "type":
			o.Type = FromString(
				v.([]interface{})[0].(string),
			)
			break
		case "properties":
			if o.Properties == nil {
				o.Properties = make(map[string]NestedProperty)
			}
			props := v.(map[string]interface{})
			for k, v := range props {
				o.Properties[k] = &Property{
					content: v.([]interface{}),
				}
			}
			break
		}
	}
	return nil
}

func (p Property) Object() *MicroformatObject {
	result := new(MicroformatObject)
	if len(p.content) > 0 {
		unmarshalMicroformat(p.content[0], result)
	}
	return result
}

type MicroformatObject struct {
	Type       EntryType                 `json:"type"`
	Properties map[string]NestedProperty `json:"properties"`
}

func (o *MicroformatObject) UnmarshalJSON(data []byte) error {
	var r interface{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}

	return unmarshalMicroformat(r, o)
}

func ParseForm(f *url.Values, o *MicroformatObject) error {
	if o.Properties == nil {
		o.Properties = make(map[string]NestedProperty)
	}

	for k, values := range *f {
		k = strings.TrimSuffix(k, "[]")
		if k == "h" {
			o.Type = FromString(values[0])
		} else {
			o.Properties[k] = Property{
				strContent: values,
			}
		}
	}
	return nil
}
