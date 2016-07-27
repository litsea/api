package signature

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/litsea/api/parameter"
)

const (
	SIGNATURE_METHOD_HMAC = "HMAC-"
)

var HASH_METHOD_MAP = map[crypto.Hash]string{
	crypto.SHA1:   "SHA1",
	crypto.SHA256: "SHA256",
}

type Signer interface {
	Sign(message string) (string, error)
	Verify(message string, signature string) error
	SignatureMethod() string
	HashFunc() crypto.Hash
	Debug(enabled bool)
}

type HMACSigner struct {
	consumerSecret string
	hashFunc       crypto.Hash
	debug          bool
}

func NewSigner(consumerSecret string, hashFunc crypto.Hash) *HMACSigner {
	return &HMACSigner{
		consumerSecret: consumerSecret,
		hashFunc:       hashFunc,
	}
}

func NewHMACSigner(consumerSecret string) *HMACSigner {
	return NewSigner(consumerSecret, crypto.SHA1)
}

func (s *HMACSigner) Debug(enabled bool) {
	s.debug = enabled
}

func (s *HMACSigner) Sign(message string) (string, error) {
	secret := parameter.Escape(s.consumerSecret)
	if s.debug {
		fmt.Println("Signing:", message)
		fmt.Println("Secret:", secret)
	}

	h := hmac.New(s.HashFunc().New, []byte(secret))
	h.Write([]byte(message))
	rawSignature := h.Sum(nil)

	base64signature := base64.StdEncoding.EncodeToString(rawSignature)
	if s.debug {
		fmt.Println("Base64 signature:", base64signature)
	}
	return base64signature, nil
}

func (s *HMACSigner) Verify(message string, signature string) error {
	if s.debug {
		fmt.Println("Verifying Base64 signature:", signature)
	}
	validSignature, err := s.Sign(message)
	if err != nil {
		return err
	}

	if validSignature != signature {
		return errors.New("signature did not match")
	}

	return nil
}

func (s *HMACSigner) SignatureMethod() string {
	return SIGNATURE_METHOD_HMAC + HASH_METHOD_MAP[s.HashFunc()]
}

func (s *HMACSigner) HashFunc() crypto.Hash {
	return s.hashFunc
}
