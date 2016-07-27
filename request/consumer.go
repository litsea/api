package request

import (
	"crypto"
	"net/http"

	"github.com/litsea/api/signature"
)

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type Consumer struct {
	consumerKey    string
	debug          bool
	HttpClient     HttpClient
	clock          clock
	nonceGenerator nonceGenerator
	signer         signature.Signer

	AdditionalParams  map[string]string
}

func newConsumer(consumerKey string, httpClient *http.Client) *Consumer {
	clock := &defaultClock{}
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Consumer{
		consumerKey:    consumerKey,
		clock:          clock,
		HttpClient:     httpClient,
		nonceGenerator: newLockedNonceGenerator(clock),
	}
}

func (c *Consumer) Debug(enabled bool) {
	c.debug = enabled
	c.signer.Debug(enabled)
}

func NewConsumer(consumerKey string, consumerSecret string) *Consumer {
	consumer := newConsumer(consumerKey, nil)
	consumer.signer = signature.NewHMACSigner(consumerSecret)

	return consumer
}

func NewCustomHttpClientConsumer(consumerKey string, consumerSecret string,
	httpClient *http.Client) *Consumer {
	consumer := newConsumer(consumerKey, httpClient)

	consumer.signer = signature.NewHMACSigner(consumerSecret)

	return consumer
}

func NewCustomConsumer(consumerKey string, consumerSecret string,
	hashFunc crypto.Hash, httpClient *http.Client) *Consumer {
	consumer := newConsumer(consumerKey, httpClient)

	consumer.signer = signature.NewHMACSigner(consumerSecret)

	return consumer
}
