package infura

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	provider "github.com/wealdtech/go-ipfs-provider"
)

const (
	testFileHash          = "QmeeLUVdiSTTKQqhWqsffYDtNvvvcTfJdotkNyi1KDEJtQ"
	testFileDirectoryHash = "QmP7RfPwpB8GgK5zhqQYxDRerUBxqVPkNruvyck864cRwy"
)

func TestPinContent(t *testing.T) {
	p, err := NewProvider()
	require.Nil(t, err, "unexpected error")

	file, err := os.Open("resources/testfile")
	require.Nil(t, err, "unexpected error")

	hash, err := p.PinContent("test file", file, nil)
	require.Nil(t, err, "unexpected error")

	assert.Equal(t, testFileHash, hash)
}

func TestPinContentOpts(t *testing.T) {
	p, err := NewProvider()
	require.Nil(t, err, "unexpected error")

	file, err := os.Open("resources/testfile")
	require.Nil(t, err, "unexpected error")

	hash, err := p.PinContent("testfile", file, &provider.ContentOpts{StoreInDirectory: true})
	require.Nil(t, err, "unexpected error")

	assert.Equal(t, testFileDirectoryHash, hash)
}

func TestItemStats(t *testing.T) {
	p, err := NewProvider()
	require.Nil(t, err, "unexpected error")

	item, err := p.ItemStats(testFileHash)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, testFileHash, item.Hash)
	assert.Equal(t, uint64(22), item.Size)
}

func TestItemStatsBadHash(t *testing.T) {
	p, err := NewProvider()
	require.Nil(t, err, "unexpected error")

	_, err = p.ItemStats("QmeeLUVdiSTTKQqhWqsffYDtNvvvcTfJdotkNyi1KD")
	require.NotNil(t, err, "missing expected error")
	require.Equal(t, "invalid path \"QmeeLUVdiSTTKQqhWqsffYDtNvvvcTfJdotkNyi1KD\": selected encoding not supported", err.Error())
}

func TestPin(t *testing.T) {
	p, err := NewProvider()
	require.Nil(t, err, "unexpected error")

	err = p.Pin(testFileHash)
	assert.Nil(t, err, "unexpected error")
}

func TestGatewayURL(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result string
		err    error
	}{
		{
			name:  "empty",
			input: "",
			err:   errors.New("unrecognised format"),
		},
		{
			name:  "bad",
			input: "bad",
			err:   errors.New("unrecognised format"),
		},
		{
			name:   "raw hash",
			input:  "QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD",
			result: "https://ipfs.infura.io/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD",
		},
		{
			name:  "raw hash with path",
			input: "QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD/index.html",
			err:   errors.New("unrecognised format"),
		},
		{
			name:   "IPFS multiaddr",
			input:  "/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD",
			result: "https://ipfs.infura.io/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD",
		},
		{
			name:   "IPFS multiaddr with path",
			input:  "/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD/index.html",
			result: "https://ipfs.infura.io/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD/index.html",
		},
		{
			name:   "IPFS URI",
			input:  "ipfs://QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD",
			result: "https://ipfs.infura.io/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD",
		},
		{
			name:   "IPFS URI with path",
			input:  "ipfs://QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD/index.html",
			result: "https://ipfs.infura.io/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD/index.html",
		},
		{
			name:   "IPNS URI",
			input:  "ipns://QmQ4QZh8nrsczdUEwTyfBope4THUhqxqc1fx6qYhhzZQei",
			result: "https://ipfs.infura.io/ipns/QmQ4QZh8nrsczdUEwTyfBope4THUhqxqc1fx6qYhhzZQei",
		},
		{
			name:   "IPNS URI with path",
			input:  "ipns://QmQ4QZh8nrsczdUEwTyfBope4THUhqxqc1fx6qYhhzZQei/index.html",
			result: "https://ipfs.infura.io/ipns/QmQ4QZh8nrsczdUEwTyfBope4THUhqxqc1fx6qYhhzZQei/index.html",
		},
		{
			name:   "Other gateway IPFS URL",
			input:  "https://some.other.gateway.com/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD",
			result: "https://ipfs.infura.io/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD",
		},
		{
			name:   "Other gateway IPFS URL with path",
			input:  "https://some.other.gateway.com/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD/index.html",
			result: "https://ipfs.infura.io/ipfs/QmbydiPQXL6YYMbsArTVVg9jjK9RzUbjUYX1xiw6XYwDoD/index.html",
		},
		{
			name:   "Other gateway IPNS URL",
			input:  "https://some.other.gateway.com/ipns/QmQ4QZh8nrsczdUEwTyfBope4THUhqxqc1fx6qYhhzZQei",
			result: "https://ipfs.infura.io/ipns/QmQ4QZh8nrsczdUEwTyfBope4THUhqxqc1fx6qYhhzZQei",
		},
		{
			name:   "Other gateway IPNS URL with path",
			input:  "https://some.other.gateway.com/ipns/QmQ4QZh8nrsczdUEwTyfBope4THUhqxqc1fx6qYhhzZQei/index.html",
			result: "https://ipfs.infura.io/ipns/QmQ4QZh8nrsczdUEwTyfBope4THUhqxqc1fx6qYhhzZQei/index.html",
		},
	}

	p, err := NewProvider()
	require.Nil(t, err, "unexpected error")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := p.GatewayURL(test.input)
			if test.err != nil {
				require.NotNil(t, err, "failed to obtain expected error")
				if err != nil {
					assert.Equal(t, test.err.Error(), err.Error(), "unexpected error value")
				}
			} else {
				require.Nil(t, err, "unexpected error")
				assert.Equal(t, test.result, result, "unexpected value")
			}
		})
	}
}
