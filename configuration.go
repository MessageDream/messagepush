package main

import (
//"github.com/Unknwon/goconfig"
)

type Config struct {
	DbConn string
}

var (
	Conf     *Config
	confFile string
)

func InitConf() {
	Conf = &Config{
		DbConn: "",
	}
}
