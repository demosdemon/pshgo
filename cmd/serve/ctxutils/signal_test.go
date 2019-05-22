package ctxutils_test

import (
	"context"
	"github.com/demosdemon/pshgo/cmd/serve/ctxutils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestCancelContextWithSignal(t *testing.T) {
	ctx, _ := ctxutils.CancelContextWithSignal(context.Background(), os.Interrupt, os.Kill)
	err := raise(os.Interrupt)
	require.NoError(t, err)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
	case <-ticker.C:
		require.Fail(t, "timeout elapsed")
	}

	require.EqualError(t, ctx.Err(), context.Canceled.Error())
}

func raise(sig os.Signal) error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Signal(sig)
}
