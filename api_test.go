package infura

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var infuraAPIKey string
var infuraAPISecret string

func init() {
}

func TestNewProvider(t *testing.T) {
	tests := []struct {
		err error
	}{
		{ // 0
		},
	}

	for i, test := range tests {
		_, err := NewProvider()
		if test.err != nil {
			assert.NotNil(t, err, fmt.Sprintf("missing expected error at test %d", i))
			assert.Equal(t, test.err.Error(), err.Error(), fmt.Sprintf("unexpected error value at test %d", i))
		} else {
			assert.Nil(t, err, fmt.Sprintf("unexpected error at test %d", i))
		}
	}
}
