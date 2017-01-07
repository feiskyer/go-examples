package main

import (
	"fmt"

	"gopkg.in/gcfg.v1"
)

func main() {
	cfgStr := `
	[global]
	driver = cfg
	[section]
	options = option1
	options = option2
	log-dir = /var/log/gcfg`
	cfg := struct {
		Global struct {
			Driver string `gcfg:"driver"`
		}
		Section struct {
			LogDir  string   `gcfg:"log-dir"`
			Options []string `gcfg:"options"`
		}
	}{}

	if err := gcfg.ReadStringInto(&cfg, cfgStr); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", cfg)
}
