package twitch

import (
	"net/http"
	"testing"

	"github.com/aidenwallis/go-twitch-client/internal/testutils/assert"
)

func TestError(t *testing.T) {
	t.Run("empty message", func(t *testing.T) {
		err := NewError("", http.StatusBadRequest)
		assert.Equal(t, "[400] <unknown>", err.Error())
	})

	t.Run("with message", func(t *testing.T) {
		err := NewError("test", http.StatusBadRequest)
		assert.Equal(t, "[400] test", err.Error())
	})
}
