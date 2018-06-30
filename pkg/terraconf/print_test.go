package terraconf

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

// Tests were initially generated with https://github.com/cweill/gotests

func Test_getAttributeNamesSorted(t *testing.T) {
	type args struct {
		attrMap map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"should return sorted map keys",
			args{map[string]string{"b": "", "a": "", "c": ""}},
			[]string{"a", "b", "c"}},
		{
			"should return single item map",
			args{map[string]string{"a": ""}},
			[]string{"a"},
		},
		{
			"should return empty slice for empty map",
			args{map[string]string{}},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAttributeNamesSorted(tt.args.attrMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAttributeNamesSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatConfig(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"should format valid hcl",
			args{"resource \"test\" \"test\"{    }"},
			"resource \"test\" \"test\" {}\n",
			false,
		},
		{
			"should return error for invalid hcl",
			args{"resource \"test\" \"test\"{"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatConfig(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("formatConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAttrName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"should return name for input with no delimiters",
			args{"test"},
			"test",
		},
		{
			"should return name for input with delimiters",
			args{"test.one"},
			"test",
		},
		{
			"should return name for input with multiple delimiters",
			args{"test.one.two"},
			"test",
		},
		{
			"should return name for input with multiple delimiters",
			args{".test"},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAttrName(tt.args.name); got != tt.want {
				t.Errorf("getAttrName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tfStringValue(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"should return quoted terraform string for input string",
			args{"test"},
			"\"test\"",
		},
		{
			"should return quoted terraform string for input bool",
			args{true},
			"\"true\"",
		},
		{
			"should return quoted terraform string for input int",
			args{42},
			"\"42\"",
		},
		{
			"should return quoted terraform string for input json",
			args{`{
	"one": 1,
	"two": "two",
}`},
			strconv.Quote("{\n\t\"one\": 1,\n\t\"two\": \"two\",\n}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tfStringValue(tt.args.i); got != tt.want {
				t.Errorf("tfStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tfAttributeValue(t *testing.T) {
	type args struct {
		name interface{}
		i    interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"should return tf string for simple string input",
			args{"test", "test"},
			"test = \"test\"\n",
		},
		{
			"should return tf string for simple bool input",
			args{"test", true},
			"test = \"true\"\n",
		},
		{
			"should return tf string for simple int input",
			args{"test", 42},
			"test = \"42\"\n",
		},
		{
			"should return tf string for simple string interface slice input",
			args{"test", []interface{}{"one", "two", "three"}},
			"test = [\n\"one\",\n\"two\",\n\"three\",\n]\n",
		},
		{
			"should return empty string for empty slice",
			args{"test", []interface{}{}},
			"",
		},
		{
			"should return tf string for simple bool interface slice input",
			args{"test", []interface{}{true, false, true}},
			"test = [\n\"true\",\n\"false\",\n\"true\",\n]\n",
		},
		{
			"should return tf string for simple int interface slice input",
			args{"test", []interface{}{2, 3, 1}},
			"test = [\n\"2\",\n\"3\",\n\"1\",\n]\n",
		},
		{
			"should return empty tf string for unsupported slice of slice",
			args{"test", []interface{}{[]interface{}{1, 2, 3}, []interface{}{1, 2, 3}}},
			"",
		},
		{
			"should return tf string for slice of maps",
			args{"test", []interface{}{
				map[string]interface{}{"one": 1, "two": 2, "three": 3},
				map[string]interface{}{"one": 1, "two": 2, "three": 3},
			}},
			"test {\none = \"1\"\nthree = \"3\"\ntwo = \"2\"\n}\ntest {\none = \"1\"\nthree = \"3\"\ntwo = \"2\"\n}\n",
		},
		{
			"should return empty string for empty map",
			args{"test", map[string]interface{}{}},
			"",
		},
		{
			"should return tf string for simple map",
			args{"test", map[string]interface{}{"one": 1, "two": 2, "three": 3}},
			"test {\none = \"1\"\nthree = \"3\"\ntwo = \"2\"\n}\n",
		},
		{
			"should return tf string for complex map",
			args{"test", map[string]interface{}{
				"one":   1,
				"two":   2,
				"three": []interface{}{1, 2, 3},
				"four":  map[string]interface{}{"one": 1, "two": 2},
			}},
			"test {\nfour {\none = \"1\"\ntwo = \"2\"\n}\none = \"1\"\nthree = [\n\"1\",\n\"2\",\n\"3\",\n]\ntwo = \"2\"\n}\n",
		},
		{
			"should return tf string with quoted key for map with key of *",
			args{"test", map[string]interface{}{
				"*": 0,
			}},
			"test {\n\"*\" = \"0\"\n}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tfAttributeValue(tt.args.name, tt.args.i); got != tt.want {
				t.Errorf("tfAttributeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

var simpleState = &terraform.State{
	Modules: []*terraform.ModuleState{
		{
			Resources: map[string]*terraform.ResourceState{
				"test_type.test_id": {
					Type: "test_type",
					Primary: &terraform.InstanceState{
						ID: "test_id",
						Attributes: map[string]string{
							"one": "1",
						},
					},
				},
			},
		},
	},
}
var simpleStateWithInvalidResourceName = &terraform.State{
	Modules: []*terraform.ModuleState{
		{
			Resources: map[string]*terraform.ResourceState{
				"test_invalid_name": {},
			},
		},
	},
}
var simpleStateWithDataSource = &terraform.State{
	Modules: []*terraform.ModuleState{
		{
			Resources: map[string]*terraform.ResourceState{
				"data.test_type.test_id": {
					Type: "test_type",
					Primary: &terraform.InstanceState{
						ID: "test_id",
						Attributes: map[string]string{
							"one": "1",
						},
					},
				},
			},
		},
	},
}
var simpleStateWithDependencies = &terraform.State{
	Modules: []*terraform.ModuleState{
		{
			Resources: map[string]*terraform.ResourceState{
				"test_type.test_id": {
					Type: "test_type",
					Primary: &terraform.InstanceState{
						ID: "test_id",
						Attributes: map[string]string{
							"one": "1",
						},
					},
					Dependencies: []string{
						"test_dep_one",
						"test_dep_two",
					},
				},
			},
		},
	},
}

func TestGetStateConfigString(t *testing.T) {
	type args struct {
		state *terraform.State
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"should not generate tf config for invalid resource names",
			args{simpleStateWithInvalidResourceName},
			"\n",
			false,
		},
		{
			"should generate tf config string from simple state file input",
			args{simpleState},
			"resource \"test_type\" \"test_id\" {\n  one = \"1\"\n}\n",
			false,
		},
		{
			"should not generate tf config for data sources",
			args{simpleStateWithDataSource},
			"\n",
			false,
		},
		{
			"should generate tf config string with dependencies from simple state file input",
			args{simpleStateWithDependencies},
			"resource \"test_type\" \"test_id\" {\n  one = \"1\"\n\n  depends_on = [\n    \"test_dep_one\",\n    \"test_dep_two\",\n  ]\n}\n",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStateConfigString(tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStateConfigString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetStateConfigString() = %v, want %v", got, tt.want)
			}
		})
	}
}
