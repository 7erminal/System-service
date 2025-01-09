package api

import (
	"net/http"

	"github.com/astaxie/beego/httplib"
)

type Client struct {
	Request *Request
	Type_   string
}

func (c *Client) SendRequest() (*http.Response, error) {
	// TODO: BaseURLとPathのvalidation
	beegoRequest := httplib.NewBeegoRequest(
		c.Request.BaseURL+c.Request.Path,
		c.Request.Method.String())
	for key, value := range c.Request.HeaderField {
		beegoRequest.Header(key, value)
	}
	if c.Type_ == "params" {
		for key, value := range c.Request.Params {
			beegoRequest.Param(key, value)
		}
	} else if c.Type_ == "body" {
		beegoRequest.JSONBody(c.Request.Params)
	}
	res, err := beegoRequest.Response()
	if err != nil {
		return nil, err
	}
	return res, nil
}
