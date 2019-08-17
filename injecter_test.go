package injecter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// general

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

func TestInject(t *testing.T) {
	handle, err := Connect("192.168.1.13:3306", "root", "123456")
	require.NoError(t, err)
	err = Inject(handle, nil, nil)
	require.NoError(t, err)
	err = DropFunc(handle, "test")
	require.NoError(t, err)
}
