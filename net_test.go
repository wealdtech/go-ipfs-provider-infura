package infura

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBadPath(t *testing.T) {
	provider, err := NewProvider()
	require.Nil(t, err, "unexpected error")

	res, err := provider.get("https://ipfs.infura.io:5001/api/v0/pin/bad", "")
	require.Nil(t, err)
	require.Equal(t, "ipfs method not allowed", res["message"])
}
