package playlist_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	obj, err := k.Read(bytes.NewReader(raw), true)
	require.NoError(t, err)

	// Expect the same value after writing it out
	out, err := json.MarshalIndent(obj, "", "  ")
	require.NoError(t, err)
	//require.JSONEq(t, string(raw), string(out))

	fmt.Printf("AFTER: %s\n", string(out))

	// INVALID

	// raw, err = os.ReadFile("testdata/invalid-v0-0.json")
	// require.NoError(t, err)

	// obj, err = k.Read(bytes.NewReader(raw), true)
	// require.NoError(t, err)
}
