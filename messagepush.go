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
	m.Post("/push/:appkey", func(params martini.Params, req *http.Request) string {
		appkey := params["appkey"]
		return pushHandler(&appkey, req)
	})
	m.Run()
}
