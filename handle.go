package tq

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/BurntSushi/toml"
)

func Handle(m map[string]interface{}, md toml.MetaData, selector, value string) (string, error) {
	// check selector and value
	if selector == "" {
		printMeta(md)
		return "", nil
	}

	if value == "" {
		return fetch(m, selector)
	}

	return replace(m, selector, value)
}

func replace(m map[string]interface{}, selector, value string) (string, error) {
	first, end, err := splitSelectorByEnd(selector)
	if err != nil {
		return "", err
	}

	ret, err := parse(m, first)
	if err != nil {
		return "", err
	}

	mm := (ret).(map[string]interface{})
	switch mm[end].(type) {
	case bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return "", err
		}
		mm[end] = b
	case []interface{}:
		value = strings.TrimLeft(value, "[")
		value = strings.TrimRight(value, "]")
		arr := strings.Split(value, ",")
		res := make([]interface{}, 0)
		for _, a := range arr {
			a = strings.Trim(a, " ")
			a = strings.Trim(a, "\"")
			res = append(res, a)
		}
		mm[end] = res
	default:
		mm[end] = value
	}

	return encode(m)
}

// Parse parses input with regex
func fetch(m map[string]interface{}, selector string) (string, error) {
	ret, err := parse(m, selector)
	if err != nil {
		return "", err
	}

	return encode(ret)
}

func parse(input map[string]interface{}, selector string) (interface{}, error) {
	m := interface{}(input)

	for {
		v, s, err := find(m, selector)
		if err != nil {
			return nil, err
		}

		if s == "" {
			return v, nil
		}

		m = v
		selector = s
	}
}

func find(m interface{}, selector string) (interface{}, string, error) {
	first, other, err := splitSelector(selector)
	if err != nil {
		return nil, "", err
	}

	if isArraySelector(first) {

	}

	mm := m.(map[string]interface{})
	v, ok := mm[first]
	if !ok {
		return nil, "", fmt.Errorf("can't find key %s", selector)
	}

	return v, other, nil
}

func encode(ret interface{}) (string, error) {
	switch r := ret.(type) {
	case map[string]interface{}:
		w := bytes.NewBuffer(nil)
		f := toml.NewEncoder(w)
		if err := f.Encode(ret); err != nil {
			return "", err
		}

		return w.String(), nil
	case []map[string]interface{}:
		w := bytes.NewBuffer(nil)
		f := toml.NewEncoder(w)
		for _, v := range r {
			if err := f.Encode(v); err != nil {
				return "", err
			}
		}

		return w.String(), nil
	default:
		return fmt.Sprintf("%v", ret), nil
	}

}

func printMeta(meta toml.MetaData) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, key := range meta.Keys() {
		fmt.Fprintf(tw, "%s%s\t%s\n",
			strings.Repeat("    ", len(key)-1), key, meta.Type(key...))
	}
	tw.Flush()
}

func splitSelector(s string) (string, string, error) {
	if s == "" {
		return "", "", nil
	}

	arr := strings.Split(s, ".")

	return arr[0], strings.Join(arr[1:], "."), nil
}

func splitSelectorByEnd(s string) (string, string, error) {
	if s == "" {
		return "", "", nil
	}

	arr := strings.Split(s, ".")
	return strings.Join(arr[:len(arr)-1], "."), arr[len(arr)-1], nil
}

func isArraySelector(s string) bool {
	if s == "" {
		return false
	}

	if s[0] == '[' && s[len(s)-1] == ']' {
		return true
	}

	return false
}
