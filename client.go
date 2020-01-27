package zirael

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 30 * time.Second

type IDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type IClient interface {
	Get(url string, headers http.Header) (*http.Response, error)
	Post(url string, body io.Reader, headers http.Header) (*http.Response, error)
}

type Client struct {
	apiID   string
	apiKey  string
	nonce   string
	client  IDoer
	timeout time.Duration
}

func NewClient(apiKey, apiID, nonce string, opts... Option) *Client{
	client := Client{
		apiID:   apiID,
		apiKey:  apiKey,
		nonce:   nonce,
		timeout: defaultTimeout,
	}
	for _, opt := range opts{
		opt(&client)
	}
	if client.client == nil {
		client.client = &http.Client{
			Timeout: client.timeout,
		}
	}
	return &client
}

func (c *Client) Get(url string, headers http.Header) (*http.Response, error){
	var response *http.Response
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}

	amx, err := c.AmxTokenGenerate(url, http.MethodGet, nil)

	if err != nil {
		return response, errors.Wrap(err, "GET - amx token creation failed")
	}

	if headers != nil {
		request.Header = headers
	}
	request.Header.Add("Authorization", fmt.Sprintf("amx %s", amx))

	return c.Do(request)
}

func (c *Client) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return response, errors.Wrap(err, "POST - request creation failed")
	}

	amx, err := c.AmxTokenGenerate(url, http.MethodPost, nil)

	if err != nil {
		return response, errors.Wrap(err, "POST - amx token creation failed")
	}

	if headers != nil {
		request.Header = headers
	}

	request.Header.Add("Authorization", fmt.Sprintf("amx %s", amx))

	return c.Do(request)
}

func (c *Client) Do(request *http.Request) (response *http.Response, err error)  {
	response, err = c.client.Do(request)
	return
}

func (c *Client) AmxTokenGenerate(urlPath, method string, body []byte) (string, error) {
	method = strings.ToUpper(method)
	timestamp := time.Now().UnixNano() / 1000000000
	percentEncoded := strings.ToLower(url.QueryEscape(urlPath))
	if body == nil {
		body = MD5Hash([]byte("\"\""))
	}else{
		body = MD5Hash(body)
	}
	md64 := base64.StdEncoding.EncodeToString(body)
	signature := fmt.Sprintf("%s%s%s%d%s%s", c.apiID, method, percentEncoded, timestamp, c.nonce, md64)
	signature64, err := ComputeHMAC256(signature, c.apiKey)

	if err != nil {
		return "", errors.Wrap(err, "Failed to compute base64 signature!")
	}

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v:%v:%v", c.apiID, signature64, c.nonce, timestamp))), nil

}

func MD5Hash(data [] byte) [] byte {
	str := string(data)
	res := strings.ReplaceAll(str, " ", "")
	res = strings.ReplaceAll(res, "\r", "")
	res = strings.ReplaceAll(res, "\n", "")
	res = strings.ReplaceAll(res, "\t", "")
	h := md5.New()
	h.Write([]byte(res))
	return h.Sum(nil)
}

func ComputeHMAC256(payload string, secret string) (string, error) {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	_, err := h.Write([]byte(payload))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
