package thema

import (
	"bytes"
	"testing"
	"time"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
	"github.com/stretchr/testify/require"
)

func TestThemaResource(t *testing.T) {
	var testkind = `
name: "TestKind"
description: "Blammo!"
maturity: "experimental"
lineage: schemas: [{
	version: [0, 0]
	schema: {
		spec: aSpecField: int32
	}
}]
`

	var testresource = `
{
	"apiVersion": "core.grafana.com/v0",
	"kind": "TestKind",
	"metadata": {
		"name": "test",
		"namespace": "default",
		"annotations": {
			"grafana.com/createdBy": "me",
			"grafana.com/updatedBy": "you",
			"grafana.com/updateTimestamp": "2023-07-06T03:08:01Z"
		}
	},
	"spec": {
		"aSpecField": 42
	}
}`

	ctx := cuecontext.New()

	rt := thema.NewRuntime(ctx)
	cv := ctx.CompileString(testkind)
	def, err := kindsys.ToDef[kindsys.CoreProperties](cv)
	require.NoError(t, err)

	k, err := NewThemaResourceKind(rt, def)
	require.NoError(t, err)

	res, err := k.Read(bytes.NewReader([]byte(testresource)), true)
	require.NoError(t, err)

	require.Equal(t, "me", res.CommonMetadata().CreatedBy)
	require.Equal(t, "you", res.CommonMetadata().UpdatedBy)
	require.Equal(t, "2023-07-06T03:08:01Z", res.CommonMetadata().UpdateTimestamp.Format(time.RFC3339))
}
