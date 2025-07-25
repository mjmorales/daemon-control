package plist

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDict_AddString(t *testing.T) {
	dict := &Dict{}
	dict.AddString("TestKey", "TestValue")

	assert.Len(t, dict.Items, 2)
	key, ok := dict.Items[0].(Key)
	assert.True(t, ok)
	assert.Equal(t, "TestKey", key.Value)

	str, ok := dict.Items[1].(String)
	assert.True(t, ok)
	assert.Equal(t, "TestValue", str.Value)
}

func TestDict_AddInteger(t *testing.T) {
	dict := &Dict{}
	dict.AddInteger("Count", 42)

	assert.Len(t, dict.Items, 2)
	key, ok := dict.Items[0].(Key)
	assert.True(t, ok)
	assert.Equal(t, "Count", key.Value)

	integer, ok := dict.Items[1].(Integer)
	assert.True(t, ok)
	assert.Equal(t, 42, integer.Value)
}

func TestDict_AddBool(t *testing.T) {
	tests := []struct {
		name  string
		value bool
		want  interface{}
	}{
		{
			name:  "true value",
			value: true,
			want:  True{},
		},
		{
			name:  "false value",
			value: false,
			want:  False{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := &Dict{}
			dict.AddBool("BoolKey", tt.value)

			assert.Len(t, dict.Items, 2)
			key, ok := dict.Items[0].(Key)
			assert.True(t, ok)
			assert.Equal(t, "BoolKey", key.Value)

			// Check type matches expected
			assert.IsType(t, tt.want, dict.Items[1])
		})
	}
}

func TestDict_AddDict(t *testing.T) {
	parent := &Dict{}
	child := &Dict{}
	child.AddString("ChildKey", "ChildValue")

	parent.AddDict("SubDict", child)

	assert.Len(t, parent.Items, 2)
	key, ok := parent.Items[0].(Key)
	assert.True(t, ok)
	assert.Equal(t, "SubDict", key.Value)

	subDict, ok := parent.Items[1].(*Dict)
	assert.True(t, ok)
	assert.Len(t, subDict.Items, 2)
}

func TestDict_AddStringArray(t *testing.T) {
	dict := &Dict{}
	values := []string{"value1", "value2", "value3"}
	dict.AddStringArray("StringArray", values)

	assert.Len(t, dict.Items, 2)
	key, ok := dict.Items[0].(Key)
	assert.True(t, ok)
	assert.Equal(t, "StringArray", key.Value)

	array, ok := dict.Items[1].(*Array)
	assert.True(t, ok)
	assert.Len(t, array.Items, 3)

	for i, item := range array.Items {
		str, ok := item.(String)
		assert.True(t, ok)
		assert.Equal(t, values[i], str.Value)
	}
}

func TestDict_AddDictArray(t *testing.T) {
	parent := &Dict{}

	dict1 := &Dict{}
	dict1.AddString("Key1", "Value1")

	dict2 := &Dict{}
	dict2.AddString("Key2", "Value2")

	dicts := []*Dict{dict1, dict2}
	parent.AddDictArray("DictArray", dicts)

	assert.Len(t, parent.Items, 2)
	key, ok := parent.Items[0].(Key)
	assert.True(t, ok)
	assert.Equal(t, "DictArray", key.Value)

	array, ok := parent.Items[1].(*Array)
	assert.True(t, ok)
	assert.Len(t, array.Items, 2)
}

func TestDict_MarshalXML(t *testing.T) {
	dict := &Dict{}
	dict.AddString("Name", "TestDaemon")
	dict.AddInteger("Port", 8080)
	dict.AddBool("Enabled", true)

	data, err := xml.Marshal(dict)
	require.NoError(t, err)

	xmlStr := string(data)
	assert.Contains(t, xmlStr, "<dict>")
	assert.Contains(t, xmlStr, "</dict>")
	assert.Contains(t, xmlStr, "<key>Name</key>")
	assert.Contains(t, xmlStr, "<string>TestDaemon</string>")
	assert.Contains(t, xmlStr, "<key>Port</key>")
	assert.Contains(t, xmlStr, "<integer>8080</integer>")
	assert.Contains(t, xmlStr, "<key>Enabled</key>")
	assert.Contains(t, xmlStr, "<true></true>")
}

