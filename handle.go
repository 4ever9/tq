package tq

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml"
)

func Handle(tree *toml.Tree, key, value string) (string, error) {
	// check selector and value
	if key == "" {
		return tree.ToTomlString()
	}

	if value == "" {
		v := fmt.Sprintf("%v", tree.Get(key))
		return v, nil
	}

	return replace(tree, key, value)
}

func replace(tree *toml.Tree, key string, value string) (string, error) {
	v := tree.Get(key)
	if v == nil {
		return tree.ToTomlString()
	}

	switch v.(type) {
	case string:
		tree.Set(key, value)
		return tree.ToTomlString()
	case int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return "", err
		}
		tree.Set(key, i)
		return tree.ToTomlString()
	case float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return "", err
		}
		tree.Set(key, f)
		return tree.ToTomlString()
	case []interface{}:
		arr := strings.Split(value, ",")
		switch v.([]interface{})[0].(type) {
		case string:
			tree.Set(key, arr)
			return tree.ToTomlString()
		case int64:
			ia := make([]int64, 0, len(v.([]interface{})))
			for _, a := range arr {
				i, err := strconv.ParseInt(a, 10, 64)
				if err != nil {
					return "", err
				}

				ia = append(ia, i)
			}

			tree.Set(key, ia)
			return tree.ToTomlString()
		case float64:
			fa := make([]float64, 0, len(v.([]interface{})))
			for _, a := range arr {
				f, err := strconv.ParseFloat(a, 64)
				if err != nil {
					return "", err
				}

				fa = append(fa, f)
			}

			tree.Set(key, fa)
			return tree.ToTomlString()
		}

		return tree.ToTomlString()
	}

	return tree.ToTomlString()
}
