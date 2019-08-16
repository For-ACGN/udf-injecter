package injecter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	handle := testHanlde(t)
	rows, err := handle.Query("show variables like '%max_allowed_packet%'")
	require.NoError(t, err)
	for i := 0; i < len(rows); i++ {
		for k, v := range rows[i] {
			t.Log(k, v)
		}
	}
	err = handle.Exec("show variables like '%max_allowed_packet%'")
	require.NoError(t, err)
}

// 8.0.15
func testHanlde(t *testing.T) Handle {
	handle, err := Connect("127.0.0.1:3306", "root", "123456")
	require.NoError(t, err)
	return handle
}
