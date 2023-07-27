package playlist_test

import (
	"encoding/json"
	"fmt"
	"kindsys2/example/playlist"
	"testing"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestModelGeneration(t *testing.T) {
	r := new(jsonschema.Reflector)
	if err := r.AddGoComments("kindsys2/example/playlist", "./"); err != nil {
		t.Fatal(err)
	}
	s := r.Reflect(&playlist.Spec{})
	data, err := json.MarshalIndent(s, "", "  ")
	require.NoError(t, err)

	fmt.Println(string(data))

	fmt.Printf("hello: %d", time.Now().UnixMilli())

	//t.FailNow()
}
