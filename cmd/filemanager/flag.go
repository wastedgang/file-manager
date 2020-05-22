package main

import (
	"flag"
)

var (
	configPath string
	isDebug bool
)

func initFlag() {
	flag.StringVar(&configPath, "c", "run.ini", "配置文件路径 `config`")
	flag.BoolVar(&isDebug, "d", false, "set `debug`")
	flag.Parse()
}