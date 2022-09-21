package twitch

import (
	"testing"

	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func TestPointerUtils(t *testing.T) {
	const val = "abc123"
	ptr := Pointer(val)
	assert.Equal(t, val, PointerValue(ptr))

	ptr = nil
	assert.Equal(t, "", PointerValue(ptr))
}