func TestArray_MarshalXML(t *testing.T) {
	array := &Array{
		Items: []interface{}{
			String{Value: "item1"},
			String{Value: "item2"},
			Integer{Value: 42},
		},
	}

	data, err := xml.Marshal(array)
	require.NoError(t, err)

	xmlStr := string(data)
	assert.Contains(t, xmlStr, "<array>")
	assert.Contains(t, xmlStr, "</array>")
	assert.Contains(t, xmlStr, "<string>item1</string>")
	assert.Contains(t, xmlStr, "<string>item2</string>")
	assert.Contains(t, xmlStr, "<integer>42</integer>")
}

func TestPlist_Validate(t *testing.T) {
	tests := []struct {
		name      string
		plist     *Plist
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid plist with label",
			plist: &Plist{
				Version: "1.0",
				Dict: &Dict{
					Items: []interface{}{
						Key{Value: "Label"},
						String{Value: "com.example.test"},
					},
				},
			},
			wantError: false,
		},
		{
			name: "missing dict",
			plist: &Plist{
				Version: "1.0",
				Dict:    nil,
			},
			wantError: true,
			errorMsg:  "plist dict is nil",
		},
		{
			name: "missing label",
			plist: &Plist{
				Version: "1.0",
				Dict: &Dict{
					Items: []interface{}{
						Key{Value: "Program"},
						String{Value: "/usr/bin/test"},
					},
				},
			},
			wantError: true,
			errorMsg:  "plist missing required Label key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.plist.Validate()

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlist_MarshalXML(t *testing.T) {
	dict := &Dict{}
	dict.AddString("Label", "com.example.test")
	dict.AddString("Program", "/usr/bin/test")
	dict.AddBool("RunAtLoad", true)

	plist := &Plist{
		Version: "1.0",
		Dict:    dict,
	}

	data, err := xml.Marshal(plist)
	require.NoError(t, err)

	xmlStr := string(data)
	assert.Contains(t, xmlStr, `<plist version="1.0">`)
	assert.Contains(t, xmlStr, "</plist>")
	assert.Contains(t, xmlStr, "<dict>")
	assert.Contains(t, xmlStr, "<key>Label</key>")
}

func TestNestedStructures(t *testing.T) {
	// Test complex nested structure
	root := &Dict{}

	// Add environment variables dict
	envDict := &Dict{}
	envDict.AddString("PATH", "/usr/local/bin:/usr/bin")
	envDict.AddString("HOME", "/Users/test")
	root.AddDict("EnvironmentVariables", envDict)

	// Add array of dicts
	intervals := []*Dict{
		{Items: []interface{}{Key{Value: "Hour"}, Integer{Value: 9}}},
		{Items: []interface{}{Key{Value: "Hour"}, Integer{Value: 17}}},
	}
	root.AddDictArray("StartCalendarInterval", intervals)

	// Add string array
	root.AddStringArray("ProgramArguments", []string{"/usr/bin/python", "script.py"})

	// Marshal and check
	data, err := xml.Marshal(root)
	require.NoError(t, err)

	xmlStr := string(data)

	// Verify nested structures
	assert.Contains(t, xmlStr, "<key>EnvironmentVariables</key>")
	assert.Contains(t, xmlStr, "<dict><key>PATH</key>")
	assert.Contains(t, xmlStr, "<key>StartCalendarInterval</key>")
	assert.Contains(t, xmlStr, "<array><dict>")

	// Count occurrences
	assert.Equal(t, 2, strings.Count(xmlStr, "<key>Hour</key>"))
}

func TestXMLIndentation(t *testing.T) {
	dict := &Dict{}
	dict.AddString("Label", "com.example.test")

	plist := &Plist{
		Version: "1.0",
		Dict:    dict,
	}

	// Marshal with indentation
	var buf strings.Builder
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "    ")

	err := encoder.Encode(plist)
	require.NoError(t, err)

	// Check indentation is present
	xmlStr := buf.String()
	assert.Contains(t, xmlStr, "\n    <dict>")
	assert.Contains(t, xmlStr, "\n        <key>")
}

func TestEmptyStructures(t *testing.T) {
	tests := []struct {
		name     string
		items    []interface{}
		expected string
	}{
		{
			name:     "empty dict",
			items:    []interface{}{&Dict{}},
			expected: "<dict></dict>",
		},
		{
			name:     "empty array",
			items:    []interface{}{&Array{}},
			expected: "<array></array>",
		},
		{
			name:     "empty string",
			items:    []interface{}{String{Value: ""}},
			expected: "<string></string>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(tt.items[0])
			require.NoError(t, err)
			assert.Contains(t, string(data), tt.expected)
		})
	}
}
