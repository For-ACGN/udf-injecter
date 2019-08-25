package main

import (
	"flag"
)

func main() {
	var (
		mode     string // "direct", "pma"(phpMyAdmin)
		host     string
		username string
		password string
		config   string
	)
	flag.StringVar(&mode, "m", "direct", "inject method")
	flag.StringVar(&host, "h", "127.0.0.1:3306", "MySQL server address or phpMyAdmin url")
	flag.StringVar(&username, "u", "root", "MySQL server username")
	flag.StringVar(&password, "p", "", "MySQL server password")
	flag.StringVar(&config, "c", "config.toml", "UDF config file path")
	flag.Parse()
}
