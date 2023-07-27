package kindsys2

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnstructuredJSON(t *testing.T) {
	simple := &UnstructuredResource{}
	simple.StaticMeta.Group = "ext.something.grafana.com"
	simple.StaticMeta.Version = "v1-1"
	simple.StaticMeta.Kind = "Example"
	simple.StaticMeta.Name = "test"
	simple.StaticMeta.Namespace = "default"
	simple.CommonMeta.CreatedBy = "ryan"
	simple.CommonMeta.CreationTimestamp = time.Date(2020, time.January, 21, 1, 10, 30, 0, time.UTC)
	simple.CommonMeta.UpdateTimestamp = time.Date(2022, time.January, 21, 1, 10, 30, 0, time.UTC)
	simple.Spec = map[string]any{
		"hello":  "world",
		"number": 1.234,
		"int":    25,
	}

	out, err := json.MarshalIndent(simple, "", "  ")
	require.NoError(t, err)
	fmt.Printf("%s\n", string(out))
	require.JSONEq(t, `{
		"apiVersion": "ext.something.grafana.com/v1-1",
		"kind": "Example",
		"metadata": {
		  "name": "test",
		  "namespace": "default",
		  "annotations": {
			"grafana.com/createdBy": "ryan"
		  }
		},
		"spec": {
		  "hello": "world",
		  "int": 25,
		  "number": 1.234
		}
	  }`, string(out))

	copy := &UnstructuredResource{}
	json.Unmarshal(out, copy)
	require.NoError(t, err)

	after, err := json.Marshal(simple)
	require.NoError(t, err)
	require.JSONEq(t, string(out), string(after))
}
