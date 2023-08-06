package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func PrintData(label string, v interface{}) {
	fmt.Println(SPrintData(label, v))
}
func SPrintData(label string, v interface{}) string {
	if r, ok := v.(io.Reader); ok {
		bytes, err := io.ReadAll(r)
		panicOnErr(label, err)
		if len(bytes) == 0 {
			return ""
		}
		var v2 interface{}
		panicOnErr(label, json.Unmarshal(bytes, &v2))
		v = v2
	}
	b, err := json.MarshalIndent(v, "", "  ")
	panicOnErr(label, err)
	return label + ": " + string(b)
}
func panicOnErr(label string, err error) {
	if err != nil {
		panic(label + ": " + err.Error())
	}
}

//lint:ignore U1000 keep it for dev testing
func TLogJson(t testing.TB, label string, v any) {
	t.Helper()
	t.Logf("%s: %s", label, JSONString(t, v))
}
func JSONString(t testing.TB, v any) string {
	t.Helper()
	if r, ok := v.(io.Reader); ok {
		b, err := io.ReadAll(r)
		require.NoError(t, err)
		// require.NotEmpty(t, b)
		return string(b)
	}
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	return string(b)
}
