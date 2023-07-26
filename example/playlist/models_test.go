package playlist_test

import (
	"encoding/json"
	"fmt"
	"kindsys2"
	"kindsys2/example/playlist"
	"testing"
	"time"

	"github.com/invopop/jsonschema"
)

func TestJSONSchema(t *testing.T) {
	// Playlist
	k := kindsys2.Manifest{
		KindInfo: kindsys2.KindInfo{
			Group:       "ext.playlists.grafana.com",
			Name:        "Playlist",
			Description: "Describes a set of dashboards that should be displayed in a loop",
		},
		MachineNames: &kindsys2.MachineNames{ // can be created manually
			Plural:   "playlists",
			Singular: "playlist",
		},
		Versions: []kindsys2.VersionInfo{
			{
				Version: "v0-0-alpha",
			}, {
				Version: "v0-1-alpha",
				Changelog: []string{
					"adding the dashboard_by_uid type",
					"deprecating the dashboard_by_id type",
					"deprecating the PlaylistItem.title property (now optional and unused)",
				},
			}, {
				Version: "v1-0-alpha",
				Changelog: []string{
					"removed the dashboard_by_id type",
					"removed the PlaylistItem.title property",
				},
			},
		},
	}

	data, err := json.MarshalIndent(k, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(data))

	r := new(jsonschema.Reflector)
	if err := r.AddGoComments("kindsys2/example/playlist", "./"); err != nil {
		t.Fatal(err)
	}
	s := r.Reflect(&playlist.Spec{})
	data, err = json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(data))

	fmt.Printf("hello: %d", time.Now().UnixMilli())

	t.FailNow()
}
