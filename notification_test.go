package notify

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpiration(t *testing.T) {
	require.EqualValues(t, 0, ExpireTimeoutNever.Milliseconds())
	require.EqualValues(t, -1, ExpireTimeoutSetByNotificationServer.Milliseconds())

	// test assignment compiles:
	n := Notification{}
	n.ExpireTimeout = ExpireTimeoutNever
	n.ExpireTimeout = ExpireTimeoutSetByNotificationServer
}
