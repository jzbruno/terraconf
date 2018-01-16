package terraconf

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/terraform"
)

const (
	tfResourceNameDelimiter = "."
	tfDataSourcePrefix      = "data" + tfResourceNameDelimiter
)

func formatConfig(s string) (string, error) {
	b, err := printer.Format([]byte(s))
	if err != nil {
		return "", fmt.Errorf("error formatting terraform config string: %s", err)
	}

	return string(b), nil
}

func getAttrName(name string) string {
	return strings.SplitN(name, tfResourceNameDelimiter, 2)[0]
}

func getAttributeNamesSorted(attrMap map[string]string) []string {
	nameMap := map[string]struct{}{}

	for k := range attrMap {
		nameMap[getAttrName(k)] = struct{}{}
	}

	names := make([]string, 0, len(nameMap))

	for n := range nameMap {
		names = append(names, n)
	}

	sort.Strings(names)

	return names
}

func tfStringValue(i interface{}) string {
	return strconv.Quote(fmt.Sprintf("%v", i))
}

func tfAttributeValue(name, i interface{}) string {
	rawType := reflect.TypeOf(i)

	s := ""

	switch rawType.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if name == "" {
			s += fmt.Sprintf("%s,\n", tfStringValue(i))
		} else {
			s += fmt.Sprintf("%s = %s\n", name, tfStringValue(i))
		}
	case reflect.Map:
		rawMap := i.(map[string]interface{})

		if len(rawMap) == 0 {
			break
		}

		// Sort map by key to create consistent config between terraconf runs
		keys := make([]string, 0, len(rawMap))
		for rawMapKey := range rawMap {
			keys = append(keys, rawMapKey)
		}
		sort.Strings(keys)

		s += fmt.Sprintf("%s {\n", name)
		for _, key := range keys {
			s += tfAttributeValue(key, rawMap[key])
		}
		s += "}\n"
	case reflect.Slice:
		rawSlice := i.([]interface{})

		if len(rawSlice) == 0 {
			break
		}

		rawSliceType := reflect.TypeOf(rawSlice[0])

		switch rawSliceType.Kind() {
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			s += fmt.Sprintf("%s = [\n", name)
			for _, rawSliceValue := range rawSlice {
				s += tfAttributeValue("", rawSliceValue)
			}
			s += "]\n"
		case reflect.Map:
			for _, rawSliceValue := range rawSlice {
				s += tfAttributeValue(name, rawSliceValue)
			}
		}
	}

	return s
}

func GetStateConfigString(state *terraform.State) (string, error) {
	s := ""

	for _, module := range state.Modules {
		for resourceID, resource := range module.Resources {
			// Skip data source types. No way to tell which attributes are input vs. output.
			if strings.HasPrefix(resourceID, tfDataSourcePrefix) {
				continue
			}

			resourceNameParts := strings.Split(resourceID, tfResourceNameDelimiter)

			// Skip invalid resource names.
			if len(resourceNameParts) < 2 {
				continue
			}

			resourceName := resourceNameParts[1]
			attrs := resource.Primary.Attributes
			attrNames := getAttributeNamesSorted(attrs)

			s += fmt.Sprintf("resource \"%s\" \"%s\" {\n", resource.Type, resourceName)

			for _, attrName := range attrNames {
				rawAttr := flatmap.Expand(attrs, attrName)
				s += tfAttributeValue(attrName, rawAttr)
			}

			if len(resource.Dependencies) > 0 {
				s += "depends_on = [\n"
				for _, dependency := range resource.Dependencies {
					s += fmt.Sprintf("%s,\n", tfStringValue(dependency))
				}
				s += "]\n"
			}

			s += "}\n"
		}
	}

	s, err := formatConfig(s)
	if err != nil {
		return s, fmt.Errorf("error formatting config string for state: %s", err)
	}

	return s, nil
}
