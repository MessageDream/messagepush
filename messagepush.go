package main

import (
	"github.com/go-martini/martini"
	"net/http"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	m := martini.Classic()
	m.Post("/push/:appid", func(params martini.Params, req *http.Request) string {
		appid := params["appid"]
		return pushHandler(&appid, req)
	})
	m.Run()
}
