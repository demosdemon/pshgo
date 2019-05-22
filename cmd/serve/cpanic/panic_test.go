package cpanic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/demosdemon/pshgo/cmd/serve/cpanic"
)

func TestRecover(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		assert.PanicsWithValue(t, assert.AnError, func() {
			defer Recover(nil)
			panic(assert.AnError)
		})
	})

	t.Run("no panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			defer Recover(func(p *Panic) {
				assert.FailNow(t, "recovered panic")
			})
			assert.True(t, true)
		})
	})

	t.Run("recovers", func(t *testing.T) {
		assert.NotPanics(t, func() {
			defer Recover(func(p *Panic) {
				assert.Equal(t, assert.AnError, p.Value)
			})
			panic(assert.AnError)
		})
	})
}

func TestPanic_String(t *testing.T) {
	p := Panic{}
	assert.Equal(t, "panic: <nil>\n\n", p.String())
}
