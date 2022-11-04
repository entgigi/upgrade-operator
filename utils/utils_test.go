package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyMap(t *testing.T) {
	m1 := map[string]interface{}{
		"a": "bbb",
		"b": map[string]interface{}{
			"c": 123,
		},
		"c": map[string]interface{}{
			"c": 123,
		},
	}

	m2 := CopyMap(m1)

	m1["a"] = "zzz"
	delete(m1, "b")
	m1["c"].(map[string]interface{})["c"] = 99

	require.Equal(t, map[string]interface{}{
		"a": "zzz",
		"c": map[string]interface{}{"c": 99},
	}, m1)
	require.Equal(t, map[string]interface{}{
		"a": "bbb",
		"b": map[string]interface{}{
			"c": 123,
		},
		"c": map[string]interface{}{"c": 123},
	}, m2)
}
