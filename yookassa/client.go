// Package yookassa implements all the necessary methods for working with YooMoney.
package yookassa

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	BaseURL = "https://api.yookassa.ru/v3/"
)

// Client works with YooMoney API.
type Client struct {
	client    http.Client
	accountId string
	secretKey string
}

func NewClient(accountId string, secretKey string) *Client {
	return &Client{
		client:    http.Client{},
		accountId: accountId,
		secretKey: secretKey,
	}
}

func (c *Client) makeRequest(
	ctx context.Context,
	method string,
	endpoint string,
	body []byte,
	params map[string]interface{},
	idempotencyKey string,
) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s", BaseURL, endpoint)

	var (
		req *http.Request
		err error
	)
	if ctx == nil {
		req, err = http.NewRequest(method, uri, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(body))
	}

	if err != nil {
		return nil, err
	}

	if idempotencyKey == "" {
		idempotencyKey = uuid.NewString()
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotence-Key", idempotencyKey)
	}

	req.SetBasicAuth(c.accountId, c.secretKey)

	if params != nil {
		q := req.URL.Query()
		for paramName, paramVal := range params {
			q.Add(paramName, fmt.Sprintf("%v", paramVal))
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
