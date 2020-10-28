package drand

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/Secured-Finance/dione/config"
	drandClient "github.com/drand/drand/client/http"
	"github.com/stretchr/testify/assert"
)

func TestPrintGroupInfo(t *testing.T) {
	cfg := config.NewDrandConfig()
	ctx := context.Background()
	drandServer := cfg.Servers[0]
	client, err := drandClient.New(drandServer, nil, nil)
	assert.NoError(t, err)
	drandResult, err := client.Get(ctx, 266966)
	assert.NoError(t, err)
	stringSha256 := hex.EncodeToString(drandResult.Randomness())
	assert.Equal(t, stringSha256, "cb67e13477cad0e54540980a3b621dfd9c5fcd7c92ed42626289a1de6f25c3d1")
}
