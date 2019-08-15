package injecter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	address  = "127.0.0.1:3306"
	username = "root"
	password = "123456"
)

func TestHandle_Query(t *testing.T) {
	handle := testHanlde(t)
	rows, err := handle.Query("show variables like '%max_allowed_packet%'")
	require.NoError(t, err)
	for i := 0; i < len(rows); i++ {
		for k, v := range rows[i] {
			t.Log(k, v)
		}
	}
}

func testHanlde(t *testing.T) Handle {
	handle, err := Connect(address, username, password)
	require.NoError(t, err)
	return handle
}
