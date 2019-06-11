package tq

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	data, err := ioutil.ReadFile("test/table.toml")
	require.Nil(t, err)

	cases := []struct {
		input  string
		output string
	}{
		{input: "", output: string(data)},
		{input: "table.key", output: "value"},
		{input: "table.subtable", output: `key = "another value"`},
		{input: "table.inline.name.first", output: "Tom"},
		{input: "table.inline.point.y", output: "2"},
		{
			input: "table",
			output: `key = "value"

[inline]
  [inline.name]
    first = "Tom"
    last = "Preston-Werner"
  [inline.point]
    x = 1
    y = 2

[subtable]
  key = "another value"`,
		},
	}

	for _, c := range cases {
		ret, err := Find(string(data), c.input)
		require.Nil(t, err)
		require.Equal(t, c.output, ret)
	}
}

func TestErrorToml(t *testing.T) {
	data, err := ioutil.ReadFile("test/error.toml")
	require.Nil(t, err)

	_, err = Find(string(data), "fish")
	require.NotNil(t, err)
	require.Equal(t, "toml decode: Near line 0 (last key parsed ''): bare keys cannot contain '\\x00'", err.Error())
}

func TestKVToml(t *testing.T) {
	data, err := ioutil.ReadFile("test/kv.toml")
	require.Nil(t, err)

	ret, err := Find(string(data), "name")
	require.Nil(t, err)
	require.Equal(t, "Tom Preston-Werner", ret)
}

func TestFuck(t *testing.T) {
	conf := `data = ["xcc", "yii"]`
	m := make(map[string]interface{})
	md, err := toml.Decode(conf, &m)
	require.Nil(t, err)
	fmt.Println(md)
	fmt.Println(md.Keys())
}

func TestSplitSelector(t *testing.T) {
	cases := []struct {
		selector string
		first    string
		other    string
	}{
		{selector: ".", first: "", other: ""},
		{selector: "..", first: "", other: "."},
		{selector: "", first: "", other: ""},
		{selector: "tq", first: "tq", other: ""},
		{selector: "fruits.apple", first: "fruits", other: "apple"},
		{selector: "a.b.c", first: "a", other: "b.c"},
	}

	for _, c := range cases {
		first, other, err := splitSelector(c.selector)
		require.Nil(t, err)
		require.Equal(t, c.first, first)
		require.Equal(t, c.other, other)
	}
}

func TestSplitSelectorByEnd(t *testing.T) {
	cases := []struct {
		selector string
		first    string
		other    string
	}{
		// {selector: ".", first: "", other: ""},
		// {selector: "..", first: "", other: "."},
		// {selector: "", first: "", other: ""},
		// {selector: "tq", first: "tq", other: ""},
		{selector: "fruits.apple", first: "fruits", other: "apple"},
		{selector: "a.b.c", first: "a.b", other: "c"},
	}

	for _, c := range cases {
		first, other, err := splitSelectorByEnd(c.selector)
		require.Nil(t, err)
		require.Equal(t, c.first, first)
		require.Equal(t, c.other, other)
	}
}

func TestIsArraySelector(t *testing.T) {
	cases := []struct {
		input string
		ret   bool
	}{
		{input: "[]", ret: true},
		{input: "[", ret: false},
		{input: "]", ret: false},
	}

	for _, c := range cases {
		r := isArraySelector(c.input)
		require.Equal(t, c.ret, r)
	}
}
