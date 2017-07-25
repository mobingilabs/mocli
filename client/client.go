package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/parnurzeal/gorequest"
)

type GrClient struct {
	requester *gorequest.SuperAgent
	config    *Config
}

func NewGrClient(cnf *Config) *GrClient {
	return &GrClient{
		requester: gorequest.New(),
		config:    cnf,
	}
}

func (c *GrClient) Get(path string) (gorequest.Response, []byte, []error) {
	return c.requester.Get(c.url()+path).Set("Authorization", "Bearer "+c.config.AccessToken).EndBytes()
}

func (c *GrClient) PostU(path, payload string) (gorequest.Response, []byte, []error) {
	return c.requester.Post(c.url() + path).Send(payload).EndBytes()
}

func (c *GrClient) Put(path, payload string) (gorequest.Response, []byte, []error) {
	return c.requester.Put(c.url()+path).Set("Authorization", "Bearer "+c.config.AccessToken).Send(payload).EndBytes()
}

func (c *GrClient) Del(path string) (gorequest.Response, []byte, []error) {
	return c.requester.Delete(c.url()+path).Set("Authorization", "Bearer "+c.config.AccessToken).EndBytes()
}

func (c *GrClient) url() string {
	return c.config.RootUrl + "/" + c.config.ApiVersion
}

var Timeout int64

type Client struct {
	client *http.Client
	config *Config
}

func NewClient(cnf *Config) *Client {
	return &Client{
		client: &http.Client{
			Timeout: time.Second * time.Duration(Timeout),
		},
		config: cnf,
	}
}

func (c *Client) url() string {
	return c.config.RootUrl + "/" + c.config.ApiVersion
}

func (c *Client) GetHeaders(path string, values url.Values, hdrs http.Header) (http.Header, error) {
	req, err := http.NewRequest("GET", c.url()+path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
	for n, h := range hdrs {
		req.Header.Add(n, h[0])
	}

	req.URL.RawQuery = values.Encode()
	if d.Verbose {
		for n, h := range req.Header {
			d.Info(fmt.Sprintf("[headers-in] %s: %s", n, h))
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if d.Verbose {
		for n, h := range resp.Header {
			d.Info(fmt.Sprintf("[headers-out] %s: %s", n, h))
		}
	}

	defer resp.Body.Close()
	ret := resp.Header
	return ret, nil
}

func (c *Client) Get(path string, values url.Values, hdrs http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", c.url()+path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
	for n, h := range hdrs {
		req.Header.Add(n, h[0])
	}

	req.URL.RawQuery = values.Encode()
	if d.Verbose {
		for n, h := range req.Header {
			d.Info(fmt.Sprintf("[get-in] %s: %s", n, h))
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if d.Verbose {
		for n, h := range resp.Header {
			d.Info(fmt.Sprintf("[get-out] %s: %s", n, h))
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return body, errors.New(resp.Status)
	}

	return body, nil
}

func (c *Client) Del(path string, values url.Values) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.url()+path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
	req.URL.RawQuery = values.Encode()
	if d.Verbose {
		for n, h := range req.Header {
			d.Info(fmt.Sprintf("[del-in] %s: %s", n, h))
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if d.Verbose {
		for n, h := range resp.Header {
			d.Info(fmt.Sprintf("[del-out] %s: %s", n, h))
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
