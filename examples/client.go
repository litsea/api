package main

import (
	"fmt"

	"github.com/litsea/api/request"
	"github.com/litsea/api/parameter"
)

func main() {
	consumer := request.NewConsumer("xx", "oo")
	consumer.Debug(true)
	addon_params := make(map[string]string)
	addon_params["_user"] = "lostsnow"
	addon_params["_user_realname"] = parameter.Escape("张三")
	addon_params["_userid"] = "1001"
	addon_params["_userip"] = "1.2.3.4"
	addon_params["alt"] = "json"
	consumer.AdditionalParams = addon_params
	c, _ := request.NewClient(consumer)

	// get
	url := "http://api.com/role/list"
	params := make(map[string]string)
	params["page_size"] = "10"
	params["page"] = "1"

	r := c.NewRequest(url)
	response, _ := r.Get(params)

	fmt.Println("response: " + response)

	// post
	url = "http://api.com/role/create"
	params = make(map[string]string)
	params["role"] = parameter.Escape("测试 1")

	r = c.NewRequest(url)
	response, _ = r.Post(params)

	fmt.Println("response: " + response)
}
