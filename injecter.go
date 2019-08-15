package injecter

import (
	"strconv"

	"github.com/pkg/errors"
)

type Func struct {
	Name   string
	Return string // return type
}

// Inject
// udf: os_arch:library
// windows_386:[]byte{0}
// windows_amd64:[]byte{0}
// linux_386:[]byte{0}
// linux_amd64:[]byte{0}
func Inject(handle Handle, UDF map[string][]byte, Funcs []Func) error {

	return nil
}

// GetVersion
// "8.0.15"
func GetVersion(handle Handle) (string, error) {
	result, err := handle.Query("select @@version")
	if err != nil {
		return "", errors.WithMessage(err, "get version failed")
	}
	if ver, ok := result[0]["@@version"]; ok {
		return ver, nil
	}
	return "", errors.New("no version")
}

// GetMaxAllowedPacket
// 1073741824
func GetMaxAllowedPacket(handle Handle) (int, error) {
	const errPrefix = "get max_allowed_packet failed"
	result, err := handle.Query("select @@max_allowed_packet")
	if err != nil {
		return 0, errors.WithMessage(err, errPrefix)
	}
	value, ok := result[0]["@@max_allowed_packet"]
	if !ok {
		return 0, errors.New("no max_allowed_packet")
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.WithMessage(err, errPrefix)
	}
	return v, nil
}

func SetMaxAllowedPacket(handle Handle, value int) error {
	err := handle.Exec("set global max_allowed_packet = " + strconv.Itoa(value))
	if err != nil {
		return errors.WithMessage(err, "set max_allowed_packet failed")
	}
	return nil
}
