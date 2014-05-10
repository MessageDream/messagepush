package main

import (
	"fmt"
	"github.com/MessageDream/go_apns"
	"net/http"
)

func pushHandler(appkey *string, req *http.Request) string {
	msg := &apns.Notification{}
	ParseJsonFromRequest(req, &msg)
	fmt.Println(msg)
	fmt.Println(appkey)
	return ""
}
