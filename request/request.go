package request

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/litsea/api/parameter"
)

const (
	AUTH_PARAM_CONSUMER_KEY     = "auth_consumer_key"
	AUTH_PARAM_NONCE            = "auth_nonce"
	AUTH_PARAM_SIGNATURE_METHOD = "auth_signature_method"
	AUTH_PARAM_SIGNATURE        = "auth_signature"
	AUTH_PARAM_TIMESTAMP        = "auth_timestamp"
)

type clock interface {
	Seconds() int64
	Nanos() int64
}

type defaultClock struct{}

func (*defaultClock) Seconds() int64 {
	return time.Now().Unix()
}

func (*defaultClock) Nanos() int64 {
	return time.Now().UnixNano()
}

type nonceGenerator interface {
	Int63() int64
}

type lockedNonceGenerator struct {
	nonceGenerator nonceGenerator
	lock           sync.Mutex
}

func newLockedNonceGenerator(c clock) *lockedNonceGenerator {
	return &lockedNonceGenerator{
		nonceGenerator: rand.New(rand.NewSource(c.Nanos())),
	}
}

func (n *lockedNonceGenerator) Int63() int64 {
	n.lock.Lock()
	r := n.nonceGenerator.Int63()
	n.lock.Unlock()
	return r
}

type RoundTripper struct {
	consumer *Consumer
}

func (c *Consumer) MakeRoundTripper() (*RoundTripper, error) {
	return &RoundTripper{consumer: c}, nil
}

func (c *Consumer) MakeHttpClient() (*http.Client, error) {
	return &http.Client{
		Transport: &RoundTripper{consumer: c},
	}, nil
}

func cloneReq(src *http.Request) *http.Request {
	dst := &http.Request{}
	*dst = *src

	dst.Header = make(http.Header, len(src.Header))
	for k, s := range src.Header {
		dst.Header[k] = append([]string(nil), s...)
	}

	if src.URL != nil {
		dst.URL = cloneURL(src.URL)
	}

	return dst
}

func cloneURL(src *url.URL) *url.URL {
	dst := &url.URL{}
	*dst = *src

	return dst
}

func (rt *RoundTripper) RoundTrip(userRequest *http.Request) (*http.Response, error) {
	serverRequest := cloneReq(userRequest)

	allParams := rt.consumer.baseParams(rt.consumer.consumerKey, rt.consumer.AdditionalParams)

	authParams := allParams.Clone()

	userParams, err := parseBody(serverRequest)
	if err != nil {
		return nil, err
	}
	paramPairs := parameter.ParamsToSortedPairs(userParams)

	for i := range paramPairs {
		allParams.AddUnescaped(paramPairs[i].Key, paramPairs[i].Value)
	}

	signingURL := cloneURL(serverRequest.URL)
	if host := serverRequest.Host; host != "" {
		signingURL.Host = host
	}
	baseString := rt.consumer.requestBaseString(serverRequest.Method, normalizeUrl(signingURL), allParams)

	signature, err := rt.consumer.signer.Sign(baseString)
	if err != nil {
		return nil, err
	}

	authParams.AddUnescaped(AUTH_PARAM_SIGNATURE, signature)

	q := serverRequest.URL.Query()
	for _, key := range authParams.Keys() {
		for _, value := range authParams.Get(key) {
			q.Add(key, value)
		}
	}
	serverRequest.URL.RawQuery = q.Encode()

	if rt.consumer.debug {
		fmt.Printf("Request: %v\n", serverRequest)
	}

	resp, err := rt.consumer.HttpClient.Do(serverRequest)

	if err != nil {
		return resp, err
	}

	return resp, nil
}

func normalizeUrl(u *url.URL) string {
	var buf bytes.Buffer
	buf.WriteString(u.Scheme)
	buf.WriteString("://")
	buf.WriteString(u.Host)
	buf.WriteString(u.Path)

	return buf.String()
}

func parseBody(request *http.Request) (map[string]string, error) {
	userParams := map[string]string{}

	if request.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		// Most of the time we get parameters from the query string:
		for k, vs := range request.URL.Query() {
			if len(vs) != 1 {
				return nil, errors.New("Must have exactly one value per param")
			}

			userParams[k] = vs[0]
		}
	} else {
		// x-www-form-urlencoded parameters come from the body instead:
		defer request.Body.Close()
		originalBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}

		// If there was a body, we have to re-install it
		// (because we've ruined it by reading it).
		request.Body = ioutil.NopCloser(bytes.NewReader(originalBody))

		params, err := url.ParseQuery(string(originalBody))
		if err != nil {
			return nil, err
		}

		for k, vs := range params {
			if len(vs) != 1 {
				return nil, errors.New("Must have exactly one value per param")
			}

			userParams[k] = vs[0]
		}
	}

	return userParams, nil
}

func (c *Consumer) baseParams(consumerKey string, additionalParams map[string]string) *parameter.OrderedParams {
	params := parameter.NewOrderedParams()
	params.Add(AUTH_PARAM_SIGNATURE_METHOD, c.signer.SignatureMethod())
	params.Add(AUTH_PARAM_TIMESTAMP, strconv.FormatInt(c.clock.Seconds(), 10))
	params.Add(AUTH_PARAM_NONCE, strconv.FormatInt(c.nonceGenerator.Int63(), 10))
	params.Add(AUTH_PARAM_CONSUMER_KEY, consumerKey)
	for key, value := range additionalParams {
		params.AddUnescaped(key, value)
	}
	return params
}

func (c *Consumer) requestBaseString(method string, url string, params *parameter.OrderedParams) string {
	result := method + "&" + parameter.Escape(url) + "&"
	param := ""
	for pos, key := range params.Keys() {
		for innerPos, value := range params.Get(key) {
			if pos+innerPos != 0 {
				param += "&"
			}

			param += fmt.Sprintf("%s=%s", parameter.Escape(key), parameter.Escape(value))
		}
	}
	result += parameter.Escape(param)
	return result
}
