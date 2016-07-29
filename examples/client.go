package main

import (
	"crypto"
	"flag"
	"fmt"
	"os"

	"github.com/litsea/api/parameter"
	"github.com/litsea/api/request"
)

func main() {
	var consumerKey *string = flag.String(
		"consumerkey",
		"",
		"Consumer Key")
	var consumerSecret *string = flag.String(
		"consumersecret",
		"",
		"Consumer Secret")
	var domain *string = flag.String(
		"domain",
		"",
		"API domain")

	flag.Parse()

	if len(*consumerKey) == 0 || len(*consumerSecret) == 0 || len(*domain) == 0 {
		fmt.Println("You must set the --consumerkey, --consumersecret and --domain flags.")
		os.Exit(1)
	}

	consumer := request.NewConsumer(*consumerKey, *consumerSecret)
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
	url := "http://" + *domain + "/role/list"
	params := make(map[string]string)
	params["page_size"] = "10"
	params["page"] = "1"

	r := c.NewRequest(url)
	response, _ := r.Get(params)

	fmt.Println("response: " + response)

	// post
	url = "http://" + *domain + "/role/create"
	params = make(map[string]string)
	params["role"] = parameter.Escape("测试 1")

	r = c.NewRequest(url)
	response, _ = r.Post(params)

	fmt.Println("response: " + response)

	consumer = request.NewCustomConsumer(*consumerKey, *consumerSecret, crypto.SHA256)
	consumer.Debug(true)
	addon_params = make(map[string]string)
	addon_params["_user"] = "lostsnow"
	addon_params["_user_realname"] = parameter.Escape("张三")
	addon_params["_userid"] = "1001"
	addon_params["_userip"] = "1.2.3.4"
	addon_params["alt"] = "json"
	consumer.AdditionalParams = addon_params
	c, _ = request.NewClient(consumer)

	// get
	url = "http://" + *domain + "/role/list"
	params = make(map[string]string)
	params["page_size"] = "5"
	params["page"] = "1"

	r = c.NewRequest(url)
	response, _ = r.Get(params)

	fmt.Println("response: " + response)
}
