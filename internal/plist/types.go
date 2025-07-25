package plist

import (
	"encoding/xml"
	"fmt"
)

// Plist represents the root plist structure
type Plist struct {
	XMLName xml.Name `xml:"plist"`
	Version string   `xml:"version,attr"`
	Dict    *Dict    `xml:"dict"`
}

// Dict represents a dictionary in plist
type Dict struct {
	Items []interface{}
}

// MarshalXML custom marshaler for Dict
func (d *Dict) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "dict"}}); err != nil {
		return err
	}
	
	for _, item := range d.Items {
		if err := e.Encode(item); err != nil {
			return err
		}
	}
	
	return e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "dict"}})
}

// Key represents a key element
type Key struct {
	XMLName xml.Name `xml:"key"`
	Value   string   `xml:",chardata"`
}

// String represents a string element
type String struct {
	XMLName xml.Name `xml:"string"`
	Value   string   `xml:",chardata"`
}

// Integer represents an integer element
type Integer struct {
	XMLName xml.Name `xml:"integer"`
	Value   int      `xml:",chardata"`
}

// True represents a boolean true element
type True struct {
	XMLName xml.Name `xml:"true"`
}

// False represents a boolean false element
type False struct {
	XMLName xml.Name `xml:"false"`
}

// Array represents an array element
type Array struct {
	XMLName xml.Name `xml:"array"`
	Items   []interface{}
}

// AddString adds a string key-value pair to the dictionary
func (d *Dict) AddString(key, value string) {
	d.Items = append(d.Items, Key{Value: key}, String{Value: value})
}

// AddInteger adds an integer key-value pair to the dictionary
func (d *Dict) AddInteger(key string, value int) {
	d.Items = append(d.Items, Key{Value: key}, Integer{Value: value})
}

// AddBool adds a boolean key-value pair to the dictionary
func (d *Dict) AddBool(key string, value bool) {
	d.Items = append(d.Items, Key{Value: key})
	if value {
		d.Items = append(d.Items, True{})
	} else {
		d.Items = append(d.Items, False{})
	}
}

// AddDict adds a dictionary key-value pair to the dictionary
func (d *Dict) AddDict(key string, dict *Dict) {
	d.Items = append(d.Items, Key{Value: key}, dict)
}

// AddStringArray adds a string array key-value pair to the dictionary
func (d *Dict) AddStringArray(key string, values []string) {
	array := &Array{}
	for _, v := range values {
		array.Items = append(array.Items, String{Value: v})
	}
	d.Items = append(d.Items, Key{Value: key}, array)
}

// AddDictArray adds a dictionary array key-value pair to the dictionary
func (d *Dict) AddDictArray(key string, dicts []*Dict) {
	array := &Array{}
	for _, dict := range dicts {
		array.Items = append(array.Items, dict)
	}
	d.Items = append(d.Items, Key{Value: key}, array)
}

// MarshalXML custom marshaler for Array
func (a *Array) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "array"}}); err != nil {
		return err
	}
	
	for _, item := range a.Items {
		if err := e.Encode(item); err != nil {
			return err
		}
	}
	
	return e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "array"}})
}

// Validate checks if the plist structure is valid
func (p *Plist) Validate() error {
	if p.Dict == nil {
		return fmt.Errorf("plist dict is nil")
	}
	
	// Check for required Label key
	hasLabel := false
	for i := 0; i < len(p.Dict.Items); i += 2 {
		if key, ok := p.Dict.Items[i].(Key); ok && key.Value == "Label" {
			hasLabel = true
			break
		}
	}
	
	if !hasLabel {
		return fmt.Errorf("plist missing required Label key")
	}
	
	return nil
}