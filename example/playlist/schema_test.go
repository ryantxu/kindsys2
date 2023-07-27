package playlist_test

import (
	"bytes"
	"kindsys2/example/playlist/schema"
	"kindsys2/jsks"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlaylistKind(t *testing.T) {
	k, err := jsks.NewResourceKind(schema.ManifestFS)
	require.NoError(t, err)

	// out, err := jsks.HackGetManifestJSON(k)
	// require.NoError(t, err)
	// fmt.Printf("GOT: %s\n", string(out))

	raw, err := os.ReadFile("testdata/valid-v0-0.json")
	require.NoError(t, err)

	obj, err := k.Parse(bytes.NewReader(raw))
	require.NoError(t, err)

	err = k.Validate(obj)
	require.NoError(t, err)

	// INVALID

	raw, err = os.ReadFile("testdata/invalid-v0-0.json")
	require.NoError(t, err)

	obj, err = k.Parse(bytes.NewReader(raw))
	require.NoError(t, err)

	err = k.Validate(obj)
	require.NoError(t, err)
}
