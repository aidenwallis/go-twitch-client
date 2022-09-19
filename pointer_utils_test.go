package twitch

import (
	"testing"

	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func TestPointer(t *testing.T) {
	const val = "abc123"
	ptr := Pointer(val)
	assert.Equal(t, val, *ptr)
}
