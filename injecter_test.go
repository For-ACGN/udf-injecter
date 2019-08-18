package injecter

import (
	"io/ioutil"
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

func TestInject50(t *testing.T) {
	handle, err := Connect("192.168.1.13:3306", "root", "123456")
	require.NoError(t, err)
	defer handle.Close()
	testInject(t, handle)
}

func TestInject51(t *testing.T) {
	handle, err := Connect("192.168.1.13:3307", "root", "123456")
	require.NoError(t, err)
	defer handle.Close()
	testInject(t, handle)
}

func testInject(t *testing.T, h Handle) {
	udfdata, err := ioutil.ReadFile("udf/windows/udf.dll")
	require.NoError(t, err)
	name := RandomStr(8) + "." + RandomStr(3)
	udfmap := map[string]*UDF{"windows_386": {name, udfdata}}
	funcs := []Func{{"udf_add", "integer"}}
	err = Inject(h, udfmap, funcs)
	require.NoError(t, err)
	err = DropFunc(h, "udf_add")
	require.NoError(t, err)
}
