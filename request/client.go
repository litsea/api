package request

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type client struct {
	httpClient *http.Client
	additionalParams  map[string]string
}

func NewClient(c *Consumer) (*client, error) {
	httpClient, err := c.MakeHttpClient()
	if err != nil {
		return nil, err
	}

	return &client{
		httpClient: httpClient,
	}, nil
}

type request struct {
	client *client
	url        string
	userParams map[string]string
}

func (c *client) NewRequest(urlStr string) (*request) {
	return &request{
		client: c,
		url: urlStr,
	}
}

func (c *client) AdditionalParams(userParams map[string]string) {
	c.additionalParams = userParams
}

func (r *request) AddGetParameters(userParams map[string]string) error {
	u, err := url.Parse(r.url)
	if err != nil {
		return err
	}

	q := u.Query()
	for k, v := range userParams {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	r.url = u.String()

	return nil
}

func (r *request) AddPostParameters(userParams map[string]string) error {
	u, err := url.Parse(r.url)
	if err != nil {
		return err
	}

	q := u.Query()
	for k, v := range userParams {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	r.url = u.String()

	return nil
}

func (r *request) Get(userParams map[string]string) (string, error) {
	r.AddGetParameters(userParams)
	response, err := r.client.httpClient.Get(r.url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	return string(bits), err
}

func (r *request) Post(userParams map[string]string) (string, error) {
	vals := url.Values{}
		for k, v := range userParams {
			vals.Add(k, v)
	}

	response, err := r.client.httpClient.PostForm(r.url, vals)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	return string(bits), err
}
