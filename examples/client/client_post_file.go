package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/litsea/api/parameter"
	"github.com/litsea/api/request"
	"io/ioutil"
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
	var filename *string = flag.String(
		"filename",
		"",
		"File path")

	flag.Parse()

	if len(*consumerKey) == 0 || len(*consumerSecret) == 0 || len(*domain) == 0 || len(*filename) == 0 {
		fmt.Println("You must set the --consumerkey, --consumersecret, --domain and --filename flags.")
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

	// post file
	fileBuf, _ := ioutil.ReadFile(*filename)

	url := "http://" + *domain + "/file/upload?xxx="
	params := make(map[string]string)
	params["__file_filename"] = "test.png"
	params["__file_formname"] = "files"
	params["__file_data"] = string(fileBuf)
	params["aa"] = "222"

	r := c.NewRequest(url)
	response, _ := r.Post(params)
	fmt.Println("response: " + response)
}
