package injecter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	ver, err := GetVersion(testHanlde(t))
	require.NoError(t, err)
	t.Log(ver)
}

func TestGetMaxAllowedPacket(t *testing.T) {
	v, err := GetMaxAllowedPacket(testHanlde(t))
	require.NoError(t, err)
	t.Log(v)
}

func TestSetMaxAllowedPacket(t *testing.T) {
	handle := testHanlde(t)
	err := SetMaxAllowedPacket(handle, 1024)
	require.NoError(t, err)
	v, err := GetMaxAllowedPacket(testHanlde(t))
	require.NoError(t, err)
	require.Equal(t, 1024, v)
}
