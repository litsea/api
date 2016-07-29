package request

import (
	"crypto"
	"net/http"

	"github.com/litsea/api/signature"
)

type Consumer struct {
	consumerKey    string
	debug          bool
	HttpClient     *http.Client
	clock          clock
	nonceGenerator nonceGenerator
	signer         signature.Signer

	AdditionalParams map[string]string
}

func newConsumer(consumerKey string) *Consumer {
	clock := &defaultClock{}
	httpClient := &http.Client{}
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
	consumer := newConsumer(consumerKey)
	consumer.signer = signature.NewHMACSigner(consumerSecret)

	return consumer
}

func NewCustomConsumer(consumerKey string, consumerSecret string, hashFunc crypto.Hash) *Consumer {
	consumer := newConsumer(consumerKey)

	consumer.signer = signature.NewSigner(consumerSecret, hashFunc)

	return consumer
}
