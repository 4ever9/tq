package tq

import (
	"testing"

	"github.com/pelletier/go-toml"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	tree, err := toml.LoadFile("test/table.toml")
	require.Nil(t, err)

	cases := []struct {
		key    string
		output string
	}{
		{key: "", output: tree.String()},
		{key: "table.key", output: "value"},
		{key: "table.subtable", output: "key = \"another value\"\n"},
		{key: "table.inline.name.first", output: "Tom"},
		{key: "table.inline.point.y", output: "2"},
	}

	for _, c := range cases {
		ret, err := Handle(tree, c.key, "")
		require.Nil(t, err)
		require.Equal(t, c.output, ret)
	}
}
