package dataframes_test

import (
	"bytes"
	"kindsys2/example/dataframes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDFKinds(t *testing.T) {
	k := dataframes.DataFramesResourceKind{}

	raw, err := os.ReadFile("testdata/valid-v1-0.json")
	require.NoError(t, err)

	obj, err := k.ReadFrames(bytes.NewReader(raw), true)
	require.NoError(t, err)

	require.Equal(t, 1, len(obj.Spec))
	require.Equal(t, 3, len(obj.Spec[0].Fields))
	require.Equal(t, 2, obj.Spec[0].Fields[0].Len())
}
